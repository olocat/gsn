package gsn

import (
	"fmt"
	"gsn/util"
	"log"
	"net"
	"sync"
	"time"
)

const (
	ListenStateReady   = 0
	ListenStateRunning = 1
	ListenStateStop    = 2
)

type NetworkListen struct {
	Network string
	Address string

	OnConnect func(connId uint64, ctx *Context)
	OnClose   func(connId uint64)
	OnReceive func(connId uint64, ctx *Context)

	netListener net.Listener
	listenState uint8

	// conn deadline
	deadLine      time.Time
	writeDeadLine time.Time
	readDeadLine  time.Time

	rLock sync.RWMutex
}

func ListenTCP(address string) *NetworkListen {
	return &NetworkListen{
		Network: "tcp",
		Address: address,
	}
}

/* -- NetworkListen -- */

func (n *NetworkListen) Start() {
	n.SetState(ListenStateReady)
	if n.netListener != nil {
		_ = n.netListener.Close()
	}

	listener, err := net.Listen(n.Network, n.Address)
	if err != nil {
		fmt.Println("监听失败")
		log.Println(fmt.Sprintf(err.Error()))
		return
	}
	n.netListener = listener
	n.SetState(ListenStateRunning)
	n.accept()
}

func (n *NetworkListen) SetState(state uint8) {
	n.rLock.Lock()
	defer n.rLock.Unlock()
	n.listenState = state
}

func (n *NetworkListen) GetState() uint8 {
	n.rLock.RLock()
	defer n.rLock.RUnlock()
	return n.listenState
}

func (n *NetworkListen) accept() {
	defer n.netListener.Close()
	for {
		if n.GetState() == ListenStateStop {
			break
		}

		conn, err := n.netListener.Accept()
		if err != nil {
			fmt.Println("连接失败")
			log.Println(fmt.Sprintf(err.Error()))
			break
		}

		_ = conn.SetDeadline(n.deadLine)
		_ = conn.SetReadDeadline(n.readDeadLine)
		_ = conn.SetWriteDeadline(n.writeDeadLine)

		context := NewContext(conn)
		_, err = GlobalClientManager.AddContext(context)
		if err != nil {
			fmt.Println("增加连接失败")
			log.Println(fmt.Sprintf(err.Error()))
			continue
		}
	}
}

// 单独抽象出一个 packer 用于解析 bytes 数据
func (n *NetworkListen) receivePackListener(conn *Context) {
	fmt.Println("startConnListener")
	defer conn.Close()

	packLenBytes := make([]byte, PackHeadSize)
	for {

		size, err := conn.Read(packLenBytes)

		if err != nil {
			conn.Close()
		}

		packLen, err := util.BytesToUInt16(packLenBytes)
		if err != nil {
			fmt.Println("packLen is Error")
			continue
		}

		if packLen == 0 {
			fmt.Println("heartbeat")
			if GlobalConfig.HeartbeatMode {
				conn.HeartbeatChan <- true
			}
			continue
		}

		data := make([]byte, packLen-PackHeadSize)
		size, err = conn.Read(data)
		if err != nil {
			fmt.Println("Conn Error read data:", size)
			continue
		}

	}
}

/* ------------ */
