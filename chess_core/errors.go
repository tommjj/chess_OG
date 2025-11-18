package chess_core

import (
	"errors"
	"fmt"
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

	ErrNoMovesAvailable = errors.New("no moves available")

	ErrMatchEnd = errors.New("the match ended")

	ErrInvalidMove      = errors.New("invalid move")
	ErrInvalidPromotion = fmt.Errorf("%w: invalid promotion", ErrInvalidMove)
	ErrMoveIntoCheck    = fmt.Errorf("%w: move into check", ErrInvalidMove)
	ErrMoveOutOfTurn    = fmt.Errorf("%w: move out of turn", ErrInvalidMove)

	ErrInvalidUndoMoves = errors.New("invalid undo moves")
)

func wrapError(err error, message string) error {
	return errors.New(message + ": " + err.Error())
}
