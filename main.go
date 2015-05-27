package main

import (
	"fmt"
	"github.com/coreos/go-iptables/iptables"
	"github.com/tatsushid/go-fastping"
	"sync"
	"time"
)

type TunnelInfo struct {
	Ip             string
	Mark           string
	Weight         int
	ChainName      string
	RecoverCommand string
	DownCommand    string
	status         bool
	lastlive       time.Time
}

var (
	tunnels       map[string]*TunnelInfo
	generatorLock sync.Mutex
	ipt           *iptables.IPTables
	pin           *fastping.Pinger
	ip2tunnel     map[string]string
)

func init() {
	tunnels = make(map[string]*TunnelInfo)
	ip2tunnel = make(map[string]string)
	ipt, _ = iptables.New()
	pin = fastping.NewPinger()
	initConfig()
	readConfig()
	initLogger()
}

func main() {
	defer logger.Close()
	logger.Info("System start")
	for name, t := range tunnels {
		ip2tunnel[t.Ip] = name
		pin.AddIP(t.Ip)
	}
	for n, t := range tunnels {
		logger.Debug("[tunnel config]%s: %s", n, fmt.Sprint(*t))
	}
	for _, tunnel := range tunnels {
		ipt.NewChain("mangle", tunnel.ChainName)
	}
	monitor()
}
