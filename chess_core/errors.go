package chess_core

import (
	"errors"
)

var (
	ErrInvalidFEN        = errors.New("invalid FEN string")
	ErrInvalidPiece      = errors.New("invalid piece in FEN string")
	ErrInvalidCastling   = errors.New("invalid castling rights in FEN string")
	ErrInvalidEnPassant  = errors.New("invalid en passant square in FEN string")
	ErrInvalidHalfmove   = errors.New("invalid halfmove clock in FEN string")
	ErrInvalidFullmove   = errors.New("invalid fullmove number in FEN string")
	ErrInvalidSideToMove = errors.New("invalid side to move in FEN string")
	ErrMultipleKings     = errors.New("multiple kings found for one color")
	ErrNoKing            = errors.New("no king found for one color")
	ErrPawnOnFirstOrLast = errors.New("pawn on first or last rank")
	ErrTooManyPieces     = errors.New("too many pieces on the board")
	ErrNegativeHalfmove  = errors.New("negative halfmove clock")
	ErrNegativeFullmove  = errors.New("negative fullmove number")
	ErrInvalidMove       = errors.New("invalid move")
	ErrIllegalMove       = errors.New("illegal move")
	ErrNoMovesAvailable  = errors.New("no moves available")

	ErrCheckmate = errors.New("checkmate")
	ErrStalemate = errors.New("stalemate")
	ErrMatchEnd  = errors.New("the match ended")

	ErrInsufficientMaterial = errors.New("insufficient material to continue the game")

	ErrInvalidPromotion = errors.New("invalid promotion")
)

func wrapError(err error, message string) error {
	return errors.New(message + ": " + err.Error())
}
