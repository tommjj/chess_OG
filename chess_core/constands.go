// constants used in the chess

package chess_core

// Color constants

type Color int

const (
	White Color = iota
	Black
	Both
	None
)

func (c Color) Opposite() Color {
	switch c {
	case White:
		return Black
	case Black:
		return White
	case Both:
		return Both
	default:
		return None
	}
}

func (c Color) String() string {
	switch c {
	case White:
		return "White"
	case Black:
		return "Black"
	case Both:
		return "Both"
	default:
		return "None"
	}
}

type Square int

func (s Square) String() string {
	if !s.IsValid() {
		return "Invalid Square"
	}

	return SquareToCoordinates[s]
}

func (s Square) Rank() int {
	return int(s)/8 + 1
}

func (s Square) File() int {
	return int(s)%8 + 1
}

func (s Square) IsValid() bool {
	return s < 64 && s >= 0
}

func (s Square) ToBB() BitBoard {
	if !s.IsValid() {
		return BitBoard(0)
	}

	return BitBoard(1) << s
}

// Square constants
const (
	SquareA1 Square = iota
	SquareB1
	SquareC1
	SquareD1
	SquareE1
	SquareF1
	SquareG1
	SquareH1
	SquareA2
	SquareB2
	SquareC2
	SquareD2
	SquareE2
	SquareF2
	SquareG2
	SquareH2
	SquareA3
	SquareB3
	SquareC3
	SquareD3
	SquareE3
	SquareF3
	SquareG3
	SquareH3
	SquareA4
	SquareB4
	SquareC4
	SquareD4
	SquareE4
	SquareF4
	SquareG4
	SquareH4
	SquareA5
	SquareB5
	SquareC5
	SquareD5
	SquareE5
	SquareF5
	SquareG5
	SquareH5
	SquareA6
	SquareB6
	SquareC6
	SquareD6
	SquareE6
	SquareF6
	SquareG6
	SquareH6
	SquareA7
	SquareB7
	SquareC7
	SquareD7
	SquareE7
	SquareF7
	SquareG7
	SquareH7
	SquareA8
	SquareB8
	SquareC8
	SquareD8
	SquareE8
	SquareF8
	SquareG8
	SquareH8
)

// Piece constants

type PieceType byte

const (
	Pawn PieceType = iota + 1
	Knight
	Bishop
	Rook
	Queen
	King
)

type Piece byte

// encode pieces
const (
	WPawn   Piece = 'P'
	WKnight Piece = 'N'
	WBishop Piece = 'B'
	WRook   Piece = 'R'
	WQueen  Piece = 'Q'
	WKing   Piece = 'K'

	BPawn   Piece = 'p'
	BKnight Piece = 'n'
	BBishop Piece = 'b'
	BRook   Piece = 'r'
	BQueen  Piece = 'q'
	BKing   Piece = 'k'

	Empty Piece = '.'
)

func (p Piece) Color() Color {
	return getSideByPiece(p)
}

func (p Piece) Type() PieceType {
	return getPieceType(p)
}

func isWhitePiece(p Piece) bool {
	return p >= 'A' && p <= 'Z'
}

func isBlackPiece(p Piece) bool {
	return p >= 'a' && p <= 'z'
}

func getSideByPiece(p Piece) Color {
	if p == Empty {
		return None
	}

	if isBlackPiece(p) {
		return Black
	} else {
		return White
	}
}

func getPieceBySide(p PieceType, side Color) Piece {
	switch side {
	case White:
		return ASCIIPieces[p]
	case Black:
		return ASCIIPieces[p+6]
	default:
		panic("getPieceBySide side must be back or white")
	}
}

func getPieceType(p Piece) PieceType {
	switch p {
	case WPawn, BPawn:
		return Pawn
	case WBishop, BBishop:
		return Bishop
	case WRook, BRook:
		return Rook
	case WKnight, BKnight:
		return Knight
	case WQueen, BQueen:
		return Queen
	case WKing, BKing:
		return King
	default:
		panic("invalid piece")
	}
}

func IsValidPiece(p Piece) bool {
	return isWhitePiece(p) || isBlackPiece(p)
}

var ASCIIPieces = [13]Piece{' ', // None
	WPawn, WKnight, WBishop, WRook, WQueen, WKing,
	BPawn, BKnight, BBishop, BRook, BQueen, BKing,
}

func encodedPieceIndex(p Piece) int {
	switch p {
	case WPawn:
		return 1
	case WKnight:
		return 2
	case WBishop:
		return 3
	case WRook:
		return 4
	case WQueen:
		return 5
	case WKing:
		return 6
	case BPawn:
		return 7
	case BKnight:
		return 8
	case BBishop:
		return 9
	case BRook:
		return 10
	case BQueen:
		return 11
	case BKing:
		return 12
	}
	return 0
}

/*
castling permissions

	0001    1  white king can castle to the king side
	0010    2  white king can castle to the queen side
	0100    4  black king can castle to the king side
	1000    8  black king can castle to the queen side

- examples

	1111       both sides an castle both directions
	1001       black king => queen side

	white king => king side
*/
const (
	WK = 1
	WQ = 2
	BK = 4
	BQ = 8
)
const AllCastling = WK | WQ | BK | BQ

// castling rights update constants
var CastlingRights = [64]int{
	7, 15, 15, 15, 3, 15, 15, 11,
	15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15,
	13, 15, 15, 15, 12, 15, 15, 14,
}

// board squares to coordinates
var SquareToCoordinates = [64]string{
	"a1", "b1", "c1", "d1", "e1", "f1", "g1", "h1",
	"a2", "b2", "c2", "d2", "e2", "f2", "g2", "h2",
	"a3", "b3", "c3", "d3", "e3", "f3", "g3", "h3",
	"a4", "b4", "c4", "d4", "e4", "f4", "g4", "h4",
	"a5", "b5", "c5", "d5", "e5", "f5", "g5", "h5",
	"a6", "b6", "c6", "d6", "e6", "f6", "g6", "h6",
	"a7", "b7", "c7", "d7", "e7", "f7", "g7", "h7",
	"a8", "b8", "c8", "d8", "e8", "f8", "g8", "h8",
}

// ASCII pieces
var PieceToASCII = map[Piece]string{
	BPawn:   "♟",
	BKnight: "♞",
	BBishop: "♝",
	BRook:   "♜",
	BQueen:  "♛",
	BKing:   "♚",

	WPawn:   "♙",
	WKnight: "♘",
	WBishop: "♗",
	WRook:   "♖",
	WQueen:  "♕",
	WKing:   "♔",

	Empty: ".",
}

var PromotedPieces = map[Piece]byte{
	BQueen:  'q',
	BRook:   'r',
	BBishop: 'b',
	BKnight: 'n',

	WQueen:  'q',
	WRook:   'r',
	WBishop: 'b',
	WKnight: 'n',
}

var SquareToPosition = [64][2]int{
	0: {7, 0}, 1: {7, 1}, 2: {7, 2}, 3: {7, 3}, 4: {7, 4}, 5: {7, 5}, 6: {7, 6}, 7: {7, 7},
	8: {6, 0}, 9: {6, 1}, 10: {6, 2}, 11: {6, 3}, 12: {6, 4}, 13: {6, 5}, 14: {6, 6}, 15: {6, 7},
	16: {5, 0}, 17: {5, 1}, 18: {5, 2}, 19: {5, 3}, 20: {5, 4}, 21: {5, 5}, 22: {5, 6}, 23: {5, 7},
	24: {4, 0}, 25: {4, 1}, 26: {4, 2}, 27: {4, 3}, 28: {4, 4}, 29: {4, 5}, 30: {4, 6}, 31: {4, 7},
	32: {3, 0}, 33: {3, 1}, 34: {3, 2}, 35: {3, 3}, 36: {3, 4}, 37: {3, 5}, 38: {3, 6}, 39: {3, 7},
	40: {2, 0}, 41: {2, 1}, 42: {2, 2}, 43: {2, 3}, 44: {2, 4}, 45: {2, 5}, 46: {2, 6}, 47: {2, 7},
	48: {1, 0}, 49: {1, 1}, 50: {1, 2}, 51: {1, 3}, 52: {1, 4}, 53: {1, 5}, 54: {1, 6}, 55: {1, 7},
	56: {0, 0}, 57: {0, 1}, 58: {0, 2}, 59: {0, 3}, 60: {0, 4}, 61: {0, 5}, 62: {0, 6}, 63: {0, 7},
}

const NoEnPassant Square = 64

const StartingFEN string = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

const EmptyFEN string = "8/8/8/8/8/8/8/8 w - - 0 1"

// GameResult
type GameResult int

const (
	ResultOngoing GameResult = iota
	ResultCheckmate
	ResultStalemate
	ResultDrawBy50Move
	ResultInsufficientMaterial
)

func (g GameResult) String() string {
	switch g {
	case ResultOngoing:
		return "Result Ongoing"
	case ResultCheckmate:
		return "Result Checkmate"
	case ResultStalemate:
		return "Result Stalemate"
	case ResultDrawBy50Move:
		return "Result Draw By 50 Move"
	case ResultInsufficientMaterial:
		return "Result Insufficient Material"
	default:
		return "Unknown Result"
	}
}
