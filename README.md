# Tunnel Monitor

Usage for monitoring iptables connection balance, and switching automatically.

## config file

	/etc/tunnelmonitor/config.ini  # basic config
	/etc/tunnelmonitor/tunnel.ini  # tunnel config

## compile

	cd /path/to/the/project
	go build

*Maybe you need use `go get` command to get the packge dependency.*

## usage

	cp -r config.example /etc/tunnelmonitor
	mkdir /path/to/log  #if you set 'log' parameter in config
	cd /path/to/the/binary
	sudo ./tunnelmonitor