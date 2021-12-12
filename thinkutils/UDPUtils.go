package thinkutils

import (
	"fmt"
	"net"
	"sync"
)

type udputils struct {
}

var (
	g_lockUDPClient sync.Mutex
	g_mapUDPClient  map[string]*net.UDPConn
)

func (this udputils) makeUDPClient(szIP string, nPort int) *net.UDPConn {
	defer g_lockUDPClient.Unlock()
	g_lockUDPClient.Lock()

	szConn := fmt.Sprintf("%s:%d", szIP, nPort)
	pConn := g_mapUDPClient[szConn]
	if nil == pConn {
		ip := net.ParseIP(szIP)

		srcAddr := &net.UDPAddr{IP: []byte{0, 0, 0, 0}, Port: 0}
		dstAddr := &net.UDPAddr{IP: ip, Port: nPort}

		conn, err := net.DialUDP("udp", srcAddr, dstAddr)
		if err != nil {
			return nil
		}

		pConn = conn
		g_mapUDPClient[szConn] = pConn
	}

	return pConn
}

func (this udputils) Send(szIP string, nPort int, data []byte) {
	go func() {
		if nil == g_mapUDPClient {
			g_mapUDPClient = make(map[string]*net.UDPConn)
		}

		szConn := fmt.Sprintf("%s:%d", szIP, nPort)
		pConn := g_mapUDPClient[szConn]
		if nil == pConn {
			pConn = this.makeUDPClient(szIP, nPort)
		}

		if nil == pConn {
			return
		}

		log.Info("%p", pConn)
		_, err := pConn.Write(data)
		if err != nil {
			pConn.Close()
			delete(g_mapUDPClient, szConn)
			return
		}
	}()

}

type OnUDPMsgCallback func(addr net.Addr, data []byte)
type UDPServer struct {
	OnMsg OnUDPMsgCallback
}

func (this *UDPServer) Start(nPort int) {
	this.StartEx(nPort, 8)
}

func (this *UDPServer) StartEx(nPort int, bufSize uint32) {
	ip := net.ParseIP("0.0.0.0")
	listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: ip, Port: nPort})
	if err != nil {
		return
	}

	for {
		buf := make([]byte, bufSize)
		_, remoteAddr, err := listener.ReadFrom(buf)
		if err != nil {
			log.Info(err.Error())
			break
		}

		if nil != this.OnMsg {
			go this.OnMsg(remoteAddr, buf)
		}
		//log.Info("<%s> %s", remoteAddr, data[:n])
	}
}
