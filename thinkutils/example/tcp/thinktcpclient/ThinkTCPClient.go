package main

import (
	"github.com/ThinkmanWang/GOThinkUtils/thinkutils"
	"github.com/ThinkmanWang/GOThinkUtils/thinkutils/logger"
	thinktcp "github.com/ThinkmanWang/GOThinkUtils/thinkutils/tcp"
	"github.com/ecofast/rtl/netutils"
	"runtime"
	"sync"
	"time"
)

var (
	log       *logger.LocalLogger = logger.DefaultLogger()
	g_pClient *thinktcp.ThinkTCPClient
)

func onConnect(c *thinktcp.TcpConn) {
	log.Info("successfully connect to server", netutils.IPFromNetAddr(c.RawConn().RemoteAddr()))
	go doSend()
}

func onClose(c *thinktcp.TcpConn) {
	log.Info("disconnect from server", netutils.IPFromNetAddr(c.RawConn().RemoteAddr()))
}

func onMsg(c *thinktcp.TcpConn, p *thinktcp.PingPacket) {
	log.Info("%d bytes from %s: %s", p.BodyLen, netutils.IPFromNetAddr(c.RawConn().RemoteAddr()), thinkutils.StringUtils.BytesToString(p.Body))
}

func doSend() {
	for {
		g_pClient.SendString(thinkutils.DateTime.CurDatetime())
		time.Sleep(time.Second)
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	g_pClient = &thinktcp.ThinkTCPClient{
		OnConnCallback:  onConnect,
		OnCloseCallback: onClose,
		OnMsgCallback:   onMsg,
		Host:            "127.0.0.1",
		Port:            8000,
	}

	g_pClient.Connect()

	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}
