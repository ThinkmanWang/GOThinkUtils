package thinkutils

import (
	"net"
)

type udputils struct {
}

func (this udputils) Send(szIP string, nPort int, data []byte) {
	go func() {
		ip := net.ParseIP(szIP)

		srcAddr := &net.UDPAddr{IP: []byte{0, 0, 0, 0}, Port: 0}
		dstAddr := &net.UDPAddr{IP: ip, Port: nPort}

		conn, err := net.DialUDP("udp", srcAddr, dstAddr)
		if err != nil {
			return
		}
		defer conn.Close()

		_, err = conn.Write(data)
		if err != nil {
			return
		}

		//log.Info("%d", nRet)
	}()

}

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
			go this.OnMsg(remoteAddr, data)
		}
		//log.Info("<%s> %s", remoteAddr, data[:n])
	}
}
