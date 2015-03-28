package main

import (
	"github.com/coreos/go-iptables/iptables"
	"github.com/tatsushid/go-fastping"
	"sync"
	"time"
	"fmt"
)

type TunnelInfo struct {
	Ip             string
	Mark           string
	Weight         int
	RecoverCommand string
	DownCommand    string
	status         bool
	lastlive       time.Time
}

var (
	tunnel        map[string]*TunnelInfo
	generatorLock sync.Mutex
	ipt           *iptables.IPTables
	pin           *fastping.Pinger
	ip2tunnel     map[string]string
)

func init() {
	tunnel = make(map[string]*TunnelInfo)
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
	for name, t := range tunnel {
		ip2tunnel[t.Ip] = name
		pin.AddIP(t.Ip)
	}
	for n,t:=range tunnel{
		logger.Debug("[tunnel config]%s: %s",n,fmt.Sprint(*t))
	}
	ipt.NewChain("mangle",chainName)
	monitor()
}
