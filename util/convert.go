package util

import "errors"

const (
	DefaultValue = 0
)

func BytesToUInt64(byteArray []byte) (uint64, error) {
	if byteArray == nil {
		return DefaultValue, errors.New("array data size insufficient")
	}

	l := len(byteArray)

	if l == 0 {
		return DefaultValue, errors.New("array data size insufficient")
	}

	if l == 1 {
		return uint64(byteArray[0]), nil
	}

	if l > 8 {
		l = 8
	}

	var value uint64 = 0
	for i := 0; i < l; i++ {
		value <<= 8
		value += uint64(byteArray[i])
	}

	return value, nil
}

func BytesToUInt32(byteArray []byte) (uint32, error) {
	if byteArray == nil {
		return DefaultValue, errors.New("array data size insufficient")
	}

	l := len(byteArray)

	if l == 0 {
		return DefaultValue, errors.New("array data size insufficient")
	}

	if l == 1 {
		return uint32(byteArray[0]), nil
	}

	if l > 4 {
		l = 4
	}

	var value uint32 = 0
	for i := 0; i < l; i++ {
		value <<= 8
		value += uint32(byteArray[i])
	}

	return value, nil
}

func BytesToUInt16(byteArray []byte) (uint16, error) {
	if byteArray == nil {
		return DefaultValue, errors.New(" array data is nil")
	}

	l := len(byteArray)

	if l == 0 {
		return DefaultValue, errors.New("array data size insufficient")
	}

	if l == 1 {
		return uint16(byteArray[0]), nil
	}

	if l > 2 {
		l = 2
	}

	var value uint16 = 0
	for i := 0; i < l; i++ {
		value <<= 8
		value += uint16(byteArray[i])
	}

	return value, nil
}

func UInt16ToBytes(u uint16) []byte {
	data := make([]byte, 2)
	data[0] = byte(u >> 8)
	data[1] = byte(u)
	return data
}

func UInt32ToBytes(u uint32) []byte {
	data := make([]byte, 4)
	for i := 0; i < 4; i++ {
		data[i] = byte(u >> (32 - (8 * (i + 1))))
	}
	return data
}

func UInt64ToBytes(u uint64) []byte {
	data := make([]byte, 8)
	for i := 0; i < 8; i++ {
		data[i] = byte(u >> (64 - (8 * (i + 1))))
	}
	return data
}
