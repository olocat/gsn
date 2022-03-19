package gsn

import (
	"fmt"
	"testing"
)

func TestPacker(t *testing.T) {
	pack := NewBinPacker()
	writeData(pack)
	unpack := NewBinUnpacker(pack.GetData())
	readData(unpack)
}

func readData(packer *Unpacker) {
	fmt.Println("readData", len(packer.DataStream))

	v1, _ := packer.UnPackUint16()
	v2, _ := packer.UnPackInt16()
	v3, _ := packer.UnPackUint32()
	v4, _ := packer.UnPackInt32()
	v5, _ := packer.UnPackUint64()
	v6, _ := packer.UnPackInt64()
	v7, _ := packer.UnPackByte()
	v8, _ := packer.UnPackBytes(int(v7))
	sz, _ := packer.UnPackUint32()
	str, err := packer.UnPackString(int(sz))
	if err != nil {
		fmt.Println("UnPackString err:", err)
	}

	fmt.Println(v1, v2, v3, v4, v5, v6)
	fmt.Println(v7, v8)
	fmt.Println(str)
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
	packer.PackUint32(uint32(len(str)))
	packer.PackString(str)
}
