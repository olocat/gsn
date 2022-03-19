package gsn

import (
	"crypto/tls"
	"errors"
	"math"
	"net"
	"sync"
)

const (
	ListenStateReady = 0
	ListenStateStop  = 1
)

type Listener struct {
	net.Listener
	TLSConfig *tls.Config
	Behavior  Behavior
	HeadSize  byte

	connMap map[uint32]net.Conn
	state   byte
	mpMux   sync.Mutex
	stMux   sync.RWMutex
}

func NewListener(listener net.Listener, behavior Behavior) *Listener {
	return &Listener{
		Listener: listener,
		HeadSize: DefaultHeadSize,
		Behavior: behavior,

		connMap: map[uint32]net.Conn{},
		state:   ListenStateReady,
	}
}

func (s *Listener) Start() {
	if s.GetState() != ListenStateReady {
		return
	}

	go s.serve()
}

func (s *Listener) Close() {
	s.SetState(ListenStateStop)
	_ = s.Listener.Close()
}

func (s *Listener) GetState() byte {
	s.stMux.RLock()
	defer s.stMux.RUnlock()

	return s.state
}

func (s *Listener) SetState(state byte) {
	s.stMux.Lock()
	defer s.stMux.Unlock()

	s.state = state
}

func (s *Listener) serve() {
	if s.Listener == nil {
		return
	}
	defer s.Close()

	for {
		if s.GetState() == ListenStateStop {
			break
		}

		conn, err := s.Accept()
		if err != nil {
			return
		}

		s.onAccept(conn)
	}
}

func (s *Listener) addConn(conn net.Conn) uint32 {
	s.mpMux.Lock()
	defer s.mpMux.Unlock()

	connId := s.genericConnId()
	s.connMap[connId] = conn
	return connId
}

func (s *Listener) closeConn(connId uint32) {
	conn, exist := s.connMap[connId]
	if !exist {
		return
	}

	if s.Behavior != nil {
		s.Behavior.OnRelease(conn)
	}
	_ = conn.Close()
	if s.Behavior != nil {
		s.Behavior.OnClose(conn)
	}
}

func (s *Listener) genericConnId() uint32 {
	for i := 1; i <= math.MaxUint32; i++ {
		if i == 0 {
			panic(errors.New("<gsn.Listener> too many conn"))
		}

		if _, exist := s.connMap[uint32(i)]; !exist {
			return uint32(i)
		}
	}

	return 0
}

func (s *Listener) onAccept(conn net.Conn) {
	if s.Behavior != nil {
		conn = s.Behavior.Convert(conn)
		s.Behavior.OnConnect(conn)
	}
	connId := s.addConn(conn)
	go s.startPackageListen(connId)
}

func (s *Listener) startPackageListen(connId uint32) {
	conn, exist := s.connMap[connId]
	if !exist {
		return
	}
	defer s.closeConn(connId)

	s.HeadSize = correctHeadSize(s.HeadSize)
	headSizeByte := make([]byte, s.HeadSize, s.HeadSize)

	for {
		if s.GetState() == ListenStateStop {
			break
		}

		err := readFull(conn, headSizeByte)
		if err != nil {
			break
		}

		headSize := transInt(headSizeByte) - uint64(s.HeadSize)
		if headSize <= 0 {
			continue
		}

		stream := make([]byte, headSize, headSize)
		err = readFull(conn, stream)
		if err != nil {
			break
		}

		if s.Behavior != nil {
			s.Behavior.OnPackage(conn, stream)
		}
	}
}
