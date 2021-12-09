package main

import (
	"GOThinkUtils/thinkutils"
	"GOThinkUtils/thinkutils/logger"
	"net"
)

var (
	log *logger.LocalLogger = logger.DefaultLogger()
)

func onMsg(addr net.Addr, data []byte) {
	log.Info("<%s> %s", addr, thinkutils.StringUtils.BytesToString(data))
}

func main() {
	log.Info("Start UDP Server.....")
	pServer := &thinkutils.UDPServer{OnMsg: onMsg}
	pServer.Start(8083)
}
