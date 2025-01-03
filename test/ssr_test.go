package test

import (
	"github.com/joho/godotenv"
	http_dialer "github.com/mwitkow/go-http-dialer"
	"github.com/zhangheng0027/shadowsocksR/client"
	"net"
	"net/url"
	"os"
	"testing"
	"time"
)

func TestShadowsocksR(t *testing.T) {
	ssrurl, b := os.LookupEnv("ssrurl")
	if !b {
		t.Error("ssrurl not found")
	}

	purl, _ := url.Parse("http://localhost:12801")
	tunnel := http_dialer.New(purl)

	ssr1, err := client.NewSSR2(ssrurl, tunnel)

	//ssr1, err := client.NewSSR1(ssrurl)
	if err != nil {
		t.Error(err)
	}
	dial, err := ssr1.Dial("tcp", "www.google.com:443")

	dial.Write([]byte("hello"))
}

func TestMain(m *testing.M) {
	loadEnv()
	listenPort()
	code := m.Run()
	os.Exit(code)
}

func loadEnv() {
	godotenv.Load("ssr.env")
}

func listenPort() {
	listen, err := net.Listen("tcp", ":12801")
	if err != nil {
		panic(err)
	}
	go func() {
		defer listen.Close()
		println("Listening on port 12801")
		for {
			conn, err := listen.Accept()
			if err != nil {
				panic(err)
			}
			go func() {
				defer conn.Close()
				buf := make([]byte, 1024)
				n, err := conn.Read(buf)
				if err != nil {
					panic(err)
				}
				println(string(buf[:n]))
			}()
		}
	}()
	// sleep 1 s
	time.Sleep(1 * time.Second)
}
