package net

import (
	"math/rand"
	"sync"
	"time"
)

type dns struct {
	host         string
	port         string
	ip           []string
	defaultIndex int
	lastCheck    int64
	nextDnsTime  time.Time

	mu sync.Mutex
}

func newDns(host string, port string, ip []string) *dns {
	d := &dns{
		host:         host,
		port:         port,
		ip:           ip,
		defaultIndex: 0,
		lastCheck:    0,
		nextDnsTime:  time.Now().Add(time.Hour * 24),
	}
	if ip == nil {
		d.ip = make([]string, 0, 5)
		d.ip = append(ip, host)
	}
	return d
}

// 是否进行dns查询
func (ns *dns) needDns() bool {
	if ns.nextDnsTime.Before(time.Now()) {
		return true
	}
	return false
}

func (ns *dns) GetIp() string {
	return ns.ip[ns.defaultIndex]
}

func (ns *dns) GetRandIp() string {
	return ns.ip[rand.Intn(len(ns.ip))]
}

func (ns *dns) removeIp(ip string) {
	if ip == ns.host {
		return
	}

	ns.mu.Lock()
	defer ns.mu.Unlock()
	if ip == ns.ip[ns.defaultIndex] {
		ns.defaultIndex++
	}
	if ns.defaultIndex >= len(ns.ip) {
		ns.defaultIndex = 0
		ns.ip[0] = ns.host
		ns.ip = ns.ip[:1]
	}

}
