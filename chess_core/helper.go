package main

import (
	"unsafe"
)

func u(i int) uint32 {
	return uint32(i)
}

func b(b bool) uint32 {
	if b {
		return 1
	}
	return 0
}

func toStringBit[T interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}](b T) string {
	size := int(unsafe.Sizeof(b)) * 8

	var result string
	for i := size - 1; i >= 0; i-- {
		if (b & (1 << i)) != 0 {
			result += "1"
		} else {
			result += "0"
		}
		if i%8 == 0 && i != 0 {
			result += "_"
		}
	}
	return result
}

// popLSB - pop least significant bit and return its index
func popLSB(bb *BitBoard) Square {
	if *bb == 0 {
		return -1
	}

	sq := bb.LeastSignificantBit()

	lsb := *bb & -*bb // isolate LS1B
	*bb ^= lsb        // xóa bit
	return Square(sq)
}

// Helper functions to get rank masks
func rankMask(rank int) BitBoard {
	return BitBoard(0xFF) << (rank * 8)
}

// Helper functions to get file masks
func fileMask(file int) BitBoard {
	var mask BitBoard = 0x0101010101010101
	return mask << file
}

func abs[T interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}](v T) T {
	if v < 0 {
		return -v
	}
	return v
}

func boolToUint32(v bool) uint32 {
	if v {
		return 1
	}
	return 0
}
