package main

import (
	"github.com/ThinkmanWang/github.com/ThinkmanWang/GOThinkUtils/thinkutils"
	"github.com/ThinkmanWang/github.com/ThinkmanWang/GOThinkUtils/thinkutils/logger"
	"fmt"
	"net"
	"time"
)

var (
	log *logger.LocalLogger = logger.DefaultLogger()
)

func onMsg(pConn *net.UDPConn, addr net.Addr, data []byte) {
	log.Info("<%s> %s", addr.String(), thinkutils.StringUtils.BytesToString(data))


	message := []byte("Hello UDP client!")
	_, err := pConn.WriteToUDP(message, addr.(*net.UDPAddr))
	if err != nil {
		return
	}
}

func doSend() {
	serverAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:8083")
	if err != nil {
		return
	}

	localAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	if err != nil {
		return
	}

	pConn, err := net.DialUDP("udp", localAddr, serverAddr)
	if err != nil {
		return
	}

	defer pConn.Close()

	szMsg := fmt.Sprintf(`{"type": "main", "port": 1024}`)
	_, err = pConn.Write(thinkutils.StringUtils.StringToBytes(szMsg))
	if err != nil {
		return
	}

	buf := make([]byte, 1024)
	n, addr, err := pConn.ReadFromUDP(buf)
	fmt.Println("Received ", string(buf[0:n]), " from ", addr)
}

func sendTest() {
	go func() {
		for {
			go doSend()
			time.Sleep(1 * time.Second)
		}
	}()
}

func main() {
	log.Info("Start UDP Server.....")
	sendTest()
	pServer := &thinkutils.UDPServer{OnMsg: onMsg}
	pServer.Start(8083)

	//wg := sync.WaitGroup{}
	//wg.Add(1)
	//wg.Wait()
}
