package gsn

import (
	"encoding/binary"
	"encoding/json"
	"errors"
)

type UnPacker struct {
	DataStream []byte
}

type Packer struct {
	DataStream []byte
}

func NewPacker() *Packer {
	return &Packer{DataStream: make([]byte, PackHeadSize)}
}

func NewUnPacker(dataStream []byte) *UnPacker {
	return &UnPacker{DataStream: dataStream}
}

/* -- UnPacker -- */

func (u *UnPacker) UnPackString(length int) (string, error) {
	byteList, err := u.UnPackBytes(length)
	if err != nil {
		return "", err
	}

	return string(byteList), nil
}

func (u *UnPacker) UnPackJsonData(length int) (map[string]interface{}, error) {
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

func (u *UnPacker) UnPackBytes(length int) ([]byte, error) {
	if u.DataStream == nil || len(u.DataStream) < length {
		return nil, errors.New("")
	}

	value := u.DataStream[:length]
	u.DataStream = u.DataStream[length:]
	return value, nil
}

func (u *UnPacker) UnPackByte() (byte, error) {
	if u.DataStream == nil || len(u.DataStream) < 1 {
		return 0, errors.New("")
	}
	value := u.DataStream[0]
	u.DataStream = u.DataStream[1:]
	return value, nil
}

func (u *UnPacker) UnPackUint16() (uint16, error) {
	if u.DataStream == nil || len(u.DataStream) < 2 {
		return 0, errors.New("")
	}

	value := binary.BigEndian.Uint16(u.DataStream)
	u.DataStream = u.DataStream[2:]
	return value, nil
}

func (u *UnPacker) UnPackInt16() (int16, error) {
	value, err := u.UnPackUint16()
	return int16(value), err
}

func (u *UnPacker) UnPackUint32() (uint32, error) {
	if u.DataStream == nil || len(u.DataStream) < 4 {
		return 0, errors.New("")
	}

	value := binary.BigEndian.Uint32(u.DataStream)
	u.DataStream = u.DataStream[4:]
	return value, nil
}

func (u *UnPacker) UnPackInt32() (int32, error) {
	value, err := u.UnPackUint32()
	return int32(value), err
}

func (u *UnPacker) UnPackUint64() (uint64, error) {
	if u.DataStream == nil || len(u.DataStream) < 8 {
		return 0, errors.New("")
	}

	value := binary.BigEndian.Uint64(u.DataStream)
	u.DataStream = u.DataStream[8:]
	return value, nil
}

func (u *UnPacker) UnPackInt64() (int64, error) {
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

func (p *Packer) computeDataLength() {
	dataLen := uint32(len(p.DataStream))
	binary.BigEndian.PutUint32(p.DataStream, dataLen)
}

func (p *Packer) GetData() []byte {
	p.computeDataLength()
	return p.DataStream
}

/* ------------ */
