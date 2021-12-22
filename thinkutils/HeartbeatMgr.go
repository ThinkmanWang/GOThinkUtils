package thinkutils

import (
	"github.com/emirpasic/gods/maps/hashmap"
	"time"
)

type OnHBTimeoutCallback func(pConn interface{})

type ConnNode struct {
	Conn      interface{}
	Timestamp int64
}

type HeartbeatMgr struct {
	m_cbTimeout OnHBTimeoutCallback
	m_nTimeout  uint32
	m_pConns    *hashmap.Map
}

func (this *HeartbeatMgr) start() {
	for {
		time.Sleep(3 * time.Second)

		//log.Info("Start Heartbeat check for %d conns", this.m_pConns.Size())
		nTimestamp := DateTime.Timestamp()
		for _, node := range this.m_pConns.Values() {
			if nTimestamp-node.(*ConnNode).Timestamp > int64(this.m_nTimeout) {
				log.Info("Remove conn %p", node.(*ConnNode).Conn)

				this.m_pConns.Remove(node.(*ConnNode).Conn)
				if this.m_cbTimeout != nil {
					go this.m_cbTimeout(node.(*ConnNode).Conn)
				}
			}
		}
	}
}

func (this *HeartbeatMgr) Init(nTimeout uint32, callback OnHBTimeoutCallback) {
	if nil == this.m_pConns {
		this.m_pConns = hashmap.New()
	}

	this.m_nTimeout = nTimeout
	this.m_cbTimeout = callback

	go this.start()
}

func (this *HeartbeatMgr) Update(pConn interface{}) {
	var pNode *ConnNode

	if _pNode, bFound := this.m_pConns.Get(pConn); bFound {
		pNode = _pNode.(*ConnNode)
	} else {
		pNode = &ConnNode{
			Conn: pConn,
		}

		this.m_pConns.Put(pConn, pNode)
	}

	pNode.Timestamp = DateTime.Timestamp()
}

func (this *HeartbeatMgr) Remove(pConn interface{}) {
	this.m_pConns.Remove(pConn)
}

func (this *HeartbeatMgr) Count() int {
	return this.m_pConns.Size()
}
