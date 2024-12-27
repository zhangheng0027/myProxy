package test

import (
	"encoding/base64"
	"fmt"
	"github.com/joho/godotenv"
	"net/url"
	"os"
	"testing"
)

func TestBase64(t *testing.T) {
	encoded := "aGVsbG8gd29ybGQ=" // This is "hello world" encoded in base64
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		fmt.Println("Error decoding base64:", err)
		return
	}
	fmt.Println("Decoded string:", string(decoded))
}

func TestCompare(t *testing.T) {
	surl, b := os.LookupEnv("ssrurl")

	if b {
		fmt.Println("ssrurl:", surl)
	}

	// 分割 surl，获取 // 后面的内容
	u, err := url.Parse(surl)
	if err != nil {
		t.Fatal(err)
	}
	println(u.Host)
	// base64 解码 u.host
	//decodeString, err := base64.URLEncoding.DecodeString(u.Host)

	decodeString, err := base64.NewEncoding("").DecodeString(surl)
	if err != nil {
		t.Fatal(err)
	}

	println(decodeString)

}

func TestMain(m *testing.M) {
	loadEnv()
	code := m.Run()
	os.Exit(code)
}

func loadEnv() {
	godotenv.Load("ssr.env")
}
