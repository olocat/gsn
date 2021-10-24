package gsn

import (
	"fmt"
	"net"
	"testing"
	"time"
)

func TestPacker(t *testing.T) {
	listen := ListenTCP(":8081")
	listen.OnReceive = func(pack *ReceivePack) {
		readData(pack)
	}
	go listen.Start()

	time.Sleep(time.Second)
	conn, err := net.Dial("tcp", ":8081")
	defer conn.Close()
	if err != nil {
		fmt.Println(err)
	}
	ctx := NewContext(conn)
	var packer = NewSendPack(ctx)
	writeData(packer)
	packer.Send()
	time.Sleep(2 * time.Second)
}

func readData(packer *ReceivePack) {
	fmt.Println("readData", len(packer.DataStream))
	v1, _ := packer.UnPackUint16()
	v2, _ := packer.UnPackInt16()
	v3, _ := packer.UnPackUint32()
	v4, _ := packer.UnPackInt32()
	v5, _ := packer.UnPackUint64()
	v6, _ := packer.UnPackInt64()
	v7, _ := packer.UnPackByte()
	v8, _ := packer.UnPackBytes(int(v7))
	str, err := packer.UnPackString()
	if err != nil {
		fmt.Println("UnPackString err:", err)
	}

	jsonData, err := packer.UnPackJsonData()
	if err != nil {
		fmt.Println("UnPackJsonData err:", err)
	}

	fmt.Println(v1, v2, v3, v4, v5, v6)
	fmt.Println(v7, v8)
	fmt.Println(str)
	fmt.Println(jsonData["zhangsan"], jsonData["lisi"], jsonData["wangwu"])
}

func writeData(packer *SendPack) {
	packer.PackUint16(65535)
	packer.PackInt16(-32768)
	packer.PackUint32(4294967295)
	packer.PackInt32(-2147483647)
	packer.PackUint64(18446744073709551615)
	packer.PackInt64(9223372036854775807)
	packer.PackByte(5)
	packer.PackBytes([]byte{0, 1, 2, 3, 255})

	str := "zhangsan lisi wangwu 都是认证啊 !`||。。，，"
	packer.PackString(str)

	jsonData := map[string]interface{}{}
	jsonData["zhangsan"] = "1001"
	jsonData["lisi"] = "1002"
	jsonData["wangwu"] = "1003"
	_ = packer.PackJsonData(jsonData)
}
