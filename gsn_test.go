package gsn

import (
	"encoding/json"
	"encoding/xml"
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
	strSize, err := unpacker.UnpackUint32()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	str, err := unpacker.UnpackString(int(strSize))
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

type TestJsonBehavior struct {
	BaseBehavior
}

func (s *TestJsonBehavior) OnPackage(conn net.Conn, stream []byte) {
	jsonData := map[string]any{}
	err := json.Unmarshal(stream, &jsonData)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for k, v := range jsonData {
		fmt.Println(k, v)
	}
}

func TestJsonData(t *testing.T) {
	listener, err := ListenTCP(":8088", &TestJsonBehavior{}, nil)
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

	jsonData := map[string]any{
		"boolData":    true,
		"numberData":  123456789,
		"stringData":  "abcdefghijklmnopqrstuvwxyz",
		"chineseData": "一地在要工上是中国同和的有人我主产不为这民了发以经",
		"mapData": map[string]any{
			"innerData": 123456,
		},
	}

	pack := NewJsonPacker()
	err = pack.PackObject(jsonData)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	conn.Write(pack.GetData())
	time.Sleep(3 * time.Second)
}

type TestXmlBehavior struct {
	BaseBehavior
}

type Address struct {
	City, State string
}
type Person struct {
	XMLName   xml.Name `xml:"person"`
	Id        int      `xml:"id,attr"`
	FirstName string   `xml:"name>first"`
	LastName  string   `xml:"name>last"`
	Age       int      `xml:"age"`
	Height    float32  `xml:"height,omitempty"`
	Married   bool
	Address
	Comment string `xml:",comment"`
}

type MyConn struct {
	net.Conn
	Id   uint32
	Name string
}

func (s *TestXmlBehavior) OnPackage(conn net.Conn, stream []byte) {
	p := Person{}
	err := xml.Unmarshal(stream, &p)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(p)
}

func TestXmlData(t *testing.T) {
	listener, err := ListenTCP(":8088", &TestXmlBehavior{}, nil)
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

	pack := NewXmlPacker()
	v := &Person{Id: 13, FirstName: "John", LastName: "Doe", Age: 42, Height: 172, Comment: " Need more details. "}
	v.Address = Address{"Hanga Roa", "Easter Island"}

	err = pack.PackObject(v)
	if err != nil {
		fmt.Println(err.Error())
	}

	conn.Write(pack.GetData())
	time.Sleep(3 * time.Second)
}
