package game

import (
	chess "github.com/tommjj/chess_OG/chess_core"
)

var (
	ErrInvalidFEN        = chess.ErrInvalidFEN
	ErrInvalidPiece      = chess.ErrInvalidPiece
	ErrInvalidCastling   = chess.ErrInvalidCastling
	ErrInvalidEnPassant  = chess.ErrInvalidEnPassant
	ErrInvalidHalfmove   = chess.ErrInvalidHalfmove
	ErrInvalidFullmove   = chess.ErrInvalidFullmove
	ErrInvalidSideToMove = chess.ErrInvalidSideToMove
	ErrMultipleKings     = chess.ErrMultipleKings
	ErrNoKing            = chess.ErrNoKing
	ErrPawnOnFirstOrLast = chess.ErrPawnOnFirstOrLast
	ErrTooManyPieces     = chess.ErrTooManyPieces
	ErrNegativeHalfmove  = chess.ErrNegativeHalfmove
	ErrNegativeFullmove  = chess.ErrNegativeFullmove
	ErrInvalidMove       = chess.ErrInvalidMove
	ErrIllegalMove       = chess.ErrIllegalMove
	ErrNoMovesAvailable  = chess.ErrNoMovesAvailable

	ErrCheckmate = chess.ErrCheckmate
	ErrStalemate = chess.ErrStalemate
	ErrMatchEnd  = chess.ErrMatchEnd

	ErrInsufficientMaterial = chess.ErrInsufficientMaterial
	ErrThreefoldRepetition  = chess.ErrThreefoldRepetition

	ErrInvalidPromotion = chess.ErrInvalidPromotion
)
