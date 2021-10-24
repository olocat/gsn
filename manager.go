package gsn

import (
	"fmt"
	"sync"
	"time"
)

const (
	PackHeadSize = 4
)

var GlobalClientManager *ClientManager

func init() {
	GlobalClientManager = &ClientManager{connMap: map[uint32]*Context{}}
}

type ClientManager struct {
	connMap map[uint32]*Context

	rLock sync.RWMutex
}

/* -- ClientManager -- */

func (p *ClientManager) AddContext(ctx *Context) (uint32, error) {
	p.rLock.Lock()
	defer p.rLock.Unlock()
	connId, err := p.getConnId()
	if err != nil {
		return 0, err
	}
	p.connMap[connId] = ctx
	p.startConnListener(ctx)
	return connId, nil
}

func (p *ClientManager) RemoveContext(ctx *Context) {
	delete(p.connMap, ctx.ConnId)
}

func (p *ClientManager) getConnId() (uint32, error) {
	for i := 1; ; i++ {
		_, isExist := p.connMap[uint32(i)]
		if !isExist {
			return uint32(i), nil
		}
	}
}

func (p *ClientManager) startConnListener(conn *Context) {
	if GlobalConfig.HeartbeatMode {
		go p.heartbeatListener(conn)
	}
}

func (p *ClientManager) heartbeatListener(conn *Context) {
	if conn.HeartbeatChan == nil {
		conn.HeartbeatChan = make(chan bool)
	}
	var timeoutCount = 0
	for {
		select {
		case <-conn.HeartbeatChan:
			timeoutCount = 0
		case <-time.After(1 * time.Second):
			fmt.Println("timeout + 1")
			timeoutCount++
			if timeoutCount == 3 {
				fmt.Println("timeout heartbeat")
				conn.Close()
			}
		}

		if timeoutCount >= 3 {
			conn.Close()
			break
		}

		time.Sleep(1 * time.Second)
	}

}

/* -------------------- */
