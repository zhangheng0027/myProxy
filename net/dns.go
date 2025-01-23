package net

import (
	"math/rand"
	"strings"
	"sync"
	"time"
)

type dns struct {
	host string
	port string

	ip    []string
	delay []time.Duration
	//sumDelay time.Duration
	count    []int
	allCount int

	// 权重
	weight []int

	defaultIndex int
	lastCheck    int64
	nextDnsTime  time.Time

	mu sync.Mutex
}

func newDns(host string, port string, ip []string) *dns {
	if ip == nil {
		ip = make([]string, 0, 5)
		ip = append(ip, host)
	} else {

		tip := make([]string, 0, len(ip))
		// 将 ipv4 加到 tip 中
		for _, s := range ip {
			if !strings.Contains(s, ":") {
				tip = append(tip, s)
			}
		}
		ip = append(tip, host)
	}
	d := &dns{
		host: host,
		port: port,

		ip:    ip,
		delay: make([]time.Duration, len(ip)),
		//sumDelay: 0,
		count:    make([]int, len(ip)),
		weight:   make([]int, 120),
		allCount: 0,

		defaultIndex: 0,
		lastCheck:    0,
		nextDnsTime:  time.Now().Add(time.Hour * 1),
	}

	for i := 0; i < len(ip); i++ {
		d.delay[i] = 0
		d.count[i] = 1
	}

	// 均匀分布
	for i := 0; i < len(d.weight); i++ {
		d.weight[i] = i % len(ip)
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

func (ns *dns) GetIp() (string, int) {

	intn := rand.Intn(len(ns.weight))

	return ns.ip[ns.weight[intn]], intn
}

func (ns *dns) GetRandIp() (string, int) {
	intn := rand.Intn(len(ns.ip))
	return ns.ip[intn], intn
}

func (ns *dns) AddOneDialed(ip string, index int, delay time.Duration) {
	go func() {
		ns.mu.Lock()
		defer ns.mu.Unlock()

		if ip != ns.ip[ns.weight[index]] {
			return
		}

		ns.delay[ns.weight[index]] = ns.delay[ns.weight[index]] + delay
		ns.count[ns.weight[index]]++
		//ns.sumDelay += delay
		ns.allCount++

		if (ns.allCount % (len(ns.ip) << 3)) == 0 {
			ns.allCount = 0
			// 计算平均延迟
			avgs := make([]int64, len(ns.delay))
			sumD := int64(0)

			minDelay := int64(10000000)

			for i := 0; i < len(ns.delay); i++ {
				avgs[i] = int64(ns.delay[i]) / int64(ns.count[i])
				if minDelay > avgs[i] {
					minDelay = avgs[i]
				}
				sumD += avgs[i]
				ns.delay[i] = 0
				ns.count[i] = 1
			}

			if minDelay == 0 || minDelay > int64(time.Millisecond)*50 {
				ns.nextDnsTime = time.Now()
			}

			wi := 0
			// 重新计算 ns.weight, 平均延迟越低, 权重越高
			for i := 0; i < len(avgs); i++ {
				c := int64(len(ns.weight)) * avgs[i] / sumD
				for ; wi < len(ns.weight) && c > 0; wi++ {
					c--
					ns.weight[wi] = i
				}

			}
		}

	}()
}

func (ns *dns) removeIp(ip string, index int) {
	if ip == ns.host {
		return
	}

	go func() {
		ns.mu.Lock()
		defer ns.mu.Unlock()
		if ns.ip[ns.weight[index]] != ip {
			return
		}

	}()
}
