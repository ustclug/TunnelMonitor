package main

import (
	"net"
	"os/exec"
	"strconv"
	"time"
)

func generator() {
	generatorLock.Lock()
	defer generatorLock.Unlock()
	logger.Info("updating new list...")
	for _, tunnel := range tunnels {
		ipt.ClearChain("mangle", tunnel.ChainName)
	}
	var weightSum, i int
	for _, tunnel := range tunnels {
		if tunnel.status {
			weightSum += tunnel.Weight
		}
	}
	for _, tunnel := range tunnels {
		if tunnel.status {
			for j := 0; j < tunnel.Weight; j++ {
				logger.Debug("iptables -t mangle -A %s -m statistic --mode nth --every %s --packet %s -j MARK --set-xmark %s", tunnel.ChainName, strconv.Itoa(weightSum), strconv.Itoa(i), tunnel.Mark)
				ipt.Append("mangle", tunnel.ChainName, "-m", "statistic", "--mode", "nth", "--every", strconv.Itoa(weightSum), "--packet", strconv.Itoa(i), "-j", "MARK", "--set-xmark", tunnel.Mark)
				i++
			}
		}
	}
	logger.Trace("updating finished")
}

func monitor() {
	pin.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		logger.Trace("ping %s %f ms", addr.String(), rtt.Seconds())
		tunInfo := tunnels[ip2tunnel[addr.String()]]
		tunInfo.lastlive = time.Now()
		if tunInfo.status == false {
			logger.Info("connection recover %s(%s)", ip2tunnel[addr.String()], tunInfo.Ip)
			tunInfo.status = true
			if tunInfo.RecoverCommand != "" {
				exec.Command("bash", "-c", tunInfo.RecoverCommand)
			}
			generator()
		}
	}
	pin.OnIdle = func() {}
	pin.MaxRTT = time.Second * time.Duration(rtt)
	go func() {
		for {
			pin.Run()
			time.Sleep(time.Second * time.Duration(rtt))
			<-pin.Done()
		}
	}()

	for {
		regenerate := false
		for tunName, tunInfo := range tunnels {
			if tunInfo.lastlive.Add(time.Second * time.Duration(determineTime)).Before(time.Now()) {
				if tunInfo.status {
					logger.Info("connection lost %s(%s)", tunName, tunInfo.Ip)
					tunInfo.status = false
					regenerate = true
					if tunInfo.DownCommand != "" {
						exec.Command("bash", "-c", tunInfo.DownCommand)
					}
				}
			}
		}
		if regenerate {
			generator()
		}
		time.Sleep(time.Second * time.Duration(detectingDuration))
	}
}
