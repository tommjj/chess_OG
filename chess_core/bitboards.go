// bit boards and game state representation
//
// references:
// - https://www.chessprogramming.org/Bitboards
// - https://www.chessprogramming.org/Magic_Bitboards
// - https://www.chessprogramming.org/Forsyth-Edwards_Notation
// - https://www.chessprogramming.org/Castling
// - https://www.chessprogramming.org/En_passant
// - https://www.chessprogramming.org/Perft_Results
// - implementation: https://github.com/maksimKorzh/bbc

package chess_core

import (
	"fmt"
	"math/bits"
	"strings"
)

// BitBoard represents a 64-bit bitboard
// Each bit corresponds to a square on the chessboard
// with the least significant bit (LSB) representing a1
// and the most significant bit (MSB) representing h8
type BitBoard uint64

func (bb BitBoard) IsSet(square Square) bool {
	return (bb & (1 << square)) != 0
}

func (bb *BitBoard) Set(square Square) {
	*bb |= (1 << square)
}

func (bb *BitBoard) Clear(square Square) {
	*bb &^= (1 << square)
}

func (bb *BitBoard) Toggle(square Square) {
	*bb ^= (1 << square)
}

func (bb BitBoard) Count() int {
	return bits.OnesCount64(uint64(bb))
}

func (bb BitBoard) IsEmpty() bool {
	return bb == 0
}

func (bb BitBoard) LeastSignificantBit() int {
	if bb == 0 {
		return -1
	}
	return bits.TrailingZeros64(uint64(bb))
}

func (bb BitBoard) MostSignificantBit() int {
	if bb == 0 {
		return -1
	}
	return 63 - bits.LeadingZeros64(uint64(bb))
}

func (bb BitBoard) String() string {
	var board strings.Builder
	for rank := 7; rank >= 0; rank-- {
		for file := 0; file < 8; file++ {
			i := rank*8 + file
			if (bb & (1 << i)) != 0 {
				board.WriteRune('1')
			} else {
				board.WriteRune('.')
			}
			if file != 7 {
				board.WriteRune(' ')
			}
		}
		board.WriteRune('\n')
	}
	return board.String()
}

type BitBoards struct {
	WhitePawns   BitBoard
	WhiteKnights BitBoard
	WhiteBishops BitBoard
	WhiteRooks   BitBoard
	WhiteQueens  BitBoard
	WhiteKing    BitBoard

	BlackPawns   BitBoard
	BlackKnights BitBoard
	BlackBishops BitBoard
	BlackRooks   BitBoard
	BlackQueens  BitBoard
	BlackKing    BitBoard

	WhitePieces BitBoard
	BlackPieces BitBoard
	AllPieces   BitBoard
}

func (bbs *BitBoards) String() string {
	var board strings.Builder

	for rank := 7; rank >= 0; rank-- {
		// In số hàng (rank + 1)
		board.WriteString(fmt.Sprintf("%d  ", rank+1))

		for file := 0; file < 8; file++ {
			i := Square(rank*8 + file)
			squarePiece := Empty

			if bbs.AllPieces.IsSet(i) {
				if bbs.WhitePieces.IsSet(i) {
					switch {
					case bbs.WhitePawns.IsSet(i):
						squarePiece = WPawn
					case bbs.WhiteKnights.IsSet(i):
						squarePiece = WKnight
					case bbs.WhiteBishops.IsSet(i):
						squarePiece = WBishop
					case bbs.WhiteRooks.IsSet(i):
						squarePiece = WRook
					case bbs.WhiteQueens.IsSet(i):
						squarePiece = WQueen
					case bbs.WhiteKing.IsSet(i):
						squarePiece = WKing
					}
				} else if bbs.BlackPieces.IsSet(i) {
					switch {
					case bbs.BlackPawns.IsSet(i):
						squarePiece = BPawn
					case bbs.BlackKnights.IsSet(i):
						squarePiece = BKnight
					case bbs.BlackBishops.IsSet(i):
						squarePiece = BBishop
					case bbs.BlackRooks.IsSet(i):
						squarePiece = BRook
					case bbs.BlackQueens.IsSet(i):
						squarePiece = BQueen
					case bbs.BlackKing.IsSet(i):
						squarePiece = BKing
					}
				}
			}

			board.WriteString(PieceToASCII[squarePiece])
			if file != 7 {
				board.WriteString(" ")
			}
		}
		board.WriteRune('\n')
	}

	// Thêm tọa độ cột
	board.WriteString("\n   a b c d e f g h\n")

	return board.String()
}

func (bbs *BitBoards) Clear() {
	bbs.WhitePawns = 0
	bbs.WhiteKnights = 0
	bbs.WhiteBishops = 0
	bbs.WhiteRooks = 0
	bbs.WhiteQueens = 0
	bbs.WhiteKing = 0

	bbs.BlackPawns = 0
	bbs.BlackKnights = 0
	bbs.BlackBishops = 0
	bbs.BlackRooks = 0
	bbs.BlackQueens = 0
	bbs.BlackKing = 0

	bbs.WhitePieces = 0
	bbs.BlackPieces = 0
	bbs.AllPieces = 0
}

func (bbs *BitBoards) UpdateAggregate() {
	bbs.WhitePieces = bbs.WhitePawns | bbs.WhiteKnights | bbs.WhiteBishops | bbs.WhiteRooks | bbs.WhiteQueens | bbs.WhiteKing
	bbs.BlackPieces = bbs.BlackPawns | bbs.BlackKnights | bbs.BlackBishops | bbs.BlackRooks | bbs.BlackQueens | bbs.BlackKing
	bbs.AllPieces = bbs.WhitePieces | bbs.BlackPieces
}

func (bbs *BitBoards) Copy() *BitBoards {
	return &BitBoards{
		WhitePawns:   bbs.WhitePawns,
		WhiteKnights: bbs.WhiteKnights,
		WhiteBishops: bbs.WhiteBishops,
		WhiteRooks:   bbs.WhiteRooks,
		WhiteQueens:  bbs.WhiteQueens,
		WhiteKing:    bbs.WhiteKing,
		BlackPawns:   bbs.BlackPawns,
		BlackKnights: bbs.BlackKnights,
		BlackBishops: bbs.BlackBishops,
		BlackRooks:   bbs.BlackRooks,
		BlackQueens:  bbs.BlackQueens,
		BlackKing:    bbs.BlackKing,
		WhitePieces:  bbs.WhitePieces,
		BlackPieces:  bbs.BlackPieces,
		AllPieces:    bbs.AllPieces,
	}
}

func (bbs *BitBoards) PiecesByType(color Color, pieceType PieceType) BitBoard {
	switch color {
	case White:
		switch pieceType {
		case Pawn:
			return bbs.WhitePawns
		case Knight:
			return bbs.WhiteKnights
		case Bishop:
			return bbs.WhiteBishops
		case Rook:
			return bbs.WhiteRooks
		case Queen:
			return bbs.WhiteQueens
		case King:
			return bbs.WhiteKing
		default:
			return 0
		}
	case Black:
		switch pieceType {
		case Pawn:
			return bbs.BlackPawns
		case Knight:
			return bbs.BlackKnights
		case Bishop:
			return bbs.BlackBishops
		case Rook:
			return bbs.BlackRooks
		case Queen:
			return bbs.BlackQueens
		case King:
			return bbs.BlackKing
		default:
			return 0
		}
	default:
		return 0
	}
}

func (bbs *BitBoards) Pieces(piece Piece) BitBoard {
	switch piece {
	case WPawn:
		return bbs.WhitePawns
	case WRook:
		return bbs.WhiteRooks
	case WBishop:
		return bbs.WhiteBishops
	case WKnight:
		return bbs.WhiteKnights
	case WQueen:
		return bbs.WhiteQueens
	case WKing:
		return bbs.WhiteKing
	case BPawn:
		return bbs.BlackPawns
	case BRook:
		return bbs.BlackRooks
	case BBishop:
		return bbs.BlackBishops
	case BKnight:
		return bbs.BlackKnights
	case BQueen:
		return bbs.BlackQueens
	case BKing:
		return bbs.BlackKing
	default:
		return 0
	}
}

func (bbs *BitBoards) setPieces(piece Piece, bb BitBoard) {
	switch piece {
	case WPawn:
		bbs.WhitePawns = bb
	case WRook:
		bbs.WhiteRooks = bb
	case WBishop:
		bbs.WhiteBishops = bb
	case WKnight:
		bbs.WhiteKnights = bb
	case WQueen:
		bbs.WhiteQueens = bb
	case WKing:
		bbs.WhiteKing = bb
	case BPawn:
		bbs.BlackPawns = bb
	case BRook:
		bbs.BlackRooks = bb
	case BBishop:
		bbs.BlackBishops = bb
	case BKnight:
		bbs.BlackKnights = bb
	case BQueen:
		bbs.BlackQueens = bb
	case BKing:
		bbs.BlackKing = bb
	}
}

func (bbs *BitBoards) SetPieces(piece Piece, bb BitBoard) {
	bbs.setPieces(piece, bb)
	bbs.UpdateAggregate()
}

func (bbs *BitBoards) GetPieceAt(square Square) Piece {
	if !bbs.AllPieces.IsSet(square) {
		return Empty
	}

	if bbs.WhitePieces.IsSet(square) {
		switch {
		case bbs.WhitePawns.IsSet(square):
			return WPawn
		case bbs.WhiteKnights.IsSet(square):
			return WKnight
		case bbs.WhiteBishops.IsSet(square):
			return WBishop
		case bbs.WhiteRooks.IsSet(square):
			return WRook
		case bbs.WhiteQueens.IsSet(square):
			return WQueen
		case bbs.WhiteKing.IsSet(square):
			return WKing
		}
	} else if bbs.BlackPieces.IsSet(square) {
		switch {
		case bbs.BlackPawns.IsSet(square):
			return BPawn
		case bbs.BlackKnights.IsSet(square):
			return BKnight
		case bbs.BlackBishops.IsSet(square):
			return BBishop
		case bbs.BlackRooks.IsSet(square):
			return BRook
		case bbs.BlackQueens.IsSet(square):
			return BQueen
		case bbs.BlackKing.IsSet(square):
			return BKing
		}
	}
	return Empty
}

func (bbs *BitBoards) clearSquare(square Square) {
	if bbs.WhitePieces.IsSet(square) {
		switch {
		case bbs.WhitePawns.IsSet(square):
			bbs.WhitePawns &^= (1 << square)
		case bbs.WhiteKnights.IsSet(square):
			bbs.WhiteKnights &^= (1 << square)
		case bbs.WhiteBishops.IsSet(square):
			bbs.WhiteBishops &^= (1 << square)
		case bbs.WhiteRooks.IsSet(square):
			bbs.WhiteRooks &^= (1 << square)
		case bbs.WhiteQueens.IsSet(square):
			bbs.WhiteQueens &^= (1 << square)
		case bbs.WhiteKing.IsSet(square):
			bbs.WhiteKing &^= (1 << square)
		}
	} else if bbs.BlackPieces.IsSet(square) {
		switch {
		case bbs.BlackPawns.IsSet(square):
			bbs.BlackPawns &^= (1 << square)
		case bbs.BlackKnights.IsSet(square):
			bbs.BlackKnights &^= (1 << square)
		case bbs.BlackBishops.IsSet(square):
			bbs.BlackBishops &^= (1 << square)
		case bbs.BlackRooks.IsSet(square):
			bbs.BlackRooks &^= (1 << square)
		case bbs.BlackQueens.IsSet(square):
			bbs.BlackQueens &^= (1 << square)
		case bbs.BlackKing.IsSet(square):
			bbs.BlackKing &^= (1 << square)
		}
	}
}

func (bbs *BitBoards) ClearSquare(square Square) {
	bbs.clearSquare(square)
	bbs.UpdateAggregate()
}

func (bbs *BitBoards) setPieceAt(square Square, piece Piece) {
	bbs.clearSquare(square)

	switch piece {
	case WPawn:
		bbs.WhitePawns.Set(square)
	case WKnight:
		bbs.WhiteKnights.Set(square)
	case WBishop:
		bbs.WhiteBishops.Set(square)
	case WRook:
		bbs.WhiteRooks.Set(square)
	case WQueen:
		bbs.WhiteQueens.Set(square)
	case WKing:
		bbs.WhiteKing.Set(square)
	case BPawn:
		bbs.BlackPawns.Set(square)
	case BKnight:
		bbs.BlackKnights.Set(square)
	case BBishop:
		bbs.BlackBishops.Set(square)
	case BRook:
		bbs.BlackRooks.Set(square)
	case BQueen:
		bbs.BlackQueens.Set(square)
	case BKing:
		bbs.BlackKing.Set(square)
	}
}

func (bbs *BitBoards) SetPieceAt(square Square, piece Piece) {
	bbs.setPieceAt(square, piece)
	bbs.UpdateAggregate()
}

func (bbs *BitBoards) GetColorAt(square Square) Color {
	if !bbs.AllPieces.IsSet(square) {
		return None
	}

	if bbs.WhitePieces.IsSet(square) {
		return White
	} else if bbs.BlackPieces.IsSet(square) {
		return Black
	}
	return None
}

func (bbs *BitBoards) OccupiedBy(side Color) BitBoard {
	switch side {
	case White:
		return bbs.WhitePieces
	case Black:
		return bbs.BlackPieces
	case Both:
		return bbs.AllPieces
	default:
		return 0
	}
}

func (bbs *BitBoards) IsSquareEmpty(square Square) bool {
	return !bbs.AllPieces.IsSet(square)
}

func (bbs *BitBoards) IsSquareOccupiedByColor(square Square, color Color) bool {
	if !bbs.AllPieces.IsSet(square) {
		return false
	}

	if color == White && bbs.WhitePieces.IsSet(square) {
		return true
	} else if color == Black && bbs.BlackPieces.IsSet(square) {
		return true
	}
	return false
}

func (bbs *BitBoards) IsSquareOccupiedByOpponent(square Square, color Color) bool {
	if !bbs.AllPieces.IsSet(square) {
		return false
	}

	if color == White && bbs.BlackPieces.IsSet(square) {
		return true
	} else if color == Black && bbs.WhitePieces.IsSet(square) {
		return true
	}
	return false
}

func (bbs *BitBoards) CountPieces() int {
	return bbs.AllPieces.Count()
}

func (bbs *BitBoards) CountPiecesByColor(color Color) int {
	switch color {
	case White:
		return bbs.WhitePieces.Count()
	case Black:
		return bbs.BlackPieces.Count()
	case Both:
		return bbs.AllPieces.Count()
	default:
		return 0
	}
}

func (bbs *BitBoards) CountPiecesByType(piece Piece) int {
	count := 0
	var bitboard BitBoard

	switch piece {
	case WPawn:
		bitboard = bbs.WhitePawns
	case WKnight:
		bitboard = bbs.WhiteKnights
	case WBishop:
		bitboard = bbs.WhiteBishops
	case WRook:
		bitboard = bbs.WhiteRooks
	case WQueen:
		bitboard = bbs.WhiteQueens
	case WKing:
		bitboard = bbs.WhiteKing
	case BPawn:
		bitboard = bbs.BlackPawns
	case BKnight:
		bitboard = bbs.BlackKnights
	case BBishop:
		bitboard = bbs.BlackBishops
	case BRook:
		bitboard = bbs.BlackRooks
	case BQueen:
		bitboard = bbs.BlackQueens
	case BKing:
		bitboard = bbs.BlackKing
	default:
		return 0
	}

	for i := range Square(64) {
		if bitboard.IsSet(i) {
			count++
		}
	}
	return count
}

func (bbs *BitBoards) FromFEN(fen string) error {
	bbs.Clear()

	rank := 7
	file := 0

	for _, char := range fen {
		if char == ' ' {
			break
		}

		if char == '/' {
			rank--
			file = 0
			continue
		}

		if char >= '1' && char <= '8' {
			file += int(char - '0')
			continue
		}

		if file > 7 || rank < 0 {
			return wrapError(ErrInvalidFEN, "FromFEN")
		}

		square := rank*8 + file
		// validate piece
		if ok := IsValidPiece(Piece(char)); !ok {
			return wrapError(ErrInvalidPiece, "FromFEN")
		}

		bbs.setPieceAt(Square(square), Piece(char))
		file++
	}

	bbs.UpdateAggregate()
	return nil
}

func (bbs *BitBoards) moveUnsafePiece(piece Piece, from, to Square) {
	pieces := bbs.Pieces(piece)
	moved := (pieces & ^from.ToBB()) | to.ToBB()
	bbs.setPieces(piece, moved)
}

func (bbs *BitBoards) MoveUnsafePiece(piece Piece, from, to Square) {
	bbs.moveUnsafePiece(piece, from, to)
	bbs.UpdateAggregate()
}

func (bbs *BitBoards) movePiece(piece Piece, from, to Square) error {
	pieces := bbs.Pieces(piece)
	side := getSideByPiece(piece)
	fromBB, toBB := from.ToBB(), to.ToBB()

	if pieces&fromBB == 0 {
		return ErrInvalidMove
	}

	enemyColor := side.Opposite()
	switch bbs.GetColorAt(to) {
	case enemyColor:
		bbs.clearSquare(to)
	case side:
		return ErrInvalidMove
	}

	moved := (pieces & ^fromBB) | toBB
	bbs.setPieces(piece, moved)

	return nil
}

func (bbs *BitBoards) MovePiece(piece Piece, from, to Square) error {
	if err := bbs.movePiece(piece, from, to); err != nil {
		return err
	}
	bbs.UpdateAggregate()
	return nil
}
