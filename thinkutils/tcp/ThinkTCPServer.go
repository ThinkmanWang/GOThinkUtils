package thinktcp

import (
	"GOThinkUtils/thinkutils"
	"GOThinkUtils/thinkutils/logger"
	"sync"
)

var (
	log *logger.LocalLogger = logger.DefaultLogger()
)

type ThinkTCPServer struct {
	OnMsgCallback     func(c *TcpConn, p *PingPacket)
	OnConnCallback    func(c *TcpConn)
	OnCloseCallback   func(c *TcpConn)
	OnTimeoutCallback func(c *TcpConn)
	Port              int
	HeartbeatTime     uint32

	m_pHeartbeatMgr *thinkutils.HeartbeatMgr
}

func (this *ThinkTCPServer) Serve() {
	this.m_pHeartbeatMgr = &thinkutils.HeartbeatMgr{}

	if (this.HeartbeatTime <= 0) {
		this.HeartbeatTime = 1800
	}
	this.m_pHeartbeatMgr.Init(this.HeartbeatTime, this.onTimeout)

	if this.Port <= 0 {
		this.Port = 8080
	}
	server := NewTcpServer(this.Port, 2, this.onConn, this.onClose, this.onProtocol)
	go server.Serve()

	defer server.Close()

	wg := &sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}

func (this *ThinkTCPServer) onMsg(c *TcpConn, p *PingPacket) {
	this.m_pHeartbeatMgr.Update(c)

	if this.OnMsgCallback != nil {
		go this.OnMsgCallback(c, p)
	}
}

func (this *ThinkTCPServer) onConn(c *TcpConn) {
	if this.OnConnCallback != nil {
		go this.OnConnCallback(c)
	}
}

func (this *ThinkTCPServer) onClose(c *TcpConn) {
	if this.OnCloseCallback != nil {
		go this.OnCloseCallback(c)
	}
}

func (this *ThinkTCPServer) onProtocol() Protocol {
	proto := &PingProtocol{}
	proto.OnMessage(this.onMsg)
	return proto
}

func (this *ThinkTCPServer) onTimeout(conn interface{}) {
	if this.OnTimeoutCallback != nil {
		go this.OnTimeoutCallback(conn.(*TcpConn))
	}
}
