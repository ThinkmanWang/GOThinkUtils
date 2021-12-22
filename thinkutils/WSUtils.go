package thinkutils

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

type OnConnectCallback func(pConn *websocket.Conn)
type OnCloseCallback func(pConn *websocket.Conn)
type OnWSMsgCallback func(pConn *websocket.Conn, msg []byte)

var (
	upgrader = websocket.Upgrader{
		// 解决跨域问题
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	} // use default options

)

type WSHandler struct {
	OnConnect OnConnectCallback
	OnClose   OnCloseCallback
	OnMsg     OnWSMsgCallback
}

func (this *WSHandler) Handler(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Info("upgrade:", err)
		return
	}

	defer func() {
		if this.OnClose != nil {
			this.OnClose(ws)
		}
		ws.Close()
	}()

	if this.OnConnect != nil {
		go this.OnConnect(ws)
	}

	for {
		mt, message, err := ws.ReadMessage()
		if err != nil {
			break
		}

		switch mt {
		case websocket.BinaryMessage, websocket.TextMessage:
			if this.OnMsg != nil {
				go this.OnMsg(ws, message)
			}
			//go onMessage(c, message)
		default:
			continue
		}
	}
}

//func (this *WSHandler) Handler(w http.ResponseWriter, r *http.Request) {
//	//log.Info("%p", this)
//	c, err := wsUpgrader.Upgrade(w, r, nil)
//	if err != nil {
//		log.Info("upgrade:", err)
//		return
//	}
//

//
//	for {
//		mt, message, err := c.ReadMessage()
//		if err != nil {
//			break
//		}
//
//		switch mt {
//		case websocket.BinaryMessage, websocket.TextMessage:
//			if this.OnMsg != nil {
//				go this.OnMsg(c, message)
//			}
//			//go onMessage(c, message)
//		default:
//			continue
//		}
//	}
//}
