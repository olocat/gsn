package gsn

import (
	"fmt"
	"testing"
)

func TestBinPacker(t *testing.T) {
	pack := NewBinPacker()
	writeBinData(pack)
	unpack := NewBinUnpacker(pack.GetData())
	readBinData(unpack)
}

func readBinData(packer *BinUnpacker) {
	fmt.Println("readBinData", len(packer.DataStream))
	_, _ = packer.UnpackUint32()
	v1, _ := packer.UnpackUint16()
	v2, _ := packer.UnpackInt16()
	v3, _ := packer.UnpackUint32()
	v4, _ := packer.UnpackInt32()
	v5, _ := packer.UnpackUint64()
	v6, _ := packer.UnpackInt64()
	v7, _ := packer.UnpackByte()
	v8, _ := packer.UnpackBytes(int(v7))
	sz, _ := packer.UnpackUint32()
	str, err := packer.UnpackString(int(sz))
	if err != nil {
		fmt.Println("UnPackString err:", err)
	}

	fmt.Println(v1, v2, v3, v4, v5, v6)
	fmt.Println(v7, v8)
	fmt.Println(str)
}

func writeBinData(packer *BinPacker) {
	packer.PackUint16(65535)
	packer.PackInt16(-32768)
	packer.PackUint32(4294967295)
	packer.PackInt32(-2147483647)
	packer.PackUint64(18446744073709551615)
	packer.PackInt64(9223372036854775807)
	packer.PackByte(5)
	packer.PackBytes([]byte{0, 1, 2, 3, 255})

	str := "一地在要工上是中国同和的有人我主产不为这民了发以经"
	packer.PackUint32(uint32(len(str)))
	packer.PackString(str)
}
