package thinkutils

import (
	"github.com/gorilla/websocket"
	"net/http"
)

type OnConnectCallback func(pConn *websocket.Conn)
type OnCloseCallback func(pConn *websocket.Conn)
type OnWSMsgCallback func(pConn *websocket.Conn, msg []byte)

var (
	wsUpgrader = websocket.Upgrader{}
)

type WSHandler struct {
	OnConnect OnConnectCallback
	OnClose   OnCloseCallback
	OnMsg     OnWSMsgCallback
}

func (this WSHandler) Handler(w http.ResponseWriter, r *http.Request) {
	c, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Info("upgrade:", err)
		return
	}

	defer func() {
		if this.OnClose != nil {
			this.OnClose(c)
		}
		c.Close()
	}()

	if this.OnConnect != nil {
		go this.OnConnect(c)
	}

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			break
		}

		switch mt {
		case websocket.BinaryMessage, websocket.TextMessage:
			if this.OnMsg != nil {
				go this.OnMsg(c, message)
			}
			//go onMessage(c, message)
		default:
			continue
		}
	}
}
