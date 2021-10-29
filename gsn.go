package gsn

import (
	"encoding/binary"
	"fmt"
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

	OnConnect func(ctx *Context)
	OnClose   func(connId uint32)
	OnReceive func(pack *ReceivePack)

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
	if n.netListener != nil {
		_ = n.netListener.Close()
	}

	listener, err := net.Listen(n.Network, n.Address)
	if err != nil {
		fmt.Println("listen network fault")
		log.Println(fmt.Sprintf(err.Error()))
		return
	}
	n.netListener = listener
	n.SetState(ListenStateRunning)
	n.accept()
}

func (n *NetworkListen) Close() {
	err := n.netListener.Close()
	if err != nil {
		panic(err)
	}
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
	defer n.Close()
	for {
		if n.GetState() == ListenStateStop {
			break
		}

		conn, err := n.netListener.Accept()
		if err != nil {
			fmt.Println("conn accept fault")
			log.Println(fmt.Sprintf(err.Error()))
			break
		}

		_ = conn.SetDeadline(n.deadLine)
		_ = conn.SetReadDeadline(n.readDeadLine)
		_ = conn.SetWriteDeadline(n.writeDeadLine)

		context, err := GlobalClientManager.ManageConn(conn)
		if err != nil {
			fmt.Println("client pool add conn fault")
			log.Println(fmt.Sprintf(err.Error()))
			continue
		}

		if n.OnConnect != nil {
			n.OnConnect(context)
		}

		go n.receivePackListener(context)
	}
}

func (n *NetworkListen) receivePackListener(conn *Context) {
	fmt.Println("startConnListener")
	defer conn.Close()

	packLenBytes := make([]byte, PackHeadSize)
	for {

		size, err := conn.Read(packLenBytes)

		if err != nil {
			conn.Close()
		}

		if size < PackHeadSize {
			continue
		}

		packLen := binary.BigEndian.Uint32(packLenBytes)

		if packLen == 0 {
			fmt.Println("heartbeat")
			if GlobalConfig.HeartbeatMode {
				conn.HeartbeatChan <- true
			}
			continue
		}

		data := make([]byte, packLen-PackHeadSize)
		_, err = conn.Read(data)
		if err != nil {
			continue
		}

		receivePack := NewReceivePack(conn, data)
		if n.OnReceive != nil {
			n.OnReceive(receivePack)
		}

	}
}

/* ------------ */
