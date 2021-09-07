package gsn

import (
	"net"
)

type Context struct {
	Conn          net.Conn
	HeartbeatChan chan bool
}

func NewContext(conn net.Conn) *Context {
	return &Context{
		Conn:          conn,
		HeartbeatChan: make(chan bool),
	}
}

/* -- Context -- */

func (p *Context) Close() {
	_ = p.Conn.Close()
}

func (p *Context) RemoteAddr() string {
	return p.Conn.RemoteAddr().String()
}

func (p *Context) LocalAddr() string {
	return p.Conn.LocalAddr().String()
}

func (p *Context) Read(b []byte) (n int, err error) {
	return p.Conn.Read(b)
}

func (p *Context) Write(b []byte) (n int, err error) {
	return p.Conn.Write(b)
}

/* -------------------- */
