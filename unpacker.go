package gsn

import (
	"encoding/binary"
	"encoding/json"
	"errors"
)

type BinUnpacker struct {
	DataStream []byte
}

func NewBinUnpacker(dataStream []byte) *BinUnpacker {
	return &BinUnpacker{DataStream: dataStream}
}

func (u *BinUnpacker) UnpackString(length int) (string, error) {
	byteList, err := u.UnpackBytes(length)
	if err != nil {
		return "", err
	}

	return string(byteList), nil
}

func (u *BinUnpacker) UnpackJsonData(length int) (map[string]interface{}, error) {
	byteList, err := u.UnpackBytes(length)
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

func (u *BinUnpacker) UnpackBytes(length int) ([]byte, error) {
	if u.DataStream == nil || len(u.DataStream) < length {
		return nil, errors.New("can't unpacker byte array, insufficient data length")
	}

	value := u.DataStream[:length]
	u.DataStream = u.DataStream[length:]
	return value, nil
}

func (u *BinUnpacker) UnpackByte() (byte, error) {
	if u.DataStream == nil || len(u.DataStream) < 1 {
		return 0, errors.New("can't unpacker byte, insufficient data length")
	}
	value := u.DataStream[0]
	u.DataStream = u.DataStream[1:]
	return value, nil
}

func (u *BinUnpacker) UnpackUint16() (uint16, error) {
	if u.DataStream == nil || len(u.DataStream) < 2 {
		return 0, errors.New("can't unpacker uint16, insufficient data length")
	}

	value := binary.BigEndian.Uint16(u.DataStream)
	u.DataStream = u.DataStream[2:]
	return value, nil
}

func (u *BinUnpacker) UnpackInt16() (int16, error) {
	value, err := u.UnpackUint16()
	return int16(value), err
}

func (u *BinUnpacker) UnpackUint32() (uint32, error) {
	if u.DataStream == nil || len(u.DataStream) < 4 {
		return 0, errors.New("can't unpacker uint32,insufficient data length")
	}

	value := binary.BigEndian.Uint32(u.DataStream)
	u.DataStream = u.DataStream[4:]
	return value, nil
}

func (u *BinUnpacker) UnpackInt32() (int32, error) {
	value, err := u.UnpackUint32()
	return int32(value), err
}

func (u *BinUnpacker) UnpackUint64() (uint64, error) {
	if u.DataStream == nil || len(u.DataStream) < 8 {
		return 0, errors.New("can't unpacker uint64,insufficient data length")
	}

	value := binary.BigEndian.Uint64(u.DataStream)
	u.DataStream = u.DataStream[8:]
	return value, nil
}

func (u *BinUnpacker) UnpackInt64() (int64, error) {
	value, err := u.UnpackUint64()
	return int64(value), err
}
