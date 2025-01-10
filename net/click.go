package net

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/url"
	"strings"
)

func Listen(address string) {

	ad, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	listener, err := net.ListenTCP("tcp", ad)
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}

	//listener, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	fmt.Println("Listening on", address)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
		}
		go handle(conn)
	}
}

func handle(conn net.Conn) {
	defer conn.Close()
	con := Tolimit(conn)
	var buf [10240]byte
	n, err := con.Read(buf[:])
	if err != nil {
		fmt.Println("读取错误", err)
	}

	indexN := bytes.IndexByte(buf[:], '\n')

	if indexN <= 0 {
		fmt.Println("解析失败")
		return
	}

	oneLine := string(buf[:indexN])

	var method, URL, address string
	fmt.Sscanf(oneLine, "%s%s", &method, &URL)

	if method == "CONNECT" {
		addr := oneLine[7+1 : len(oneLine)-9-1]
		address = addr
		remoteCon, err := remoteDial(con, address)
		if err != nil {
			fmt.Println("链接失败:", err, address)
			return
		}
		fmt.Fprint(con, "HTTP/1.1 200 Connection established\r\n\r\n")
		go io.Copy(remoteCon, con)
		io.Copy(con, remoteCon)
		return
	}
	address, err = ExtractAddressFromOtherRequestLine(oneLine)
	remoteCon, err := remoteDial(con, address)
	if err != nil {
		fmt.Println("链接失败:", err, address)
		return
	}

	var requestLine = string(buf[:indexN+1])
	//如果使用 http 协议，需将从客户端得到的 http 请求转发给服务端
	clienthost, _, err := net.SplitHostPort(con.RemoteAddr().String())
	if err != nil {
		log.Println(err)
		return
	}
	//log.Println("clienthost:", clienthost)
	//log.Println("clientport:", port)
	forwarded := fmt.Sprintf(
		"for=%s;host=%s;proto=%s",
		clienthost, // 代理自己的标识或IP地址
		address,    // 原始请求的目标主机名
		"http",     // 或者 "https" 根据实际协议
	)
	var headers = map[string]string{"Forwarded": forwarded}
	shouldReturn := WriteRequestLineAndHeadersWithRequestURI(requestLine, remoteCon, n, buf, headers)
	if shouldReturn {
		return
	}
	go io.Copy(remoteCon, con)
	io.Copy(con, remoteCon)
}

// WriteRequestLineAndHeadersWithRequestURI 将请求行和头部信息写入服务器连接
func WriteRequestLineAndHeadersWithRequestURI(requestLine string, server net.Conn, n int, b [10240]byte, headers map[string]string) bool {
	/*有的服务器不支持这种 "GET http://speedtest.cn/ HTTP/1.1" */
	output, err := RemoveURLPartsLeaveMethodRequestURIVersion(requestLine)
	if err != nil {
		log.Println("解析错误 ", err)
		return true
	}
	server.Write([]byte(output))

	for k, v := range headers {
		server.Write([]byte(k + ": " + v + "\r\n"))
	}
	server.Write(b[len(requestLine):n])
	return false
}

func ExtractAddressFromOtherRequestLine(line string) (string, error) {
	var address string
	domain, port, err := ExtractDomainAndPort(line)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	} else {
		address = domain + ":" + port
	}
	return address, nil
}

func ExtractAddressFromConnectRequestLine(line string) string {
	return line[7+1 : len(line)-9-1]
}
func RemoveURLPartsLeaveMethodRequestURIVersion(requestLine string) (string, error) {
	/* 有的服务器不支持这种 "GET http://speedtest.cn/ HTTP/1.1" */
	// 正则表达式用于匹配 http(s)://[domain(:port)] 部分
	parts := strings.SplitN(requestLine, " ", 3)
	if len(parts) < 3 {
		return "", fmt.Errorf("invalid http request line")
	}

	// 获取请求目标
	requestTarget := parts[1]

	// 解析URL
	u, err := url.Parse(requestTarget)
	if err != nil {
		return "", fmt.Errorf("failed to parse url: %w", err)
	}
	var cleanedUrl = parts[0] + " " + u.RequestURI() + " " + parts[2]
	/* "GET / HTTP/1.1" */
	return cleanedUrl, nil
}
func ExtractDomainAndPort(requestLine string) (string, string, error) {
	/* "GET http://speedtest.cn/ HTTP/1.1" */
	// 分割字符串以获取URL部分
	parts := strings.SplitN(requestLine, " ", 3)
	if len(parts) < 3 {
		return "", "", fmt.Errorf("invalid http request line")
	}

	// 获取请求目标
	requestTarget := parts[1]

	// 解析URL
	u, err := url.Parse(requestTarget)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse url: %w", err)
	}

	// 提取域名
	domain := u.Hostname()

	// 提取端口
	port := u.Port()
	if port == "" {
		// 如果端口未指定，则根据协议使用默认端口
		if u.Scheme == "http" {
			port = "80"
		} else if u.Scheme == "https" {
			port = "443"
		}
	}
	if IsIPv6(domain) {
		domain = "[" + domain + "]"
	}

	/* 需要识别ipv6地址 */
	/* Domain: speedtest.cn, Port: 80 */
	return domain, port, nil
}
func IsIPv6(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	return ip != nil && ip.To16() != nil && ip.To4() == nil
}
