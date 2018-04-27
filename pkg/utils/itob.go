package utils

import (
	"encoding/binary"
)

func IToB(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func BToI(b []byte) int {
	b = padOrTrim(b, 8)
	return int(binary.BigEndian.Uint64(b))
}

func padOrTrim(bb []byte, size int) []byte {
	l := len(bb)
	if l == size {
		return bb
	}
	if l > size {
		return bb[l-size:]
	}
	tmp := make([]byte, size)
	copy(tmp[size-l:], bb)
	return tmp
}
