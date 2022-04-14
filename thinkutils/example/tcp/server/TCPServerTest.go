package main

import (
	"github.com/ThinkmanWang/github.com/ThinkmanWang/GOThinkUtils/thinkutils"
	"github.com/ThinkmanWang/github.com/ThinkmanWang/GOThinkUtils/thinkutils/logger"
	"github.com/ThinkmanWang/github.com/ThinkmanWang/GOThinkUtils/thinkutils/tcp"
	"fmt"
	"github.com/ecofast/rtl/netutils"
)

var log *logger.LocalLogger = logger.DefaultLogger()

var (
	shutdown = make(chan bool, 1)

	listenPort      int = 12345
	m_pHeartbeatMgr *thinkutils.HeartbeatMgr
)

func onHBTimeout(conn interface{}) {
	pConn := conn.(*thinktcp.TcpConn)
	log.Info("%p heartbeat timeout", pConn)

	if false == pConn.Closed() {
		pConn.Close()
	}
}

func onConnConnect(conn *thinktcp.TcpConn) {
	log.Info("accept connection from %s\n", netutils.IPFromNetAddr(conn.RawConn().RemoteAddr()))
}

func onConnClose(conn *thinktcp.TcpConn) {
	log.Info("connection closed from %s\n", netutils.IPFromNetAddr(conn.RawConn().RemoteAddr()))
}

func onMsg(conn *thinktcp.TcpConn, p *thinktcp.PingPacket) {
	m_pHeartbeatMgr.Update(conn)

	szTxt := thinkutils.StringUtils.BytesToString(p.Body)
	log.Info("recved ping message from %s with %d bytes of data: %s\n", netutils.IPFromNetAddr(conn.RawConn().RemoteAddr()), p.BodyLen, szTxt)
	conn.Write(p)
}

func onProtocol() thinktcp.Protocol {
	proto := &thinktcp.PingProtocol{}
	proto.OnMessage(onMsg)
	return proto
}

func main() {

	m_pHeartbeatMgr = &thinkutils.HeartbeatMgr{}
	m_pHeartbeatMgr.Init(10, onHBTimeout)

	fmt.Printf("tcpping listening on port: %d\n", listenPort)
	server := thinktcp.NewTcpServer(listenPort, 2, onConnConnect, onConnClose, onProtocol)
	log.Info("=====service start=====")
	go server.Serve()

	<-shutdown
	log.Info("shutdown server")
	server.Close()
	log.Info("=====service stop=====")
}
