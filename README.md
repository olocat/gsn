# gsn

一个简单的Go网络工具，可以用于TCP或UDP。

支持 TCP UDP

## 描述

- 采用模块化设计，如同积木一样，可以自由控制数据格式
- 代码简单，使用方便
- 无外部依赖

- [开发进度](https://github.com/olocat/gsn/blob/master/WorkProgress.md)
- [开发文档](https://github.com/olocat/gsn/blob/master/doc.md)

## 开始

`go get github.com/olocat/gsn`

## 示例

### 1.启动一个TCP监听

```go
package main

import (
	"fmt"
	"github.com/olocat/gsn"
)

func main() {
	listener, err := gsn.ListenTCP(":8088", nil, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	listener.Start()
}
```

上面的代码虽然启动了一个监听，但并没有处理。

在 gsn 中，对行为的处理由 `Behavior` 接口来控制。

```go
type Behavior interface {
OnConnect(net.Conn)
OnPackage(net.Conn, []byte)
OnRelease(net.Conn)
OnClose(net.Conn)
Convert(net.Conn) net.Conn
}
```

OnConnect 会在连接成功后执行 OnPackage 会在收到数据时执行 OnRelease 会在关闭连接前执行 OnClose 会在关闭连接后执行 Convert 用于设置连接参数，或用于将net.Conn转成自定义的连接类型

现有一个已经空实现 Behavior 的结构体 BaseBehavior

```go
type BaseBehavior struct{}

func (s *BaseBehavior) OnConnect(net.Conn)             {}
func (s *BaseBehavior) OnPackage(net.Conn, []byte)     {}
func (s *BaseBehavior) OnRelease(net.Conn)             {}
func (s *BaseBehavior) OnClose(net.Conn)               {}
func (s *BaseBehavior) Convert(conn net.Conn) net.Conn { return conn }
```

如果自定的 Behavior 不需要实现全部的 Behavior 接口内的方法，可以选择直接继承 BaseBehavior。

### 2.启动带有行为的监听

```go
package main

import (
	"fmt"
	"github.com/olocat/gsn"
	"net"
)

type MyBehavior struct {
	gsn.BaseBehavior
}

func OnConnect(conn net.Conn) {
	fmt.Println("有一个连接进来了！", conn.RemoteAddr())
}

func OnPackage(conn net.Conn, stream []byte) {
	fmt.Println("收到数据包", stream)
}

func main() {
	listener, err := gsn.ListenTCP(":8088", nil, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	listener.Start()
}
```

如上述代码中，当有连接进来时，OnConnect 就会执行 当收到数据包后，OnPackage 就会执行

### 3.发起连接

使用 Dial 可以快捷的发起一个连接

```go
package main

import (
	"fmt"
	"github.com/olocat/gsn"
	"net"
)

type MyBehavior struct {
	gsn.BaseBehavior
}

func OnPackage(conn net.Conn, stream []byte) {
	fmt.Println("收到数据包", stream)
}

func main() {
	conn, err := gsn.Dial("tcp", ":8088", &MyBehavior{}, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//TODO
}
```

Dial 支持 TCP 和 UDP 也同样使用可以使用 Behavior 管理行为

### 4.发包与解包

gsn 默认提供三种收发包格式 bin,json,xml

- bin 即手动控制 stream 中每个字节表示的数据类型
- json 即使用 json 格式
- xml 同样，即使用 xml 格式

原本有想过支持 protobuf 的，但为了不引入每三方库，只得放弃。 但 gsn 也可以很方便的引入 protobuf，如有需要可以自行实现。

发包：

```go
package main

import (
	"fmt"
	"github.com/olocat/gsn"
)

// Packer 发包示例
func Packer() {
	packer := gsn.NewBinPacker()

	packer.PackInt16(-32768)
	packer.PackUint32(4294967295)
	packer.PackUint64(18446744073709551615)
	packer.PackBytes([]byte{0, 1, 2, 3, 255})

	str := "一地在要工上是中国同和的有人我主产不为这民了发以经"
	packer.PackUint32(uint32(len(str)))
	packer.PackString(str)

	packer.GetData() // 使用 GetData 获取封装好的数据 []byte 类型
}

// Unpacker 解包示例
func Unpacker(stream []byte) {
	unpacker := gsn.NewBinUnpacker(stream)

	v1, _ := unpacker.UnpackInt16()  // -32768
	v2, _ := unpacker.UnpackUint32() // 4294967295
	v3, _ := unpacker.UnpackUint64() // 18446744073709551615
	v4, _ := unpacker.UnpackBytes(5) // 5为 byte 数据长度

	sz, _ := unpacker.UnpackUint32() // 字符串长度
	str, _ := unpacker.UnpackString(int(sz))

	fmt.Println(v1, v2, v3, v4)
	fmt.Println(str)
	//TODO
}
```

除此外还可以使用 json 或 xml

json 示例

```go
package main

import (
	"encoding/json"
	"fmt"
	"github.com/olocat/gsn"
)

func Packer() {
	jsonData := map[string]any{
		"boolData":    true,
		"numberData":  123456789,
		"stringData":  "abcdefghijklmnopqrstuvwxyz",
		"chineseData": "一地在要工上是中国同和的有人我主产不为这民了发以经",
		"mapData": map[string]any{
			"innerData": 123456,
		},
	}

	packer := gsn.NewJsonPacker()
	err := packer.PackObject(jsonData)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	packer.GetData() // 使用 GetData 获取封装好的数据 []byte 类型
}

func Unpacker(stream []byte) {
	jsonData := map[string]any{}
	err := json.Unmarshal(stream, &jsonData)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for k, v := range jsonData {
		fmt.Println(k, v)
	}
	//TODO
}
```

xml 示例

```go
package main

import (
	"encoding/xml"
	"fmt"
	"github.com/olocat/gsn"
)

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

func Packer() {
	pack := gsn.NewXmlPacker()
	v := &Person{Id: 13, FirstName: "John", LastName: "Doe", Age: 42, Height: 172, Comment: " Need more details. "}
	v.Address = Address{"Hanga Roa", "Easter Island"}

	err := pack.PackObject(v)
	if err != nil {
		fmt.Println(err.Error())
	}

	pack.GetData() // 使用 GetData 获取封装好的数据 []byte 类型
}

func Unpacker(stream []byte) {
	p := Person{}
	err := xml.Unmarshal(stream, &p)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(p)
	//TODO
}
```

### 5.收发包示例

```go
package main

import (
	"encoding/json"
	"fmt"
	"github.com/olocat/gsn"
	"net"
	"time"
)

type MyBehavior struct {
	gsn.BaseBehavior
}

func (m *MyBehavior) OnPackage(conn net.Conn, stream []byte) {
	jsonData := map[string]any{}
	err := json.Unmarshal(stream, &jsonData)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//TODO
}

func main() {
	listener, err := gsn.ListenTCP(":8088", &MyBehavior{}, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	listener.Start()

	time.Sleep(time.Second)

	// ---------- 连接并发包 -----------
	conn, err := gsn.Dial("tcp", ":8088", nil, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	data := map[string]any{
		"cmd":     1001,
		"name":    "NAME",
		"data":    "abcdefghijklmnopqrstuvwxyz",
		"chinese": "一地在要工上是中国同和的有人我主产不为这民了发以经",
	}

	packer := gsn.NewJsonPacker()
	_ = packer.PackObject(data)

	_, _ = conn.Write(packer.GetData()) // 发送数据

	time.Sleep(3 * time.Second) // 预留足够的时间来完成包处理
}
```

### 6.设置 Deadline

gsn 可以在 Behavior 的 Convert 内设置 Deadline

```go
package main

import (
	"fmt"
	"github.com/olocat/gsn"
	"net"
	"time"
)

type MyBehavior struct {
	gsn.BaseBehavior
}

func OnPackage(conn net.Conn, stream []byte) {
	fmt.Println("收到数据包", stream)
}

func Convert(conn net.Conn) net.Conn {
	conn.SetDeadline(time.Now().Add(10 * time.Second))
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
}

func main() {
	conn, err := gsn.Dial("tcp", ":8088", &MyBehavior{}, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//TODO
}
```

### 7.自定义Conn

由于 Behavior 中的每个方法都会传递 conn， 但很多时候 conn 需要自己定义， 所以 Convert 主要用于将 net.Conn 转成自定义的 conn

```go
package main

import (
	"fmt"
	"github.com/olocat/gsn"
	"net"
)

type MyConn struct {
	net.Conn
	Id   uint32
	Name string
}

type MyBehavior struct {
	gsn.BaseBehavior
}

func OnPackage(conn net.Conn, stream []byte) {
	fmt.Printf("收到一个包 Id:%d Name:%s\n", conn.(*MyConn).Id, conn.(*MyConn).Name)
}

func Convert(conn net.Conn) net.Conn {
    return &MyConn{Conn: conn, Id: 100, Name: "this a name"}
}

func main() {
	conn, err := gsn.Dial("tcp", ":8088", &MyBehavior{}, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//TODO
}
```

## 说明

### 1.数据格式

数据格式为：

|--4b--|------|

|-包头-|--数据--|

没错，就是这么简单。 在发包的时候需要加上包头，包头表示总长（包含自己）

比如一个 12 个字节的数据包。

[ 0,0,0,12, 0,0,0,0,0,0,0,0 ]

前四位 0 0 0 12 表示该包总长度为 12 个字节，后 8 个 0 则为该包的数据

对于默认的 4b 包头来说，所能发送最大包长为 429496729 字节

如果需要调整包头，可以使用 SetHeadSize 方法。暂时支持 1b,2b,4b,8b 长度的包头

```go
package main

import (
	"fmt"
	"github.com/olocat/gsn"
)

func main() {
	listener, err := gsn.ListenTCP(":8088", nil, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	listener.SetHeadSize(gsn.FourWordHeadSize) // 包头改成 8b
	listener.Start()

	// 解包时也需要设计包头
	pack := gsn.NewBinPacker()
	pack.SetHeadSize(gsn.FourWordHeadSize) // 包头改成 8b
}
```

## 未来计划

1. 优雅的退出
2. protobuf 示例
3. 心跳示例
