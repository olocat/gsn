package gsn

import (
	"encoding/binary"
	"encoding/json"
	"encoding/xml"
)

const (
	ByteHeadSize       = 1
	WordHeadSize       = 2
	DoubleWordHeadSize = 4
	FourWordHeadSize   = 8
	DefaultHeadSize    = DoubleWordHeadSize
)

type Packer struct {
	headSize   byte
	dataStream []byte
}

type BinPacker struct {
	*Packer
}

type JsonPacker struct {
	*Packer
}

type XmlPacker struct {
	*Packer
}

func NewPacker(headSize byte) *Packer {
	headSize = correctHeadSize(headSize)
	return &Packer{dataStream: make([]byte, headSize), headSize: headSize}
}

func NewBinPacker() *BinPacker {
	return &BinPacker{NewPacker(DefaultHeadSize)}
}

func NewJsonPacker() *JsonPacker {
	return &JsonPacker{NewPacker(DefaultHeadSize)}
}

func NewXmlPacker() *XmlPacker {
	return &XmlPacker{NewPacker(DefaultHeadSize)}
}

func (p *Packer) SetHeadSize(headSize byte) {
	headSize = correctHeadSize(headSize)
	if headSize == p.headSize {
		return
	}

	length := transInt(p.dataStream[:p.headSize])
	newHeadBytes := make([]byte, headSize, headSize)
	putInt(headSize, newHeadBytes, length)
	p.dataStream = append(newHeadBytes, p.dataStream[headSize:]...)
}

func (p *Packer) Append(stream ...byte) {
	p.dataStream = append(p.dataStream, stream...)
}

func (p *Packer) GetData() []byte {
	length := len(p.dataStream)
	putInt(p.headSize, p.dataStream, uint64(length))
	return p.dataStream
}

func (p *Packer) Clean() {
	p.dataStream = make([]byte, p.headSize)
}

func (p *BinPacker) PackJsonData(jsonData map[string]interface{}) error {
	value, err := json.Marshal(jsonData)
	if err != nil {
		return err
	}
	p.dataStream = append(p.dataStream, value...)
	return nil
}

func (p *BinPacker) PackString(str string) {
	value := []byte(str)
	p.Append(value...)
}

func (p *BinPacker) PackBytes(value []byte) {
	p.Append(value...)
}

func (p *BinPacker) PackByte(value byte) {
	p.Append(value)
}

func (p *BinPacker) PackUint16(value uint16) {
	byteList := make([]byte, 2)
	binary.BigEndian.PutUint16(byteList, value)
	p.Append(byteList...)
}

func (p *BinPacker) PackInt16(value int16) {
	p.PackUint16(uint16(value))
}

func (p *BinPacker) PackUint32(value uint32) {
	byteList := make([]byte, 4)
	binary.BigEndian.PutUint32(byteList, value)
	p.Append(byteList...)
}

func (p *BinPacker) PackInt32(value int32) {
	p.PackUint32(uint32(value))
}

func (p *BinPacker) PackUint64(value uint64) {
	byteList := make([]byte, 8)
	binary.BigEndian.PutUint64(byteList, value)
	p.Append(byteList...)
}

func (p *BinPacker) PackInt64(value int64) {
	p.PackUint64(uint64(value))
}

func (p *JsonPacker) PackObject(v any) error {
	stream, err := json.Marshal(v)
	if err != nil {
		return err
	}

	p.Append(stream...)
	return nil
}

func (p *XmlPacker) PackObject(v any) error {
	stream, err := xml.Marshal(v)
	if err != nil {
		return err
	}

	p.Append(stream...)
	return nil
}
