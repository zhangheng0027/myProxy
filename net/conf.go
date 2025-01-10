package net

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"strings"
	"sync"
)

type Conf struct {
	// 默认限速
	DefaultIpReadLimit  int64
	DefaultIpWriteLimit int64

	// 白名单 ip, 不在名单中会进行限速
	WhiteListIp map[string]bool

	ProxyServer []*ProxyServer

	SSRServers []*SSRServer

	availableSSR []*SSRServer

	RouteSSRUrl map[string]bool

	dnsCache map[string]*dns

	mu sync.Mutex
}

var Configuration *Conf = &Conf{
	DefaultIpReadLimit:  5000 * 1024,
	DefaultIpWriteLimit: 5000 * 1024,
	WhiteListIp:         make(map[string]bool),
	ProxyServer:         make([]*ProxyServer, 0, 5),
	SSRServers:          make([]*SSRServer, 0, 50),
	availableSSR:        make([]*SSRServer, 0, 50),
	// 通过 ssr 进行请求的 url
	RouteSSRUrl: make(map[string]bool, 50),
	dnsCache:    make(map[string]*dns, 500),
}

func getProxyServer() *ProxyServer {
	return Configuration.ProxyServer[rand.Intn(len(Configuration.ProxyServer))]
}

func getSSRStr() *SSRServer {
	if len(Configuration.availableSSR) == 0 {
		return nil
	}
	return Configuration.availableSSR[rand.Intn(len(Configuration.availableSSR))]
}

func isSSR(addr string) bool {
	//addr = strings.Split(addr, ":")[0]
	_, ok := Configuration.RouteSSRUrl[addr]
	if ok {
		return true
	}
	return _isSSR(strings.Split(addr, "."))
}

func _isSSR(addrs []string) bool {
	if len(addrs) <= 1 {
		return false
	}
	url := "*." + strings.Join(addrs, ".")
	_, ok := Configuration.RouteSSRUrl[url]
	if ok {
		return true
	}
	return _isSSR(addrs[1:])
}

func dnsResolve(addr string, port string) *dns {
	s, ok := Configuration.dnsCache[addr]
	if ok && !s.needDns() {
		return s
	}

	Configuration.mu.Lock()
	defer Configuration.mu.Unlock()

	s, ok = Configuration.dnsCache[addr]
	if ok && !s.needDns() {
		return s
	}
	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			dialer := &net.Dialer{}
			// 使用 8.8.8.8 作为 DNS 服务器，端口为 53
			return dialer.DialContext(ctx, "udp", "8.8.8.8:53")
		},
	}

	// 查找 A 记录（IPv4 地址）
	ips, err := resolver.LookupHost(context.Background(), addr)
	var d *dns
	if err != nil {
		fmt.Println("Error resolving DNS:", err)
		d = newDns(addr, port, nil)
	} else {
		d = newDns(addr, port, ips)
	}
	Configuration.dnsCache[addr] = d
	return d
}
