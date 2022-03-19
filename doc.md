# 新版本使用文档

---

> 该文档中的用法暂未实现
> 是未来中的计划

## 开始

### 1.创建TCP监听

```go
package main

import (
	"github.com/olocat/gsn"
	"log"
	"net"
)

// ExampleBehavior  一个方法类
// gsn.Behavior     方法接口
// gsn.BaseBehavior 是一个包含Behavior初始实现方法的基类，内嵌该类可以省去实现一些不需要的方法
type ExampleBehavior struct {
	gsn.BaseBehavior
}

// OnConnect 当监听到其它请求进来时，连接成功后会调用该方法，只有 Listen 会用到该方法
func (e *ExampleBehavior) OnConnect(content *gsn.Content) {
	var conn net.Conn = content.Conn
	log.Println("connect in", conn.RemoteAddr())
}

// OnPackage 当收到包后，数据会传到该方法
// 由自己选择使用哪种解包方式
func (e *ExampleBehavior) OnPackage(content *gsn.Content, stream []byte) {
	// 使用二进制解包
	recv := gsn.NewBinUnpacker(stream)

	// 不使用泛型的方式
	var byte_1 byte = recv.UnpackByte()
	var int8_1 int8 = int8(recv.UnpackByte())

	// 使用泛型的方式
	var byte_2 byte = recv.UnpackByte[byte]()
	var int8_2 int8 = recv.UnpackByte[int8]()
	// 暂时未考虑好是否引入泛型

	var string_1 string = recv.UnpackString()

	// 使用Json解包
	recv = gsn.NewJsonUnpacker(stream)
	var json_1 map[string]any = recv.UnpackJson()
	var json_2 JsonObject = JsonObject{}
	recv.UnpackObject(&json_2)

	// 使用XML解包
	recv = gsn.NewXmlUnpacker(stream)
	var xml_1 XmlObject = XmlObject{}
	recv.UnpackObject(&xml_1)

	// 使用 protobuf 解包
	recv = gsn.NewProtobufUnpacker(stream)
	var pb_1 PBObject = PBObject{}
	recv.UnpackObject(&pb_1)
}

func example() {
	behavior := &ExampleBehavior{}
	listen, err := gsn.Listen("tcp", ":8088", behavior)
	if err != nil {
		log.Println(err.Error())
		return
	}
}
```