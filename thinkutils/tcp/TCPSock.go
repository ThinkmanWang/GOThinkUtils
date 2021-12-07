package thinktcp

import (
	"sync"
)

const (
	RecvBufLenMax = 4 * 1024
	SendBufLenMax = 4 * 1024

	SendBufCapMax = 10
	RecvBufCapMax = 10
)

type tcpSock struct {
	sendBufCap       uint32
	recvBufCap       uint32
	exitChan         chan struct{}
	waitGroup        *sync.WaitGroup
	onConnConnect    OnTcpConnCallback
	onConnClose      OnTcpConnCallback
	onCustomProtocol OnTcpCustomProtocol
}

type Protocol interface {
	Parse(b []byte, recvChan chan<- Packet)
	Process(conn *TcpConn, p Packet)
}

type Packet interface {
	Marshal() []byte
}
