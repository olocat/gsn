package gsn

import (
	"encoding/binary"
	"encoding/json"
	"errors"
)

const (
	ByteHeadSize       = 1
	WordHeadSize       = 2
	DoubleWordHeadSize = 4
	FourWordHeadSize   = 8
	DefaultHeadSize    = DoubleWordHeadSize
)

type Unpacker struct {
	DataStream []byte
}

type Packer struct {
	DataStream []byte
}

func NewBinPacker() *Packer {
	return &Packer{DataStream: []byte{}}
}

func NewBinUnpacker(dataStream []byte) *Unpacker {
	return &Unpacker{DataStream: dataStream}
}

/* -- Unpacker -- */

func (u *Unpacker) UnPackString(length int) (string, error) {
	byteList, err := u.UnPackBytes(length)
	if err != nil {
		return "", err
	}

	return string(byteList), nil
}

func (u *Unpacker) UnPackJsonData(length int) (map[string]interface{}, error) {
	byteList, err := u.UnPackBytes(length)
	if err != nil {
		return nil, err
	}
	jsonData := map[string]interface{}{}
	err = json.Unmarshal(byteList, &jsonData)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

func (u *Unpacker) UnPackBytes(length int) ([]byte, error) {
	if u.DataStream == nil || len(u.DataStream) < length {
		return nil, errors.New("can't unpacker byte array, insufficient data length")
	}

	value := u.DataStream[:length]
	u.DataStream = u.DataStream[length:]
	return value, nil
}

func (u *Unpacker) UnPackByte() (byte, error) {
	if u.DataStream == nil || len(u.DataStream) < 1 {
		return 0, errors.New("can't unpacker byte, insufficient data length")
	}
	value := u.DataStream[0]
	u.DataStream = u.DataStream[1:]
	return value, nil
}

func (u *Unpacker) UnPackUint16() (uint16, error) {
	if u.DataStream == nil || len(u.DataStream) < 2 {
		return 0, errors.New("can't unpacker uint16, insufficient data length")
	}

	value := binary.BigEndian.Uint16(u.DataStream)
	u.DataStream = u.DataStream[2:]
	return value, nil
}

func (u *Unpacker) UnPackInt16() (int16, error) {
	value, err := u.UnPackUint16()
	return int16(value), err
}

func (u *Unpacker) UnPackUint32() (uint32, error) {
	if u.DataStream == nil || len(u.DataStream) < 4 {
		return 0, errors.New("can't unpacker uint32,insufficient data length")
	}

	value := binary.BigEndian.Uint32(u.DataStream)
	u.DataStream = u.DataStream[4:]
	return value, nil
}

func (u *Unpacker) UnPackInt32() (int32, error) {
	value, err := u.UnPackUint32()
	return int32(value), err
}

func (u *Unpacker) UnPackUint64() (uint64, error) {
	if u.DataStream == nil || len(u.DataStream) < 8 {
		return 0, errors.New("can't unpacker uint64,insufficient data length")
	}

	value := binary.BigEndian.Uint64(u.DataStream)
	u.DataStream = u.DataStream[8:]
	return value, nil
}

func (u *Unpacker) UnPackInt64() (int64, error) {
	value, err := u.UnPackUint64()
	return int64(value), err
}

/* ------------ */

/* -- Packer -- */

func (p *Packer) PackJsonData(jsonData map[string]interface{}) error {
	value, err := json.Marshal(jsonData)
	if err != nil {
		return err
	}
	p.DataStream = append(p.DataStream, value...)
	return nil
}

func (p *Packer) PackString(str string) {
	value := []byte(str)
	p.DataStream = append(p.DataStream, value...)
}

func (p *Packer) PackBytes(value []byte) {
	p.DataStream = append(p.DataStream, value...)
}

func (p *Packer) PackByte(value byte) {
	p.DataStream = append(p.DataStream, value)
}

func (p *Packer) PackUint16(value uint16) {
	byteList := make([]byte, 2)
	binary.BigEndian.PutUint16(byteList, value)
	p.DataStream = append(p.DataStream, byteList...)
}

func (p *Packer) PackInt16(value int16) {
	p.PackUint16(uint16(value))
}

func (p *Packer) PackUint32(value uint32) {
	byteList := make([]byte, 4)
	binary.BigEndian.PutUint32(byteList, value)
	p.DataStream = append(p.DataStream, byteList...)
}

func (p *Packer) PackInt32(value int32) {
	p.PackUint32(uint32(value))
}

func (p *Packer) PackUint64(value uint64) {
	byteList := make([]byte, 8)
	binary.BigEndian.PutUint64(byteList, value)
	p.DataStream = append(p.DataStream, byteList...)
}

func (p *Packer) PackInt64(value int64) {
	p.PackUint64(uint64(value))
}

func (p *Packer) GetData() []byte {
	return p.DataStream
}

/* ------------ */
