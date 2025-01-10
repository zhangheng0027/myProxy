package net

import (
	"fmt"
	http_dialer "github.com/mwitkow/go-http-dialer"
	"net/http"
	"net/url"
	"testing"
	"time"
)

// 9955
func Test1(t *testing.T) {

	purl, _ := url.Parse("http://192.168.115.2:8080")
	tunnel := http_dialer.New(purl)

	//dial, err := tunnel.Dial("tcp", "www.baidu.com:80")
	//if err != nil {
	//	t.Error(err)
	//}

	//http.Get("https://www.baidu.com")

	client := http.Client{
		Transport: &http.Transport{Dial: tunnel.Dial},
	}

	start := time.Now()
	resp, err := client.Get("https://appscross.com")
	if err != nil {
		t.Error(err)
	}
	t.Log(resp)

	// 统计耗时
	t.Log(time.Since(start))

}

func Test2(t *testing.T) {
	//Listen(":9966")
	checkSSR(20)
}

func Test3(t *testing.T) {
	//
	server := NewSSRServer("")
	Configuration.SSRServers = append(Configuration.SSRServers, server)
	checkSSR(20)
}

func Test4(t *testing.T) {
	fmt.Println(dnsResolve("mcs.doubao.com", "443"))
}

func Test5(t *testing.T) {
	ListenUDP(":9977")
}
