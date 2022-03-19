package gsn

import (
	"fmt"
	"net"
	"testing"
	"time"
)

type TestListenBehavior struct {
	BaseBehavior
}

func (m *TestListenBehavior) OnConnect(conn net.Conn) {
	fmt.Println("ListenTCP OnConnect", conn.RemoteAddr())
	dataStr := "listen -> conn(含中文)"
	connPack := NewBinPacker()
	connPack.PackUint32(uint32(len(dataStr)))
	connPack.PackString(dataStr)
	_, err := conn.Write(connPack.GetData())
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func (m *TestListenBehavior) OnPackage(_ net.Conn, stream []byte) {
	unpacker := NewBinUnpacker(stream)
	strSize, err := unpacker.UnPackUint32()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	str, err := unpacker.UnPackString(int(strSize))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Receive", str)
}

type TestConnBehavior struct {
	TestListenBehavior
}

func TestNetwork(t *testing.T) {
	listener, err := ListenTCP(":8088", &TestListenBehavior{}, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	listener.Start()

	time.Sleep(time.Second)

	conn, err := Dial("tcp", ":8088", &TestConnBehavior{}, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	dataStr := "conn -> listen(含中文)"
	connPack := NewBinPacker()
	connPack.PackUint32(uint32(len(dataStr)))
	connPack.PackString(dataStr)
	_, err = conn.Write(connPack.GetData())
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	time.Sleep(2 * time.Second)
}
