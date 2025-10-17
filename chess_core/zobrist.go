package chess_core

import (
	"math/rand"
)

func init() {
	initRandomZobristKey()
}

// random piece keys [piece][square]
var pieceKeys [12][64]uint64

// random enpassant keys [square]
var enpassantKeys [64]uint64

// random castling keys
var castleKeys [16]uint64

// random side key
var sideKey uint64

func initRandomZobristKey() {
	// 12 pieces × 64 squares
	for piece := range 12 {
		for sq := range 64 {
			pieceKeys[piece][sq] = rand.Uint64() // random 64-bit int
		}
	}

	// 64 en passant target squares
	for sq := range 64 {
		enpassantKeys[sq] = rand.Uint64()
	}

	// 16 castling rights states (bitmask 0–15)
	for i := range 16 {
		castleKeys[i] = rand.Uint64()
	}

	// random side to move key
	sideKey = rand.Uint64()
}

func computeZobristHash(bbs *BitBoards, sideToMove Color, epSp Square, castlingRights int) uint64 {
	var h uint64

	var b BitBoard

	for piece := 1; piece < len(ASCIIPieces); piece++ {
		b = bbs.Pieces(ASCIIPieces[piece])
		for b != 0 {
			sq := popLSB(&b)
			h ^= pieceKeys[piece-1][sq]
		}
	}

	if epSp.IsValid() {
		h ^= enpassantKeys[epSp]
	}

	h ^= castleKeys[castlingRights]

	if sideToMove == Black {
		h ^= sideKey
	}

	return h
}
