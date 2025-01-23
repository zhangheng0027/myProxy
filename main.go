package main

import "github.com/zhangheng0027/myProxy/net"

func main() {
	net.Configuration.ProxyServer =
		append(net.Configuration.ProxyServer,
			net.NewProxyServer("http://192.168.115.1:8080", 20*1024*1024, 20*1024*1024))

	net.Configuration.ProxyServer =
		append(net.Configuration.ProxyServer,
			net.NewProxyServer("http://192.168.115.2:8080", 20*1024*1024, 20*1024*1024))

	net.Configuration.ProxyServer =
		append(net.Configuration.ProxyServer,
			net.NewProxyServer("http://192.168.115.3:8080", 20*1024*1024, 20*1024*1024))

	net.Configuration.ProxyServer =
		append(net.Configuration.ProxyServer,
			net.NewProxyServer("http://192.168.115.4:8080", 20*1024*1024, 20*1024*1024))

	net.Configuration.ProxyServer =
		append(net.Configuration.ProxyServer,
			net.NewProxyServer("http://192.168.115.5:8080", 20*1024*1024, 20*1024*1024))

	net.Configuration.ProxyServer =
		append(net.Configuration.ProxyServer,
			net.NewProxyServer("http://192.168.115.6:8080", 20*1024*1024, 20*1024*1024))

	net.Configuration.ProxyServer =
		append(net.Configuration.ProxyServer,
			net.NewProxyServer("http://192.168.115.7:8080", 20*1024*1024, 20*1024*1024))

	net.Configuration.ProxyServer =
		append(net.Configuration.ProxyServer,
			net.NewProxyServer("http://192.168.115.8:8080", 20*1024*1024, 20*1024*1024))

	net.Configuration.ProxyServer =
		append(net.Configuration.ProxyServer,
			net.NewProxyServer("http://192.168.115.9:8080", 20*1024*1024, 20*1024*1024))

	net.Configuration.ProxyServer =
		append(net.Configuration.ProxyServer,
			net.NewProxyServer("http://192.168.115.10:8080", 20*1024*1024, 20*1024*1024))

	net.Configuration.RouteSSRUrl["*.google.com"] = true
	net.Configuration.RouteSSRUrl["*.googleapis.com"] = true
	net.Configuration.RouteSSRUrl["*.github.com"] = true
	net.Configuration.RouteSSRUrl["*.docker.com"] = true
	net.Configuration.RouteSSRUrl["*.docker.io"] = true
	net.Configuration.RouteSSRUrl["*.blizzard.com"] = true
	net.Configuration.RouteSSRUrl["*.google.com.hk"] = true
	net.Configuration.RouteSSRUrl["*.youtube.com"] = true
	net.Configuration.RouteSSRUrl["*.facebook.com"] = true
	net.Configuration.RouteSSRUrl["*.githubassets.com"] = true
	net.Configuration.RouteSSRUrl["*.githubusercontent.com"] = true
	net.Configuration.RouteSSRUrl["*.googleusercontent.io"] = true
	net.Configuration.RouteSSRUrl["*.wikipedia.org"] = true
	net.Configuration.RouteSSRUrl["*.quartz-scheduler.org"] = true
	net.Configuration.RouteSSRUrl["*.huggingface.co"] = true
	net.Configuration.RouteSSRUrl["*.hooos.com"] = true
	net.Configuration.RouteSSRUrl["*.cloudfront.net"] = true

	net.Configuration.WhiteListIp["localhost"] = true
	net.Listen(":15802")

}
