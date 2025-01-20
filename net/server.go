package net

import (
	"bufio"
	"fmt"
	httpdialer "github.com/mwitkow/go-http-dialer"
	"github.com/zhangheng0027/ratelimit-plus"
	"github.com/zhangheng0027/shadowsocksR/client"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

func Tolimit(conn net.Conn) *ConnLimit {
	addr := conn.RemoteAddr()
	tcpAddr := addr.(*net.TCPAddr)
	if tcpAddr.IP.IsLoopback() {
		if _, ok := (Configuration.WhiteListIp)["localhost"]; ok {
			return ToLimitNetConn(conn, nil, nil)
		}
	}

	_, ok := (Configuration.WhiteListIp)[tcpAddr.IP.String()]
	if ok {
		return ToLimitNetConn(conn, nil, nil)
	} else {
		return ToLimitNetConn(conn, ratelimit.NewBucketWithRate(float64(Configuration.DefaultIpReadLimit), Configuration.DefaultIpReadLimit*10),
			ratelimit.NewBucketWithRate(float64(Configuration.DefaultIpWriteLimit), Configuration.DefaultIpWriteLimit*10))
	}
}

func remoteDial(conn *ConnLimit, addr string) (net.Conn, error) {
	proxyS := getProxyServer()
	conn.addUpStream(proxyS.readLimit, proxyS.writeLimit)
	sp := strings.Split(addr, ":")
	if isSSR(sp[0]) {
		ssr := getSSRStr()
		if nil == ssr {
			return remoteDialNoSSR(proxyS, sp)
		}
		proxy, err := ssr.ssr.DialProxy("tcp", addr, proxyS.tunnel)
		if err != nil {
			return nil, err
		}
		conn.addUpStream(ssr.readLimit, ssr.writeLimit)
		return proxy, err
	} else {
		return remoteDialNoSSR(proxyS, sp)
	}
}

func remoteDialNoSSR(proxyS *ProxyServer, sp []string) (net.Conn, error) {
	fmt.Println("use proxy ", proxyS.addr, " <-> ", strings.Join(sp, ":"))
	resolve := dnsResolve(sp[0], sp[1])
	ip, i := resolve.GetIp()
	// 需要记录建立连接的耗时
	now := time.Now()
	dial, err := proxyS.tunnel.Dial("tcp", ip+":"+sp[1])
	if err != nil {
		resolve.removeIp(ip, i)
	} else {
		resolve.AddOneDialed(ip, i, time.Since(now))
	}
	return dial, err
}

type ProxyServer struct {
	tunnel     *httpdialer.HttpTunnel
	readLimit  *ratelimit.Bucket
	writeLimit *ratelimit.Bucket
	addr       string
}

type SSRServer struct {
	ssr           *client.SSR
	availableFlag bool
	delayTime     time.Duration
	readLimit     *ratelimit.Bucket
	writeLimit    *ratelimit.Bucket
}

func NewSSRServer(url string) *SSRServer {
	ssr1, err := client.NewSSR1(url)
	if err != nil {
		return nil
	}
	return &SSRServer{
		ssr:        ssr1,
		readLimit:  ratelimit.NewBucketWithRate(float64(DefaultReadLimit), DefaultReadLimit*10),
		writeLimit: ratelimit.NewBucketWithRate(float64(DefaultReadLimit), DefaultReadLimit*10),
	}
}

func NewProxyServer(proxyAddr string, readSp int64, writeSp int64) *ProxyServer {
	url, _ := url.Parse(proxyAddr)
	return &ProxyServer{
		tunnel:     httpdialer.New(url),
		readLimit:  ratelimit.NewBucketWithRate(float64(readSp), readSp*10),
		writeLimit: ratelimit.NewBucketWithRate(float64(writeSp), writeSp*10),
		addr:       proxyAddr,
	}
}

func (ssr *SSRServer) DialDns(network, addr string) (net.Conn, error) {
	sp := strings.Split(addr, ":")
	resolve := dnsResolve(sp[0], sp[1])
	ip, _ := resolve.GetRandIp()
	return ssr.ssr.Dial(network, ip+":"+sp[1])
}

func (ps *ProxyServer) usageRate() float64 {
	return float64(ps.readLimit.UsageRate()) / float64(ps.readLimit.Capacity())
}

func Init(url string) {
	loadSSR(url)
	checkSSR(20)
	go func() {
		ticker := time.Tick(10 * time.Minute)
		for range ticker {
			checkSSR(20)
		}
	}()

	go func() {
		ticker := time.Tick(2 * time.Hour)
		for range ticker {
			loadSSR(url)
		}
	}()
}

func loadSSR(url string) {

	client1 := &http.Client{
		Transport: &http.Transport{
			Proxy: nil,
		},
	}
	resp, err := client1.Get(url)
	if err != nil {
		fmt.Println("更新 ssr 失败")
		file, err := os.ReadFile("ssr.conf")
		if err != nil {
			fmt.Println("读取 ssr.conf 文件错误", err)
			return
		}
		s := string(file)
		encodeSSRBody(s)
		return
	}
	defer resp.Body.Close()
	bu := make([]byte, 1024*50)
	n, err := resp.Body.Read(bu)
	str := client.DecodeBase64ToStr(string(bu[:n]))
	file, err := os.Create("ssr.conf")
	if err != nil {
		fmt.Println("create file error", err)
		return
	}
	defer file.Close()
	bufio.NewWriter(file).WriteString(str)
	encodeSSRBody(str)
}

func encodeSSRBody(str string) {
	split := strings.Split(str, "\n")
	servers := make([]*SSRServer, 0, 40)
	for _, s := range split {
		if len(s) > 0 {
			servers = append(servers, NewSSRServer(s))
		}
	}
	Configuration.SSRServers = servers
}

func checkSSR(w int64) {
	var g sync.WaitGroup
	ssrServers := Configuration.SSRServers
	g.Add(len(ssrServers))
	for _, ssr := range ssrServers {
		go func(ssr *SSRServer) {
			defer g.Done()
			client := http.Client{
				Transport: &http.Transport{Dial: ssr.DialDns},
			}
			// 设置 headers

			req, err := http.NewRequest("GET", "https://www.baidu.com", nil)

			now := time.Now()
			resp, err := client.Do(req)
			if err != nil {
				ssr.availableFlag = false
				return
			}
			// 获取延迟时间
			ssr.delayTime = time.Since(now)
			if resp.StatusCode != 200 {
				ssr.availableFlag = false
			} else {
				ssr.availableFlag = true
				fmt.Println("ssr server is available", ssr.ssr.Remarks)
			}
		}(ssr)
	}

	done := make(chan struct{})
	go func() {
		g.Wait()
		close(done)
	}()

	select {
	case <-done:
		// All goroutines completed within the timeout
	case <-time.After(time.Duration(w) * time.Second):
		// Timeout reached
	}

	servers := make([]*SSRServer, 0, 40)
	for _, ssr := range ssrServers {
		if ssr.availableFlag {
			servers = append(servers, ssr)
		}
	}
	Configuration.availableSSR = servers

	fmt.Printf("check ssr done %d 可用\n", len(servers))
}
