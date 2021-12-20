package main

import (
	"GOThinkUtils/thinkutils"
	"GOThinkUtils/thinkutils/logger"
	"net"
	"time"
)

var (
	log *logger.LocalLogger = logger.DefaultLogger()
)

func onMsg(addr net.Addr, data []byte) {
	log.Info("<%s> %s", addr.String(), thinkutils.StringUtils.BytesToString(data))
}

func sendTest() {
	go func() {
		for {
			thinkutils.UDPUtils.Send("127.0.0.1", 8083, []byte(thinkutils.DateTime.CurDatetime()))
			time.Sleep(1 * time.Second)
		}
	}()
}

func main() {
	log.Info("Start UDP Server.....")
	sendTest()
	pServer := &thinkutils.UDPServer{OnMsg: onMsg}
	pServer.Start(8083)
}
