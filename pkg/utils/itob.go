package utils

import (
	"encoding/binary"
)

func IToB(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
