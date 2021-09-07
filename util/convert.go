package util

import "errors"

const (
	DefaultValue          = 0
	BytesToUInt16NeedSize = 2
	BytesToUInt32NeedSize = 4
)

func BytesToUInt32(byteArray []byte, defaultValue uint32) uint32 {
	if byteArray == nil {
		return defaultValue
	}

	l := len(byteArray)

	if l == 0 {
		return defaultValue
	}

	if l == 1 {
		return uint32(byteArray[0])
	}

	if l > 4 {
		l = 4
	}

	var value uint32 = 0
	for i := 0; i < l; i++ {
		value <<= 8
		value += uint32(byteArray[i])
	}

	return value
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
