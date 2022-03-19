package gsn

import (
	"crypto/tls"
	"net"
)

func ListenTCP(addr string, behavior Behavior, tlsConfig *tls.Config) (*Listener, error) {
	var listener net.Listener
	var err error

	if tlsConfig == nil {
		listener, err = net.Listen("tcp", addr)
		if err != nil {
			return nil, err
		}

	} else {
		listener, err = tls.Listen("tcp", addr, tlsConfig)
		if err != nil {
			return nil, err
		}
	}

	return NewListener(listener, behavior), nil
}

func Dial(network, addr string, behavior Behavior, tlsConfig *tls.Config) (*Conn, error) {
	var conn net.Conn
	var err error

	if tlsConfig == nil {
		conn, err = net.Dial(network, addr)
		if err != nil {
			return nil, err
		}

	} else {
		conn, err = tls.Dial(network, addr, tlsConfig)
		if err != nil {
			return nil, err
		}
	}

	dialConn := NewConn(conn, behavior)
	dialConn.Start()
	return dialConn, nil
}
