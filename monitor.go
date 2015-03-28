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
	ipt.ClearChain("mangle", chainName)
	var weightSum, i int
	for _, t := range tunnel {
		if t.status {
			weightSum += t.Weight
		}
	}
	for _, t := range tunnel {
		if t.status {
			for j := 0; j < t.Weight; j++ {
				logger.Debug("iptables -t mangle -A %s -m statistic --mode nth --every %s --packet %s -j MARK --set-xmark %s", chainName, strconv.Itoa(weightSum), strconv.Itoa(i), t.Mark)
				ipt.Append("mangle", chainName, "-m", "statistic", "--mode", "nth", "--every", strconv.Itoa(weightSum), "--packet", strconv.Itoa(i), "-j", "MARK", "--set-xmark", t.Mark)
				i++
			}
		}
	}
	logger.Trace("updating finished")
}

func monitor() {
	pin.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		logger.Trace("ping %s %f ms", addr.String(), rtt.Seconds())
		tunInfo := tunnel[ip2tunnel[addr.String()]]
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
		for tunName, tunInfo := range tunnel {
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
