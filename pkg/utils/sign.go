package utils

import "encoding/binary"

func PrepareSignData(chainID uint64, data []byte) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(chainID))

	return append(buf, data...)
}
