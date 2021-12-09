package thinkutils

import (
	"net"
)

type OnUDPMsgCallback func(addr net.Addr, data []byte)
type UDPServer struct {
	OnMsg OnUDPMsgCallback
}

func (this *UDPServer) Start(nPort int) {
	this.StartEx(nPort, 1024)
}

func (this *UDPServer) StartEx(nPort int, bufSize uint32) {
	ip := net.ParseIP("0.0.0.0")
	listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: ip, Port: nPort})
	if err != nil {
		return
	}

	data := make([]byte, bufSize)
	for {
		_, remoteAddr, err := listener.ReadFrom(data)
		if err != nil {
			log.Info(err.Error())
			break
		}

		if nil != this.OnMsg {
			this.OnMsg(remoteAddr, data)
		}
		//log.Info("<%s> %s", remoteAddr, data[:n])
	}
}
