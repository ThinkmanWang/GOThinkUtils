package main

import (
	"GOThinkUtils/tcp/protocol"
	"GOThinkUtils/thinkutils/logger"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	. "github.com/ecofast/rtl/netutils"
	"github.com/ecofast/tcpsock"
)

var log *logger.LocalLogger = logger.DefaultLogger()

var (
	shutdown = make(chan bool, 1)

	listenPort int = 12345
)

func init() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signals
		shutdown <- true
	}()
}

func parseFlag() {
	flag.IntVar(&listenPort, "p", listenPort, "listen port")
	flag.Parse()
}

func main() {
	parseFlag()

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
	log.Info("recved ping message from %s with %d bytes of data\n", IPFromNetAddr(conn.RawConn().RemoteAddr()), protocol.PacketHeadSize+p.BodyLen)
	conn.Write(p)
}

func onProtocol() tcpsock.Protocol {
	proto := &protocol.PingProtocol{}
	proto.OnMessage(onMsg)
	return proto
}
