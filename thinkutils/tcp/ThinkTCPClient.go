package thinktcp

import (
	"github.com/ThinkmanWang/github.com/ThinkmanWang/GOThinkUtils/thinkutils"
	"fmt"
)

type ThinkTCPClient struct {
	OnMsgCallback   func(c *TcpConn, p *PingPacket)
	OnConnCallback  func(c *TcpConn)
	OnCloseCallback func(c *TcpConn)
	Host            string
	Port            uint32

	m_pClient *TcpClient
	m_pConn   *TcpConn
}

func (this *ThinkTCPClient) onMsg(c *TcpConn, p *PingPacket) {
	if this.OnMsgCallback != nil {
		go this.OnMsgCallback(c, p)
	}
}

func (this *ThinkTCPClient) onConn(c *TcpConn) {
	this.m_pConn = c
	if this.OnConnCallback != nil {
		go this.OnConnCallback(c)
	}
}

func (this *ThinkTCPClient) onClose(c *TcpConn) {
	if this.OnCloseCallback != nil {
		go this.OnCloseCallback(c)
	}
}

func (this *ThinkTCPClient) onProtocol() Protocol {
	proto := &PingProtocol{}
	proto.OnMessage(this.onMsg)
	return proto
}

func (this *ThinkTCPClient) Close() {
	if nil == this.m_pClient {
		return
	}

	this.m_pClient.Close()
}

func (this *ThinkTCPClient) Send(data []byte) {
	p := NewPingPacket(data)
	this.m_pConn.Write(p)
}

func (this *ThinkTCPClient) SendString(data string) {
	this.Send(thinkutils.StringUtils.StringToBytes(data))
}

func (this *ThinkTCPClient) Connect() {
	szAddr := fmt.Sprintf("%s:%d", this.Host, this.Port)
	this.m_pClient = NewTcpClient(szAddr, this.onConn, this.onClose, this.onProtocol)
	go this.m_pClient.Run()
}
