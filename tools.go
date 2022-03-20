package gsn

import (
	"encoding/binary"
	"io"
)

func readFull(reader io.Reader, buf []byte) error {
	if buf == nil || len(buf) == 0 {
		return nil
	}

	offset := 0
	for {
		n, e := reader.Read(buf[offset:])
		if e != nil {
			return e
		}
		offset += n
		if offset >= len(buf) {
			break
		}
	}
	return nil
}

func transInt(b []byte) uint64 {
	if b == nil || len(b) == 0 {
		return 0
	}

	size := len(b)
	if size == ByteHeadSize {
		return uint64(b[0])
	}

	if size == WordHeadSize {
		v := binary.BigEndian.Uint16(b)
		return uint64(v)
	}

	if size == DoubleWordHeadSize {
		v := binary.BigEndian.Uint32(b)
		return uint64(v)
	}

	if size == FourWordHeadSize {
		v := binary.BigEndian.Uint64(b)
		return v
	}

	return 0
}

func putInt(headSize byte, b []byte, v uint64) {
	if b == nil || len(b) == 0 {
		return
	}

	size := len(b)
	if headSize == ByteHeadSize {
		if size < 1 {
			return
		}
		b[0] = byte(v)
		return
	}

	if headSize == WordHeadSize {
		if size < 2 {
			return
		}
		binary.BigEndian.PutUint16(b, uint16(v))
		return
	}

	if headSize == DoubleWordHeadSize {
		if size < 4 {
			return
		}
		binary.BigEndian.PutUint32(b, uint32(v))
		return
	}

	if headSize == FourWordHeadSize {
		if size < 8 {
			return
		}
		binary.BigEndian.PutUint64(b, v)
		return
	}
}

//correctHeadSize Only 1,2,4,8 allowed
func correctHeadSize(headSize byte) byte {
	if headSize <= 0 {
		return DefaultHeadSize
	}

	if headSize == 3 || (headSize > DoubleWordHeadSize && headSize < FourWordHeadSize) {
		return DefaultHeadSize
	}

	if headSize > FourWordHeadSize {
		return FourWordHeadSize
	}

	return DefaultHeadSize
}
