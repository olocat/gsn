package gsn

import "net"

type Behavior interface {
	OnConnect(net.Conn)
	OnPackage(net.Conn, []byte)
	OnRelease(net.Conn)
	OnClose(net.Conn)
	Convert(net.Conn) net.Conn
}

type BaseBehavior struct{}

func (s *BaseBehavior) OnConnect(net.Conn)             {}
func (s *BaseBehavior) OnPackage(net.Conn, []byte)     {}
func (s *BaseBehavior) OnRelease(net.Conn)             {}
func (s *BaseBehavior) OnClose(net.Conn)               {}
func (s *BaseBehavior) Convert(conn net.Conn) net.Conn { return conn }
