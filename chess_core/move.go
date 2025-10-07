// move encoding and move list management

package main

import (
	"fmt"
	"sync"
)

// Move encoding
/*
         binary move bits                               hexidecimal constants

   0000 0000 0000 0000 0011 1111    source square       0x3f
   0000 0000 0000 1111 1100 0000    target square       0xfc0
   0000 0000 1111 0000 0000 0000    piece               0xf000
   0000 1111 0000 0000 0000 0000    promoted piece      0xf0000
   0001 0000 0000 0000 0000 0000    capture flag        0x100000
   0010 0000 0000 0000 0000 0000    double push flag    0x200000
   0100 0000 0000 0000 0000 0000    enpassant flag      0x400000
   1000 0000 0000 0000 0000 0000    castling flag       0x800000

         piece encoding PNBRQKpnbrqk

	0000 WPawn   0
	0001 WKnight 1
	0010 WBishop 2
	0011 WRook   3
	0100 WQueen  4
	0101 WKing   5

	0110 BPawn   6
	0111 BKnight 7
	1000 BBishop 8
	1001 BRook   9
	1010 BQueen  10
	1011 BKing   11
*/
type Move uint32

// NewMove create a new encoding move
/*
	src: 		source square
	dst: 		target square
	piece: 		encoding piece ex:(encodedPieceIndex(pieceToMove))
	promo: 		promotion piece type ex:uint32(Queen)
	capture: 	capture flag 1 or 0
	dbl: 		double push flag 1 or 0
	ep: 		enpassant flag 1 or 0
	castle: 	castling flag 1 or 0
*/
func NewMove(src, dst, piece, promo, capture, dbl, ep, castle uint32) Move {
	return Move(
		(src) |
			(dst << 6) |
			(piece << 12) |
			(promo << 16) |
			(capture << 20) |
			(dbl << 21) |
			(ep << 22) |
			(castle << 23),
	)
}

func (m Move) From() uint32     { return uint32(m) & 0x3F }
func (m Move) To() uint32       { return (uint32(m) >> 6) & 0x3F }
func (m Move) Piece() uint32    { return (uint32(m) >> 12) & 0xF }
func (m Move) Promoted() uint32 { return (uint32(m) >> 16) & 0xF }

func (m Move) IsPromotion() bool  { return m.Promoted() != 0 }
func (m Move) IsCapture() bool    { return (uint32(m) & 0x100000) != 0 }
func (m Move) IsDoublePush() bool { return (uint32(m) & 0x200000) != 0 }
func (m Move) IsEnPassant() bool  { return (uint32(m) & 0x400000) != 0 }
func (m Move) IsCastle() bool     { return (uint32(m) & 0x800000) != 0 }

func (m Move) Side() Color {
	return getSideByPiece(ASCIIPieces[m.Piece()])
}

func (m Move) String() string {
	if m.IsPromotion() {
		return fmt.Sprintf(
			"from=%s, to=%s, promoted=%s, capture=%v",
			SquareToCoordinates[m.From()],
			SquareToCoordinates[m.To()],
			PieceToASCII[ASCIIPieces[m.Promoted()]],
			m.IsCapture(),
		)
	}
	return fmt.Sprintf(
		"from=%s, to=%s, piece=%s, capture=%v",
		SquareToCoordinates[m.From()],
		SquareToCoordinates[m.To()],
		PieceToASCII[ASCIIPieces[m.Piece()]],
		m.IsCapture(),
	)
}

type MoveList []Move

func (ml *MoveList) Add(move Move) {
	*ml = append(*ml, move)
}

func (ml *MoveList) Clear() {
	*ml = (*ml)[:0]
}

func (ml *MoveList) Len() int {
	return len(*ml)
}

func (ml *MoveList) Get(index int) Move {
	return (*ml)[index]
}

func (ml *MoveList) Set(index int, move Move) {
	(*ml)[index] = move
}

// moveListPool is a sync.Pool to reuse MoveList slices
var moveListPool = sync.Pool{
	New: func() any {
		return make(MoveList, 0, 220) // average max moves in a position is around 218
	},
}

func generatePseudoLegalMoves(bb *BitBoards, side Color, castling int, enPassantSquare Square, ml *MoveList) {
	generatePawnMoves(bb, side, enPassantSquare, ml)
	generateBishopMoves(bb, side, ml)
	generateRookMoves(bb, side, ml)
	generateQueenMoves(bb, side, ml)
	generateKingMoves(bb, side, castling, ml)
}

// // generateLegalMove generates all legal moves for the given position and side to move
// func generateLegalMove(bb BitBoards, side Color, castling uint8, enPassantSquare Square) MoveList {
// 	ml := moveListPool.Get().(MoveList)
// 	ml.Clear()

// 	// Generate pseudo-legal moves
// 	generatePseudoLegalMoves(&bb, side, castling, enPassantSquare, ml)

// 	// Filter out illegal moves
// 	legalMoves := ml[:0]

// generatePawnMoves
func generatePawnMoves(bb *BitBoards, side Color, enPassantSquare Square, ml *MoveList) {
	var pawnPieces, enemyPieces BitBoard
	var forwardStep, startRank, promotionRank int
	var leftOffset, rightOffset int // Captures offset

	var pawnTypeCode uint32

	if side == White {
		leftOffset = 1
		rightOffset = -1
		pawnTypeCode = uint32(encodedPieceIndex(WPawn))
		pawnPieces = bb.WhitePawns
		enemyPieces = bb.BlackPieces
		forwardStep = 8
		startRank = 1     // rank 2
		promotionRank = 8 // rank 8
	} else {
		leftOffset = -1
		rightOffset = 1
		pawnTypeCode = uint32(encodedPieceIndex(BPawn))
		pawnPieces = bb.BlackPawns
		enemyPieces = bb.WhitePieces
		forwardStep = -8
		startRank = 6     // rank 7
		promotionRank = 1 // rank 1
	}

	// Single and double forward moves
	var singlePush, doublePush BitBoard
	if side == White {
		singlePush = (pawnPieces << 8) & ^bb.AllPieces
		doublePush = ((pawnPieces & rankMask(startRank)) << 16) & ^bb.AllPieces & ^(bb.AllPieces << 8)
	} else {
		singlePush = (pawnPieces >> 8) & ^bb.AllPieces
		doublePush = ((pawnPieces & rankMask(startRank)) >> 16) & ^bb.AllPieces & ^(bb.AllPieces >> 8)
	}

	for sq := singlePush; sq != 0; {
		toSquare := popLSB(&sq)
		fromSquare := Square(int(toSquare) - forwardStep)

		// Check for promotion
		if toSquare.Rank() == promotionRank {
			ml.Add(NewMove(uint32(fromSquare), uint32(toSquare), pawnTypeCode, uint32(Queen), 0, 0, 0, 0))
			ml.Add(NewMove(uint32(fromSquare), uint32(toSquare), pawnTypeCode, uint32(Rook), 0, 0, 0, 0))
			ml.Add(NewMove(uint32(fromSquare), uint32(toSquare), pawnTypeCode, uint32(Bishop), 0, 0, 0, 0))
			ml.Add(NewMove(uint32(fromSquare), uint32(toSquare), pawnTypeCode, uint32(Knight), 0, 0, 0, 0))
		} else {
			ml.Add(NewMove(uint32(fromSquare), uint32(toSquare), pawnTypeCode, 0, 0, 0, 0, 0))
		}
	}

	// Double forward moves
	for sq := doublePush; sq != 0; {
		toSquare := popLSB(&sq)
		fromSquare := Square(int(toSquare) - 2*forwardStep)
		ml.Add(NewMove(uint32(fromSquare), uint32(toSquare), pawnTypeCode, 0, 0, 1, 0, 0))
	}

	// Captures
	var leftCaptures, rightCaptures BitBoard
	if side == White {
		leftCaptures = ((pawnPieces &^ fileMask(0)) << 7) & enemyPieces
		rightCaptures = ((pawnPieces &^ fileMask(7)) << 9) & enemyPieces
	} else {
		leftCaptures = ((pawnPieces &^ fileMask(7)) >> 7) & enemyPieces
		rightCaptures = ((pawnPieces &^ fileMask(0)) >> 9) & enemyPieces
	}

	for sq := leftCaptures; sq != 0; {
		toSquare := popLSB(&sq)
		fromSquare := Square(int(toSquare) - forwardStep + leftOffset)
		// Check for promotion
		if toSquare.Rank() == promotionRank {
			ml.Add(NewMove(uint32(fromSquare), uint32(toSquare), pawnTypeCode, uint32(Queen), 0, 0, 0, 0))
			ml.Add(NewMove(uint32(fromSquare), uint32(toSquare), pawnTypeCode, uint32(Rook), 0, 0, 0, 0))
			ml.Add(NewMove(uint32(fromSquare), uint32(toSquare), pawnTypeCode, uint32(Bishop), 0, 0, 0, 0))
			ml.Add(NewMove(uint32(fromSquare), uint32(toSquare), pawnTypeCode, uint32(Knight), 0, 0, 0, 0))
		} else {
			ml.Add(NewMove(uint32(fromSquare), uint32(toSquare), pawnTypeCode, 0, 1, 0, 0, 0))
		}
	}

	for sq := rightCaptures; sq != 0; {
		toSquare := popLSB(&sq)
		fromSquare := Square(int(toSquare) - forwardStep + rightOffset)
		// Check for promotion
		if toSquare.Rank() == promotionRank {
			ml.Add(NewMove(uint32(fromSquare), uint32(toSquare), pawnTypeCode, uint32(Queen), 0, 0, 0, 0))
			ml.Add(NewMove(uint32(fromSquare), uint32(toSquare), pawnTypeCode, uint32(Rook), 0, 0, 0, 0))
			ml.Add(NewMove(uint32(fromSquare), uint32(toSquare), pawnTypeCode, uint32(Bishop), 0, 0, 0, 0))
			ml.Add(NewMove(uint32(fromSquare), uint32(toSquare), pawnTypeCode, uint32(Knight), 0, 0, 0, 0))
		} else {
			ml.Add(NewMove(uint32(fromSquare), uint32(toSquare), pawnTypeCode, 0, 1, 0, 0, 0))
		}
	}

	// En passant captures
	if enPassantSquare != NoEnPassant {
		var leftCaptures, rightCaptures BitBoard
		if side == White {
			leftCaptures = ((pawnPieces &^ fileMask(0)) << 7) & enPassantSquare.ToBB()  // left capture
			rightCaptures = ((pawnPieces &^ fileMask(7)) << 9) & enPassantSquare.ToBB() // right capture
		} else {
			leftCaptures = ((pawnPieces &^ fileMask(7)) >> 7) & enPassantSquare.ToBB()  // left capture
			rightCaptures = ((pawnPieces &^ fileMask(0)) >> 9) & enPassantSquare.ToBB() // right capture
		}

		for sq := leftCaptures; sq != 0; {
			toSquare := popLSB(&sq)
			fromSquare := Square(int(toSquare) - forwardStep + leftOffset)
			ml.Add(NewMove(uint32(fromSquare), uint32(toSquare), pawnTypeCode, 0, 1, 0, 1, 0))
		}

		for sq := rightCaptures; sq != 0; {
			toSquare := popLSB(&sq)
			fromSquare := Square(int(toSquare) - forwardStep + rightOffset)
			ml.Add(NewMove(uint32(fromSquare), uint32(toSquare), pawnTypeCode, 0, 1, 0, 1, 0))
		}
	}
}

// generateRookMoves
func generateRookMoves(bb *BitBoards, side Color, ml *MoveList) {
	var rookPieces, enemyPieces, allyPieces BitBoard
	var pawnTypeCode uint32

	if side == White {
		pawnTypeCode = uint32(encodedPieceIndex(WRook))
		rookPieces = bb.WhiteRooks
		enemyPieces = bb.BlackPieces
		allyPieces = bb.WhitePieces
	} else {
		pawnTypeCode = uint32(encodedPieceIndex(BRook))
		rookPieces = bb.BlackRooks
		enemyPieces = bb.WhitePieces
		allyPieces = bb.BlackPieces
	}

	for rp := rookPieces; rp != 0; {
		fromSq := Square(popLSB(&rp))

		attacks := RookAttacks(fromSq, bb.AllPieces) & ^allyPieces

		for p := attacks; p != 0; {
			toSq := popLSB(&p)

			isCapture := (enemyPieces >> toSq) & 1
			ml.Add(NewMove(uint32(fromSq), uint32(toSq), pawnTypeCode, 0, uint32(isCapture), 0, 0, 0))
		}
	}
}

// generateBishopMoves
func generateBishopMoves(bb *BitBoards, side Color, ml *MoveList) {
	var bishopPieces, enemyPieces, allyPieces BitBoard
	var pawnTypeCode uint32

	if side == White {
		pawnTypeCode = uint32(encodedPieceIndex(WBishop))
		bishopPieces = bb.WhiteBishops
		enemyPieces = bb.BlackPieces
		allyPieces = bb.WhitePieces
	} else {
		pawnTypeCode = uint32(encodedPieceIndex(BBishop))
		bishopPieces = bb.BlackBishops
		enemyPieces = bb.WhitePieces
		allyPieces = bb.BlackPieces
	}

	for rp := bishopPieces; rp != 0; {
		fromSq := Square(popLSB(&rp))

		attacks := BishopAttacks(fromSq, bb.AllPieces) & ^allyPieces

		for p := attacks; p != 0; {
			toSq := popLSB(&p)

			isCapture := (enemyPieces >> toSq) & 1
			ml.Add(NewMove(uint32(fromSq), uint32(toSq), pawnTypeCode, 0, uint32(isCapture), 0, 0, 0))
		}
	}
}

// generateQueenMoves
func generateQueenMoves(bb *BitBoards, side Color, ml *MoveList) {
	var queenPieces, enemyPieces, allyPieces BitBoard
	var pawnTypeCode uint32

	if side == White {
		pawnTypeCode = uint32(encodedPieceIndex(WQueen))
		queenPieces = bb.WhiteQueens
		enemyPieces = bb.BlackPieces
		allyPieces = bb.WhitePieces
	} else {
		pawnTypeCode = uint32(encodedPieceIndex(BQueen))
		queenPieces = bb.BlackQueens
		enemyPieces = bb.WhitePieces
		allyPieces = bb.BlackPieces
	}

	for rp := queenPieces; rp != 0; {
		fromSq := Square(popLSB(&rp))

		attacks := QueenAttacks(fromSq, bb.AllPieces) & ^allyPieces

		for p := attacks; p != 0; {
			toSq := popLSB(&p)

			isCapture := (enemyPieces >> toSq) & 1
			ml.Add(NewMove(uint32(fromSq), uint32(toSq), pawnTypeCode, 0, uint32(isCapture), 0, 0, 0))
		}
	}
}

// generateKingMoves
func generateKingMoves(bb *BitBoards, side Color, castling int, ml *MoveList) {
	var kingPieces, enemyPieces, allyPieces BitBoard
	var pawnTypeCode uint32

	if side == White {
		pawnTypeCode = uint32(encodedPieceIndex(WKing))
		kingPieces = bb.WhiteKing
		enemyPieces = bb.BlackPieces
		allyPieces = bb.WhitePieces
	} else {
		pawnTypeCode = uint32(encodedPieceIndex(BKing))
		kingPieces = bb.BlackKing
		enemyPieces = bb.WhitePieces
		allyPieces = bb.BlackPieces
	}

	fromSq := Square(kingPieces.LeastSignificantBit())

	attacks := KingAttacks(fromSq) & ^allyPieces

	for p := attacks; p != 0; {
		toSq := popLSB(&p)

		isCapture := (enemyPieces >> toSq) & 1
		ml.Add(NewMove(uint32(fromSq), uint32(toSq), pawnTypeCode, 0, uint32(isCapture), 0, 0, 0))
	}

	// castling
	if side == White {
		// short castle (O-O)
		if castling&WK == WK &&
			(bb.AllPieces&(SquareF1.ToBB()|SquareG1.ToBB())) == 0 && // đường trống
			!IsAttacked(SquareE1, White, bb) &&
			!IsAttacked(SquareF1, White, bb) &&
			!IsAttacked(SquareG1, White, bb) {
			ml.Add(NewMove(uint32(fromSq), uint32(SquareG1), pawnTypeCode, 0, 0, 0, 0, 1))
		}

		// long castle (O-O-O)
		if castling&WQ == WQ &&
			(bb.AllPieces&(SquareB1.ToBB()|SquareC1.ToBB()|SquareD1.ToBB())) == 0 &&
			!IsAttacked(SquareE1, White, bb) &&
			!IsAttacked(SquareD1, White, bb) &&
			!IsAttacked(SquareC1, White, bb) {
			ml.Add(NewMove(uint32(fromSq), uint32(SquareC1), pawnTypeCode, 0, 0, 0, 0, 1))
		}
	} else {
		// short castle (O-O)
		if castling&BK == BK &&
			(bb.AllPieces&(SquareF8.ToBB()|SquareG8.ToBB())) == 0 && // đường trống
			!IsAttacked(SquareE8, Black, bb) &&
			!IsAttacked(SquareF8, Black, bb) &&
			!IsAttacked(SquareG8, Black, bb) {
			ml.Add(NewMove(uint32(fromSq), uint32(SquareG8), pawnTypeCode, 0, 0, 0, 0, 1))
		}

		// long castle (O-O-O)
		if castling&BQ == BQ &&
			(bb.AllPieces&(SquareB8.ToBB()|SquareC8.ToBB()|SquareD8.ToBB())) == 0 &&
			!IsAttacked(SquareE8, Black, bb) &&
			!IsAttacked(SquareD8, Black, bb) &&
			!IsAttacked(SquareC8, Black, bb) {
			ml.Add(NewMove(uint32(fromSq), uint32(SquareC8), pawnTypeCode, 0, 0, 0, 0, 1))
		}
	}
}

// makeUnsafeMove don't check any rule when do move
func makeUnsafeMove(bbs *BitBoards, move Move) {
	pieceToMove := ASCIIPieces[move.Piece()]
	side := getSideByPiece(pieceToMove)
	from := Square(move.From())
	to := Square(move.To())

	switch {
	case move.IsCastle(): //
		switch to {
		case SquareG1: // white short castle
			bbs.moveUnsafePiece(WRook, SquareH1, SquareF1)
		case SquareC1: // white long castle
			bbs.moveUnsafePiece(WRook, SquareA1, SquareD1)
		case SquareG8: // black short castle
			bbs.moveUnsafePiece(BRook, SquareH8, SquareF8)
		case SquareC8: // black long castle
			bbs.moveUnsafePiece(BRook, SquareA8, SquareD8)
		}
	case move.IsPromotion():
		if move.IsCapture() {
			bbs.clearSquare(to)
		}
		pieceToMove = getPieceBySide(PieceType(move.Promoted()), side)
		bbs.setPieceAt(from, pieceToMove)
	case move.IsEnPassant():
		if side == White {
			bbs.clearSquare(to - 8)
		} else {
			bbs.clearSquare(to + 8)
		}
	case move.IsCapture():
		bbs.clearSquare(to)
	}

	bbs.moveUnsafePiece(pieceToMove, from, to)
	bbs.UpdateAggregate()
}

func hasAnyLegalMove(bb *BitBoards, side Color, castlingRights int, enPassant Square) bool {
	enemy := side.Opposite()
	allyPieces := bb.OccupiedBy(side)
	enemyPieces := bb.OccupiedBy(enemy)
	allPieces := bb.AllPieces

	kingSq := Square(bb.PiecesByType(side, King).LeastSignificantBit())
	if !kingSq.IsValid() {
		panic(ErrNoKing)
	}

	attackers := AttackersTo(kingSq, side, bb)

	// double chess
	if attackers.Count() >= 2 {
		kingMoves := KingAttacks(kingSq) & ^allyPieces
		for kingMoves != 0 {
			to := popLSB(&kingMoves)
			bbCopy := bb.Copy()
			move := NewMove(uint32(kingSq), uint32(to), uint32(encodedPieceIndex(getPieceBySide(King, side))), 0, 0, 0, 0, 0)
			makeUnsafeMove(bbCopy, move)
			if !IsKingAttacked(side, bbCopy) {
				return true
			}
		}
		return false
	}

	blockMask := BitBoard(0)
	if a := attackers; a != 0 {
		attSq := popLSB(&a)
		blockMask = BetweenBits(kingSq, attSq)
	}
	blockMask |= attackers

	for _, pieceType := range []PieceType{Pawn, Knight, Bishop, Rook, Queen, King} {
		pieces := bb.PiecesByType(side, pieceType)

		for pieces != 0 {
			from := popLSB(&pieces)
			var moves BitBoard

			switch pieceType {
			case Pawn:
				moves = PawnAttacks(from, side) & (enemyPieces | enPassant.ToBB())
				if side == White {
					oneStep := (from.ToBB() << 8) & ^allPieces
					moves |= oneStep

					if from.Rank() == 2 {
						twoStep := (oneStep << 8) & ^allPieces
						moves |= twoStep
					}
				} else {
					oneStep := (from.ToBB() >> 8) & ^allPieces
					moves |= oneStep

					if from.Rank() == 7 {
						twoStep := (oneStep >> 8) & ^allPieces
						moves |= twoStep
					}
				}
			case Knight:
				moves = KnightAttacks(from) & ^allyPieces
			case Bishop:
				moves = BishopAttacks(from, allPieces) & ^allyPieces
			case Rook:
				moves = RookAttacks(from, allPieces) & ^allyPieces
			case Queen:
				moves = QueenAttacks(from, allPieces) & ^allyPieces
			case King:
				moves = KingAttacks(from) & ^allyPieces
				// kiểm tra nhập thành
				// castling
				if side == White {
					// short castle (O-O)
					if castlingRights&WK == WK &&
						(bb.AllPieces&(SquareF1.ToBB()|SquareG1.ToBB())) == 0 && // đường trống
						!IsAttacked(SquareE1, White, bb) &&
						!IsAttacked(SquareF1, White, bb) &&
						!IsAttacked(SquareG1, White, bb) {
						moves |= SquareG1.ToBB()
					}

					// long castle (O-O-O)
					if castlingRights&WQ == WQ &&
						(bb.AllPieces&(SquareB1.ToBB()|SquareC1.ToBB()|SquareD1.ToBB())) == 0 &&
						!IsAttacked(SquareE1, White, bb) &&
						!IsAttacked(SquareD1, White, bb) &&
						!IsAttacked(SquareC1, White, bb) {
						moves |= SquareD1.ToBB()
					}
				} else {
					// short castle (O-O)
					if castlingRights&BK == BK &&
						(bb.AllPieces&(SquareF8.ToBB()|SquareG8.ToBB())) == 0 && // đường trống
						!IsAttacked(SquareE8, Black, bb) &&
						!IsAttacked(SquareF8, Black, bb) &&
						!IsAttacked(SquareG8, Black, bb) {
						moves |= SquareG8.ToBB()
					}

					// long castle (O-O-O)
					if castlingRights&BQ == BQ &&
						(bb.AllPieces&(SquareB8.ToBB()|SquareC8.ToBB()|SquareD8.ToBB())) == 0 &&
						!IsAttacked(SquareE8, Black, bb) &&
						!IsAttacked(SquareD8, Black, bb) &&
						!IsAttacked(SquareC8, Black, bb) {
						moves |= SquareD8.ToBB()
					}
				}
			}

			for moves != 0 {
				to := popLSB(&moves)

				isCapture := (enemyPieces>>to)&1 == 1
				isDoublePush := false
				isEnPassant := false
				isCastle := false

				switch pieceType {
				case Pawn:
					forward := 8
					if side == Black {
						forward = -8
					}
					// double push
					if to == from+Square(forward*2) {
						isDoublePush = true
					}
					// en passant
					if to == enPassant {
						isEnPassant = true
						isCapture = true
					}

				case King:
					// castling
					if abs(int(to)-int(from)) == 2 {
						isCastle = true
					}
				}

				move := NewMove(
					uint32(from),
					uint32(to),
					uint32(encodedPieceIndex(getPieceBySide(pieceType, side))),
					0, // promo
					boolToUint32(isCapture),
					boolToUint32(isDoublePush),
					boolToUint32(isEnPassant),
					boolToUint32(isCastle),
				)

				bbCopy := bb.Copy()
				makeUnsafeMove(bbCopy, move)

				if !IsKingAttacked(side, bbCopy) {
					return true
				}
			}
		}
	}

	return false
}
