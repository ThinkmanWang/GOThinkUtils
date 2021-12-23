package thinkutils

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

type OnConnectCallback func(pConn *websocket.Conn)
type OnCloseCallback func(pConn *websocket.Conn)
type OnTimeoutCallback func(pConn *websocket.Conn)
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
	OnTimeout OnTimeoutCallback

	m_pHeartbeatMgr  *HeartbeatMgr
	HeartbeatTimeout uint32
}

func (this *WSHandler) Init() {
	this.m_pHeartbeatMgr = &HeartbeatMgr{}
	if this.HeartbeatTimeout <= 0 {
		this.HeartbeatTimeout = 60
	}

	this.m_pHeartbeatMgr.Init(this.HeartbeatTimeout, this.onHBTimeout)
}

func (this *WSHandler) onHBTimeout(conn interface{}) {
	pConn := conn.(*websocket.Conn)
	log.Info("%p heartbeat timeout", pConn)

	if this.OnTimeout != nil {
		go this.OnTimeout(pConn)
	}

	pConn.Close()
}

func (this *WSHandler) Handler(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Info("upgrade:", err)
		return
	}

	defer func() {
		if this.OnClose != nil {
			go this.OnClose(ws)
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

		this.m_pHeartbeatMgr.Update(ws)
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
