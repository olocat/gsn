package gsn

import (
	"encoding/json"
	"fmt"
	"net"
	"testing"
	"time"
)

func TestPacker(t *testing.T) {
	listen := ListenTCP(":8081")
	listen.OnReceive = func(pack UnPacker) {
		readData(pack)
	}
	go listen.Start()

	time.Sleep(time.Second)
	conn, err := net.Dial("tcp", ":8081")
	defer conn.Close()
	if err != nil {
		fmt.Println(err)
	}
	var packer = Packer{}
	writeData(&packer)
	sendData := packer.GetData()
	conn.Write(sendData)
	time.Sleep(2 * time.Second)
}

func readData(packer UnPacker) {
	fmt.Println("readData", len(*packer.data))
	v1, _ := packer.UnPackUint16()
	v2, _ := packer.UnPackInt16()
	v3, _ := packer.UnPackUint32()
	v4, _ := packer.UnPackInt32()
	v5, _ := packer.UnPackUint64()
	v6, _ := packer.UnPackInt64()
	v7, _ := packer.UnPackByte()
	v8, _ := packer.UnPackBytes(int(v7))

	sl, _ := packer.UnPackUint16()
	str, err := packer.UnPackString(int(sl))
	if err != nil {
		fmt.Println("UnPackString err:", err)
	}

	jl, _ := packer.UnPackUint16()
	jsonData, err := packer.UnPackJsonData(int(jl))
	if err != nil {
		fmt.Println("UnPackJsonData err:", err)
	}

	fmt.Println(v1, v2, v3, v4, v5, v6)
	fmt.Println(v7, v8)
	fmt.Println(str)
	fmt.Println(jsonData["zhangsan"], jsonData["lisi"], jsonData["wangwu"])
}

func writeData(packer *Packer) {
	packer.PackUint16(65535)
	packer.PackInt16(-32768)
	packer.PackUint32(4294967295)
	packer.PackInt32(-2147483647)
	packer.PackUint64(18446744073709551615)
	packer.PackInt64(9223372036854775807)
	packer.PackByte(5)
	packer.PackBytes([]byte{0, 1, 2, 3, 255})

	str := "zhangsan lisi wangwu 都是认证啊 !`||。。，，"
	strBytes := []byte(str)
	strLength := uint16(len(strBytes))
	packer.PackUint16(strLength)
	packer.PackString(str)

	jsonData := map[string]interface{}{}
	jsonData["zhangsan"] = "1001"
	jsonData["lisi"] = "1002"
	jsonData["wangwu"] = "1003"
	jsonBytes, _ := json.Marshal(jsonData)
	jsonLength := uint16(len(jsonBytes))
	packer.PackUint16(jsonLength)
	_ = packer.PackJsonData(jsonData)
}
