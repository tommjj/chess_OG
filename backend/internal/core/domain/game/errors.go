package game

import (
	"errors"

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
	ErrNoMovesAvailable  = chess.ErrNoMovesAvailable

	ErrMatchEnd = chess.ErrMatchEnd

	ErrInvalidMove      = chess.ErrInvalidMove
	ErrInvalidPromotion = chess.ErrInvalidPromotion
	ErrMoveIntoCheck    = chess.ErrMoveIntoCheck
	ErrMoveOutOfTurn    = chess.ErrMoveOutOfTurn

	ErrTimeout        = errors.New("error timeout")
	ErrGamePaused     = errors.New("error game paused")
	ErrGameNotStarted = errors.New("error game not started")

	// Create game errors
	ErrInvalidGameMode = errors.New("error invalid game mode")
)
