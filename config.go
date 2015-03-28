package main

import (
	"github.com/Unknwon/goconfig"
	"log"
	"net"
	"strconv"
	"time"
)

const (
	config_file        = "/etc/tunnelmonitor/config.ini"
	config_tunnel_file = "/etc/tunnelmonitor/tunnel.ini"
)

type level int

const (
	FATAL level = iota
	ERROR
	WARMING
	INFO
)

var (
	cfg               *goconfig.ConfigFile
	cfg_tunnel        *goconfig.ConfigFile
	chainName         string
	determineTime     int64
	detectingDuration int64
	rtt               int64
//	logLevel          int
)

func initConfig() {
	var err error
	cfg, err = goconfig.LoadConfigFile(config_file)
	if err != nil {
		log.Fatalf("CANNOT load config file(%s) : %s\n", config_file, err)
	}
	cfg_tunnel, err = goconfig.LoadConfigFile(config_tunnel_file)
	if err != nil {
		log.Fatalf("CANNOT load config file(%s) : %s\n", config_file, err)
	}
}
func readConfig() {
	chainName = configCommon("chainName", FATAL)

	determineTimeIntTmp, err := strconv.Atoi(configCommon("determineTime", INFO))
	determineTime = int64(determineTimeIntTmp)
	if err != nil {
		determineTime = 10
	}

	detectingDurationIntTmp, err := strconv.Atoi(configCommon("detectingDuration", INFO))
	detectingDuration = int64(detectingDurationIntTmp)
	if err != nil {
		detectingDuration = 10
	}

	rttIntTmp, err := strconv.Atoi(configCommon("rtt", INFO))
	rtt = int64(rttIntTmp)
	if err != nil {
		rtt = 1
	}

//	switch configCommon("logLevel", WARMING) {
//	case "FINEST":
//		logLevel = log4go.FINEST
//	case "FINE":
//		logLevel = log4go.FINE
//	case "DEBUG":
//		logLevel = log4go.DEBUG
//	case "TRACE":
//		logLevel = log4go.TRACE
//	case "INFO":
//		logLevel = log4go.INFO
//	case "WARMING":
//		logLevel = log4go.WARNING
//	case "ERROR":
//		logLevel = log4go.ERROR
//	case "CRITICAL":
//		logLevel = log4go.CRITICAL
//	default:
//		logLevel = log4go.INFO
//	}

	for _, tunnelName := range cfg_tunnel.GetSectionList() {
		tunnel[tunnelName] = &TunnelInfo{}
		ipAddr, err := net.ResolveIPAddr("ip4", configTunnel(tunnelName, "peerIP", FATAL))
		if err != nil {
			logger.Critical("peerIP format error for %s: %s ", tunnelName, err.Error())
		}
		tunnel[tunnelName].Ip = ipAddr.String()

		tunnel[tunnelName].Weight, err = strconv.Atoi(configTunnel(tunnelName, "weight", INFO))
		if err != nil {
			tunnel[tunnelName].Weight = 1
		}

		tunnel[tunnelName].Mark = configTunnel(tunnelName, "mark", FATAL)

		tunnel[tunnelName].RecoverCommand = configTunnel(tunnelName, "recoverCommand", INFO)

		tunnel[tunnelName].DownCommand = configTunnel(tunnelName, "downCommand", INFO)

		tunnel[tunnelName].lastlive, _ = time.Parse(time.ANSIC, time.ANSIC)
	}
}

func config(cfg *goconfig.ConfigFile, section string, key string, lvl level) string {
	value, err := cfg.GetValue(section, key)
	if err != nil {
		switch lvl {
		case FATAL:
			logger.Critical("Can't Read config %s", err.Error())
			panic(err.Error())
		case ERROR:
			logger.Error("Can't Read config %s", err.Error())
		case WARMING:
			logger.Warn("Can't Read config %s", err.Error())
		case INFO:
			logger.Info("Can't Read config %s", err.Error())
		}
		return ""
	}
	return value
}

func configCommon(key string, lvl level) string {
	return config(cfg, goconfig.DEFAULT_SECTION, key, lvl)
}

func configTunnel(tunnelName string, key string, lvl level) string {
	return config(cfg_tunnel, tunnelName, key, lvl)
}
