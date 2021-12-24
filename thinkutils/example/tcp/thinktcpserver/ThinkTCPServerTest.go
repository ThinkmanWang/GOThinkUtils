package main

import (
	"GOThinkUtils/thinkutils"
	"GOThinkUtils/thinkutils/logger"
	thinktcp "GOThinkUtils/thinkutils/tcp"
	"github.com/ecofast/rtl/netutils"
	"runtime"
)

var (
	log *logger.LocalLogger = logger.DefaultLogger()
)

func onTimeout(pConn *thinktcp.TcpConn) {
	log.Info("%p heartbeat timeout", pConn)

	if false == pConn.Closed() {
		pConn.Close()
	}
}

func onConn(conn *thinktcp.TcpConn) {
	log.Info("accept connection from %s\n", netutils.IPFromNetAddr(conn.RawConn().RemoteAddr()))
}

func onClose(conn *thinktcp.TcpConn) {
	log.Info("connection closed from %s\n", netutils.IPFromNetAddr(conn.RawConn().RemoteAddr()))
}

func onMsg(conn *thinktcp.TcpConn, p *thinktcp.PingPacket) {
	szTxt := thinkutils.StringUtils.BytesToString(p.Body)
	log.Info("recved ping message from %s with %d bytes of data: %s\n", netutils.IPFromNetAddr(conn.RawConn().RemoteAddr()), p.BodyLen, szTxt)
	conn.Write(p)
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	pServer := &thinktcp.ThinkTCPServer{
		OnConnCallback:    onConn,
		OnCloseCallback:   onClose,
		OnMsgCallback:     onMsg,
		OnTimeoutCallback: onTimeout,
		Port:              8000,
		HeartbeatTime:     10,
	}

	pServer.Serve()
}
