package thinktcp

import (
	"encoding/binary"
)

const (
	PacketHeadSize = 4
)

type PingPacket struct {
	BodyLen uint32
	Body    []byte
}

func NewPingPacket(body []byte) *PingPacket {
	return &PingPacket{
		BodyLen: uint32(len(body)),
		Body:    body,
	}
}

func (self *PingPacket) Marshal() []byte {
	buf := make([]byte, PacketHeadSize+self.BodyLen)
	binary.LittleEndian.PutUint32(buf, self.BodyLen)
	copy(buf[PacketHeadSize:], self.Body[:])
	return buf
}

type PingProtocol struct {
	recvBuf    []byte
	recvBufLen int
	onMsg      func(c *TcpConn, p *PingPacket)
}

func (self *PingProtocol) Parse(b []byte, recvChan chan<- Packet) {
	count := len(b)
	if count+self.recvBufLen > RecvBufLenMax {
		return
	}

	self.recvBuf = append(self.recvBuf, b[0:count]...)
	self.recvBufLen += count
	offsize := 0
	offset := 0
	var pkt PingPacket
	for self.recvBufLen-offsize > PacketHeadSize {
		offset = 0
		pkt.BodyLen = binary.LittleEndian.Uint32(self.recvBuf[offsize+0 : offsize+PacketHeadSize])
		offset += PacketHeadSize
		pkglen := int(PacketHeadSize + pkt.BodyLen)
		if pkglen >= RecvBufLenMax {
			offsize = self.recvBufLen
			break
		}
		if offsize+pkglen > self.recvBufLen {
			break
		}

		recvChan <- NewPingPacket(self.recvBuf[offsize+offset : offsize+offset+int(pkt.BodyLen)])
		offsize += pkglen
	}

	self.recvBufLen -= offsize
	if self.recvBufLen > 0 {
		self.recvBuf = self.recvBuf[offsize : offsize+self.recvBufLen]
	} else {
		self.recvBuf = nil
	}
}

func (self *PingProtocol) Process(conn *TcpConn, p Packet) {
	packet := p.(*PingPacket)
	self.onMsg(conn, packet)
}

func (self *PingProtocol) OnMessage(fn func(c *TcpConn, p *PingPacket)) {
	self.onMsg = fn
}
