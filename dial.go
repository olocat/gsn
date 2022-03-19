package gsn

import (
	"crypto/tls"
	"net"
	"sync"
)

const (
	ConnStateReady = 0
	ConnStateStop  = 1
)

type Conn struct {
	net.Conn
	TLSConfig *tls.Config
	Behavior  Behavior
	HeadSize  byte

	state byte
	mpMux sync.Mutex
	stMux sync.RWMutex
}

func NewConn(conn net.Conn, behavior Behavior) *Conn {
	return &Conn{
		Conn:     conn,
		Behavior: behavior,
		HeadSize: DefaultHeadSize,
	}
}

func (s *Conn) Start() {
	if s.Conn == nil {
		return
	}

	if s.GetState() != ConnStateReady {
		return
	}

	if s.Behavior != nil {
		s.Conn = s.Behavior.Convert(s.Conn)
	}

	go s.listen()
}

func (s *Conn) Close() {
	s.SetState(ConnStateStop)
	if s.Conn == nil {
		return
	}

	if s.Behavior != nil {
		s.Behavior.OnRelease(s.Conn)
	}
	_ = s.Conn.Close()
	if s.Behavior != nil {
		s.Behavior.OnClose(s.Conn)
	}
}

func (s *Conn) GetState() byte {
	s.stMux.RLock()
	defer s.stMux.RUnlock()

	return s.state
}

func (s *Conn) SetState(state byte) {
	s.stMux.Lock()
	defer s.stMux.Unlock()

	s.state = state
}

func (s *Conn) listen() {
	defer s.Close()
	s.HeadSize = correctHeadSize(s.HeadSize)
	headSizeByte := make([]byte, s.HeadSize, s.HeadSize)

	for {
		if s.GetState() == ConnStateStop {
			break
		}

		err := readFull(s.Conn, headSizeByte)
		if err != nil {
			break
		}

		headSize := transInt(headSizeByte) - uint64(s.HeadSize)
		if headSize <= 0 {
			continue
		}

		stream := make([]byte, headSize, headSize)
		err = readFull(s.Conn, stream)
		if err != nil {
			break
		}

		if s.Behavior != nil {
			s.Behavior.OnPackage(s.Conn, stream)
		}
	}
}
