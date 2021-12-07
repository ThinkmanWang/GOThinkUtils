package main

import (
	"GOThinkUtils/tcp/protocol"
	"GOThinkUtils/thinkutils"
	"GOThinkUtils/thinkutils/logger"
	"fmt"
	. "github.com/ecofast/rtl/netutils"
	"github.com/ecofast/tcpsock"
)

var log *logger.LocalLogger = logger.DefaultLogger()

var (
	shutdown = make(chan bool, 1)

	listenPort int = 12345
)

func main() {
	fmt.Printf("tcpping listening on port: %d\n", listenPort)
	server := tcpsock.NewTcpServer(listenPort, 2, onConnConnect, onConnClose, onProtocol)
	log.Info("=====service start=====")
	go server.Serve()

	<-shutdown
	log.Info("shutdown server")
	server.Close()
	log.Info("=====service stop=====")
}

func onConnConnect(conn *tcpsock.TcpConn) {
	log.Info("accept connection from %s\n", IPFromNetAddr(conn.RawConn().RemoteAddr()))
}

func onConnClose(conn *tcpsock.TcpConn) {
	log.Info("connection closed from %s\n", IPFromNetAddr(conn.RawConn().RemoteAddr()))
}

func onMsg(conn *tcpsock.TcpConn, p *protocol.PingPacket) {
	szTxt := thinkutils.StringUtils.BytesToString(p.Body)
	log.Info("recved ping message from %s with %d bytes of data: %s\n", IPFromNetAddr(conn.RawConn().RemoteAddr()), p.BodyLen, szTxt)
	conn.Write(p)
}

func onProtocol() tcpsock.Protocol {
	proto := &protocol.PingProtocol{}
	proto.OnMessage(onMsg)
	return proto
}
