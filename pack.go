package gsn

import (
	"encoding/json"
	"errors"
	"fmt"
	"gsn/util"
)

type UnPacker struct {
	data *[]byte
}

type Packer struct {
	data *[]byte
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
	if u.data == nil || len(*u.data) < length {
		return nil, errors.New("")
	}

	value := (*u.data)[:length]
	*u.data = (*u.data)[length:]
	return value, nil
}

func (u *UnPacker) UnPackByte() (byte, error) {
	if u.data == nil || len(*u.data) < 1 {
		return 0, errors.New("")
	}
	value := (*u.data)[0]
	*u.data = (*u.data)[1:]
	return value, nil
}

func (u *UnPacker) UnPackUint16() (uint16, error) {
	if u.data == nil {
		return 0, errors.New("")
	}

	value, err := util.BytesToUInt16(*u.data)
	*u.data = (*u.data)[2:]
	return value, err
}

func (u *UnPacker) UnPackInt16() (int16, error) {
	value, err := u.UnPackUint16()
	return int16(value), err
}

func (u *UnPacker) UnPackUint32() (uint32, error) {
	if u.data == nil {
		return 0, errors.New("")
	}

	value, err := util.BytesToUInt32(*u.data)
	*u.data = (*u.data)[4:]
	return value, err
}

func (u *UnPacker) UnPackInt32() (int32, error) {
	value, err := u.UnPackUint32()
	return int32(value), err
}

func (u *UnPacker) UnPackUint64() (uint64, error) {
	if u.data == nil {
		return 0, errors.New("")
	}

	value, err := util.BytesToUInt64(*u.data)
	*u.data = (*u.data)[8:]
	return value, err
}

func (u *UnPacker) UnPackInt64() (int64, error) {
	value, err := u.UnPackUint64()
	return int64(value), err
}

/* ------------ */

/* -- Packer -- */

func (p *Packer) initData() {
	if p.data == nil {
		p.data = &[]byte{0, 0}
	}
}

func (p *Packer) PackJsonData(jsonData map[string]interface{}) error {
	value, err := json.Marshal(jsonData)
	if err != nil {
		return err
	}
	*p.data = append(*p.data, value...)
	return nil
}

func (p *Packer) PackString(str string) {
	p.initData()
	value := []byte(str)
	*p.data = append(*p.data, value...)
}

func (p *Packer) PackBytes(value []byte) {
	p.initData()
	*p.data = append(*p.data, value...)
}

func (p *Packer) PackByte(value byte) {
	p.initData()
	*p.data = append(*p.data, value)
}

func (p *Packer) PackUint16(value uint16) {
	p.initData()
	byteList := util.UInt16ToBytes(value)
	*p.data = append(*p.data, byteList...)
}

func (p *Packer) PackInt16(value int16) {
	p.initData()
	byteList := util.UInt16ToBytes(uint16(value))
	*p.data = append(*p.data, byteList...)
}

func (p *Packer) PackUint32(value uint32) {
	p.initData()
	byteList := util.UInt32ToBytes(value)
	*p.data = append(*p.data, byteList...)
}

func (p *Packer) PackInt32(value int32) {
	p.initData()
	byteList := util.UInt32ToBytes(uint32(value))
	*p.data = append(*p.data, byteList...)
}

func (p *Packer) PackUint64(value uint64) {
	p.initData()
	byteList := util.UInt64ToBytes(value)
	*p.data = append(*p.data, byteList...)
}

func (p *Packer) PackInt64(value int64) {
	p.initData()
	byteList := util.UInt64ToBytes(uint64(value))
	*p.data = append(*p.data, byteList...)
}

func (p *Packer) computeDataLength() {
	p.initData()

	fmt.Println("computeDataLength len:", len(*p.data))

	dataLen := uint16(len(*p.data))
	dataLenBytes := util.UInt16ToBytes(dataLen)
	(*p.data)[0] = dataLenBytes[0]
	(*p.data)[1] = dataLenBytes[1]
}

func (p *Packer) GetData() []byte {
	p.initData()
	p.computeDataLength()

	fmt.Println("GetData ", (*p.data)[0], (*p.data)[1])
	return *p.data
}

/* ------------ */
