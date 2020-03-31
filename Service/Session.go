package Service

import (
	"github.com/yaice-rx/yaice/network"
	"sync"
)

type connSessionData struct {
	sync.Mutex
	conns map[uint64]network.IConn
}

func NewConnSessionMgr() *connSessionData {
	return &connSessionData{
		conns: map[uint64]network.IConn{},
	}
}

func (c *connSessionData) Add(guid uint64, conn network.IConn) {
	c.Lock()
	defer c.Unlock()
	c.conns[guid] = conn
}

func (c *connSessionData) Remove(connGuid uint64) {
	c.Lock()
	defer c.Unlock()
	delete(c.conns, connGuid)
}

func (c *connSessionData) GetAllList() map[uint64]network.IConn {
	return c.conns
}

func (c *connSessionData) Get(guid uint64) network.IConn {
	return c.conns[guid]
}
