// Game state representation for a chess engine

package chess_core

import (
	"fmt"
	"strings"
	"sync"
)

// GameState represents the current state of the chess game
type GameState struct {
	BitBoards  *BitBoards
	SideToMove Color

	CastlingRights  int    // bitmask for castling rights (1=White kingside, 2=White queenside, 4=Black kingside, 8=Black queenside)
	EnPassantSquare Square // square index for en passant target square, 0-63, or 64 if none
	HalfmoveClock   int    // for fifty-move rule
	FullmoveNumber  int    // starts at 1, increments after Black

	startFen string
	history  History

	state GameResult

	mx sync.Mutex
}

func NewGame() *GameState {
	return &GameState{
		BitBoards: &BitBoards{},
		history:   make(History, 50),
	}
}

func (gs *GameState) String() string {
	return gs.BitBoards.String()
}

func (gs *GameState) Copy() *GameState {
	newGS := &GameState{
		BitBoards:       gs.BitBoards.Copy(),
		SideToMove:      gs.SideToMove,
		CastlingRights:  gs.CastlingRights,
		EnPassantSquare: gs.EnPassantSquare,
		HalfmoveClock:   gs.HalfmoveClock,
		FullmoveNumber:  gs.FullmoveNumber,

		state: gs.state,
	}
	return newGS
}

func (gs *GameState) History() History {
	gs.mx.Lock()
	defer gs.mx.Unlock()

	return gs.history.Copy()
}

func (gs *GameState) FromFEN(fen string) error {
	gs.mx.Lock()
	defer gs.mx.Unlock()

	if gs.BitBoards == nil {
		gs.BitBoards = &BitBoards{}
	}

	if gs.history == nil {
		gs.history = make(History, 0, 50)
	}
	gs.history.Clear()

	if err := gs.BitBoards.FromFEN(fen); err != nil {
		return wrapError(err, "bitboards FromFEN")
	}

	// Parse FEN string to set other fields
	var sideToMove string
	var castling string
	var enPassant string
	var halfmove int
	var fullmove int

	_, lastFen, _ := strings.Cut(fen, " ")

	_, err := fmt.Sscanf(lastFen, "%s %s %s %d %d", &sideToMove, &castling, &enPassant, &halfmove, &fullmove)
	if err != nil {
		return wrapError(ErrInvalidFEN, "arguments parsing in FromFEN")
	}

	switch sideToMove {
	case "w":
		gs.SideToMove = White
	case "b":
		gs.SideToMove = Black
	default:
		return wrapError(ErrInvalidSideToMove, "side to move in FromFEN")
	}

	gs.CastlingRights = 0
	for _, c := range castling {
		switch c {
		case 'K':
			gs.CastlingRights |= 1
		case 'Q':
			gs.CastlingRights |= 2
		case 'k':
			gs.CastlingRights |= 4
		case 'q':
			gs.CastlingRights |= 8
		case '-':
			// no castling rights
		default:
			return wrapError(ErrInvalidCastling, "castling rights in FromFEN")
		}
	}

	if enPassant == "-" {
		gs.EnPassantSquare = Square(NoEnPassant)
	} else {
		file := enPassant[0] - 'a'
		rank := enPassant[1] - '1'
		if file > 7 || rank > 7 {
			return wrapError(ErrInvalidEnPassant, "en passant square in FromFEN")
		}
		gs.EnPassantSquare = Square(rank*8 + file)
	}

	if halfmove < 0 {
		return wrapError(ErrNegativeHalfmove, "halfmove in FromFEN")
	}
	gs.HalfmoveClock = halfmove

	if fullmove <= 0 {
		return wrapError(ErrNegativeFullmove, "fullmove in FromFEN")
	}
	gs.FullmoveNumber = fullmove

	// Locate kings
	whiteKingBB := gs.BitBoards.WhiteKing.LeastSignificantBit()
	blackKingBB := gs.BitBoards.BlackKing.LeastSignificantBit()

	if whiteKingBB == -1 {
		return wrapError(ErrNoKing, "white king in FromFEN")
	}
	if blackKingBB == -1 {
		return wrapError(ErrNoKing, "black king in FromFEN")
	}
	if gs.BitBoards.WhiteKing.Count() > 1 {
		return wrapError(ErrMultipleKings, "white king in FromFEN")
	}
	if gs.BitBoards.BlackKing.Count() > 1 {
		return wrapError(ErrMultipleKings, "black king in FromFEN")
	}

	gs.startFen = fen

	return nil
}

func (gs *GameState) ToFEN() string {
	gs.mx.Lock()
	defer gs.mx.Unlock()

	bb := gs.BitBoards

	var fenBoard strings.Builder
	for rank := 7; rank >= 0; rank-- {
		emptyCount := 0
		for file := range 8 {
			sq := rank*8 + file
			piece := bb.GetPieceAt(Square(sq))

			if piece == Empty {
				emptyCount++
			} else {
				if emptyCount > 0 {
					fenBoard.WriteString(fmt.Sprintf("%d", emptyCount))
					emptyCount = 0
				}
				fenBoard.WriteByte(byte(piece))
			}
		}
		if emptyCount > 0 {
			fenBoard.WriteString(fmt.Sprintf("%d", emptyCount))
		}
		if rank > 0 {
			fenBoard.WriteByte('/')
		}
	}

	side := "w"
	if gs.SideToMove == Black {
		side = "b"
	}

	// Castling rights
	castle := ""
	if gs.CastlingRights&1 != 0 {
		castle += "K"
	}
	if gs.CastlingRights&2 != 0 {
		castle += "Q"
	}
	if gs.CastlingRights&4 != 0 {
		castle += "k"
	}
	if gs.CastlingRights&8 != 0 {
		castle += "q"
	}
	if castle == "" {
		castle = "-"
	}

	//  En passant
	enpassant := "-"
	if gs.EnPassantSquare != NoEnPassant {
		enpassant = SquareToCoordinates[gs.EnPassantSquare]
	}

	// Combine everything
	return fmt.Sprintf("%s %s %s %s %d %d",
		fenBoard.String(),
		side,
		castle,
		enpassant,
		gs.HalfmoveClock,
		gs.FullmoveNumber,
	)
}

func (gs *GameState) createMove(from, to Square, promo PieceType) (Move, error) {
	pieceToMove := gs.BitBoards.GetPieceAt(from)
	if pieceToMove == Empty {
		return 0, ErrInvalidMove
	}

	side := pieceToMove.Color()
	allyPieces := gs.BitBoards.OccupiedBy(side)
	enemyPieces := gs.BitBoards.OccupiedBy(side.Opposite())
	allPieces := gs.BitBoards.AllPieces

	switch pieceToMove.Type() {
	case Pawn:
		var promotionRank, forwardStep, doublePushRank int
		var isSingle, isDouble, isCapture bool

		if side == White {
			forwardStep = 8
			promotionRank = 8
			doublePushRank = 2
		} else {
			forwardStep = -8
			promotionRank = 1
			doublePushRank = 7
		}

		isSingle = from+Square(forwardStep) == to && to.ToBB()&allPieces == 0
		if isSingle {
			isPromotion := to.Rank() == promotionRank
			if isPromotion {
				switch promo {
				case Queen:
					return NewMove(uint32(from), uint32(to), uint32(encodedPieceIndex(pieceToMove)), uint32(Queen), 0, 0, 0, 0), nil
				case Rook:
					return NewMove(uint32(from), uint32(to), uint32(encodedPieceIndex(pieceToMove)), uint32(Rook), 0, 0, 0, 0), nil
				case Bishop:
					return NewMove(uint32(from), uint32(to), uint32(encodedPieceIndex(pieceToMove)), uint32(Bishop), 0, 0, 0, 0), nil
				case Knight:
					return NewMove(uint32(from), uint32(to), uint32(encodedPieceIndex(pieceToMove)), uint32(Knight), 0, 0, 0, 0), nil
				default:
					return 0, ErrInvalidPromotion
				}
			}
			return NewMove(uint32(from), uint32(to), uint32(encodedPieceIndex(pieceToMove)), 0, 0, 0, 0, 0), nil
		}

		midSquare := from + Square(forwardStep)
		isDouble = from+Square(forwardStep*2) == to &&
			to.ToBB()&allPieces == 0 &&
			midSquare.ToBB()&allPieces == 0 &&
			from.Rank() == doublePushRank
		if isDouble {
			return NewMove(uint32(from), uint32(to), uint32(encodedPieceIndex(pieceToMove)), 0, 0, 1, 0, 0), nil
		}

		isCapture = PawnAttacks(from, side)&to.ToBB() != 0 &&
			((to.ToBB()&enemyPieces) != 0 || to == gs.EnPassantSquare)
		if isCapture {
			isPromotion := to.Rank() == promotionRank
			if isPromotion {
				switch promo {
				case Queen:
					return NewMove(uint32(from), uint32(to), uint32(encodedPieceIndex(pieceToMove)), uint32(Queen), 1, 0, 0, 0), nil
				case Rook:
					return NewMove(uint32(from), uint32(to), uint32(encodedPieceIndex(pieceToMove)), uint32(Rook), 1, 0, 0, 0), nil
				case Bishop:
					return NewMove(uint32(from), uint32(to), uint32(encodedPieceIndex(pieceToMove)), uint32(Bishop), 1, 0, 0, 0), nil
				case Knight:
					return NewMove(uint32(from), uint32(to), uint32(encodedPieceIndex(pieceToMove)), uint32(Knight), 1, 0, 0, 0), nil
				default:
					return 0, ErrInvalidPromotion
				}
			}

			isEnPassant := to == gs.EnPassantSquare
			if isEnPassant {
				return NewMove(uint32(from), uint32(to), uint32(encodedPieceIndex(pieceToMove)), 0, 1, 0, 1, 0), nil
			}

			return NewMove(uint32(from), uint32(to), uint32(encodedPieceIndex(pieceToMove)), 0, 1, 0, 0, 0), nil
		}

		return 0, ErrInvalidMove
	case Rook:
		attacks := RookAttacks(from, allPieces) & ^allyPieces
		if to.ToBB()&attacks == 0 {
			return 0, ErrInvalidMove
		}
		isCapture := (enemyPieces >> to) & 1
		return NewMove(uint32(from), uint32(to), uint32(encodedPieceIndex(pieceToMove)), 0, uint32(isCapture), 0, 0, 0), nil
	case Bishop:
		attacks := BishopAttacks(from, allPieces) & ^allyPieces
		if to.ToBB()&attacks == 0 {
			return 0, ErrInvalidMove
		}
		isCapture := (enemyPieces >> to) & 1
		return NewMove(uint32(from), uint32(to), uint32(encodedPieceIndex(pieceToMove)), 0, uint32(isCapture), 0, 0, 0), nil
	case Knight:
		attacks := KnightAttacks(from) & ^allyPieces
		if to.ToBB()&attacks == 0 {
			return 0, ErrInvalidMove
		}
		isCapture := (enemyPieces >> to) & 1
		return NewMove(uint32(from), uint32(to), uint32(encodedPieceIndex(pieceToMove)), 0, uint32(isCapture), 0, 0, 0), nil
	case Queen:
		attacks := QueenAttacks(from, allPieces) & ^allyPieces
		if to.ToBB()&attacks == 0 {
			return 0, ErrInvalidMove
		}
		isCapture := (enemyPieces >> to) & 1
		return NewMove(uint32(from), uint32(to), uint32(encodedPieceIndex(pieceToMove)), 0, uint32(isCapture), 0, 0, 0), nil
	case King:
		attacks := KingAttacks(from) & ^allyPieces
		if to.ToBB()&attacks != 0 && !IsAttacked(to, side, gs.BitBoards) {
			isCapture := (enemyPieces >> to) & 1
			return NewMove(uint32(from), uint32(to), uint32(encodedPieceIndex(pieceToMove)), 0, uint32(isCapture), 0, 0, 0), nil
		}

		// castling

		if side == White && from == SquareE1 {
			// Kingside castling (O-O): e1 -> e1, rook h1 -> f1
			if to == SquareG1 && (gs.CastlingRights&WK) != 0 {
				if allPieces&(SquareF1.ToBB()|SquareG1.ToBB()) == 0 &&
					!IsAttacked(SquareE1, White, gs.BitBoards) &&
					!IsAttacked(SquareF1, White, gs.BitBoards) &&
					!IsAttacked(SquareG1, White, gs.BitBoards) {
					return NewMove(uint32(from), uint32(to), uint32(encodedPieceIndex(pieceToMove)), 0, 0, 0, 0, 1), nil
				}
			}

			if to == SquareC1 && (gs.CastlingRights&WQ) != 0 {
				if allPieces&(SquareB1.ToBB()|SquareC1.ToBB()|SquareD1.ToBB()) == 0 &&
					!IsAttacked(SquareE1, White, gs.BitBoards) &&
					!IsAttacked(SquareD1, White, gs.BitBoards) &&
					!IsAttacked(SquareC1, White, gs.BitBoards) {
					return NewMove(uint32(from), uint32(to), uint32(encodedPieceIndex(pieceToMove)), 0, 0, 0, 0, 1), nil
				}
			}
		}

		if side == Black && from == SquareE8 {
			// Kingside castling (O-O): e8 -> g8, rook h8 -> f8
			if to == SquareG8 && (gs.CastlingRights&BK) != 0 {
				if allPieces&(SquareF8.ToBB()|SquareG8.ToBB()) == 0 &&
					!IsAttacked(SquareE8, Black, gs.BitBoards) &&
					!IsAttacked(SquareF8, Black, gs.BitBoards) &&
					!IsAttacked(SquareG8, Black, gs.BitBoards) {
					return NewMove(uint32(from), uint32(to), uint32(encodedPieceIndex(pieceToMove)), 0, 0, 0, 0, 1), nil
				}
			}
			// Queenside castling (O-O-O): e8 -> c8, rook a8 -> d8
			if to == SquareC8 && (gs.CastlingRights&BQ) != 0 {
				if allPieces&(SquareB8.ToBB()|SquareC8.ToBB()|SquareD8.ToBB()) == 0 &&
					!IsAttacked(SquareE8, Black, gs.BitBoards) &&
					!IsAttacked(SquareD8, Black, gs.BitBoards) &&
					!IsAttacked(SquareC8, Black, gs.BitBoards) {
					return NewMove(uint32(from), uint32(to), uint32(encodedPieceIndex(pieceToMove)), 0, 0, 0, 0, 1), nil
				}
			}
		}

		return 0, ErrInvalidMove
	default:
		return 0, ErrInvalidMove
	}
}

func (gs *GameState) MakeMove(side Color, from, to Square, promo PieceType) (GameResult, error) {
	gs.mx.Lock()
	defer gs.mx.Unlock()

	if gs.state != ResultOngoing {
		return gs.state, ErrMatchEnd
	}

	// handle validate move

	if gs.SideToMove != side {
		return "", ErrMoveOutOfTurn
	}

	move, err := gs.createMove(from, to, promo)
	if err != nil {
		return "", err
	}

	bbs := gs.BitBoards.Copy()

	makeUnsafeMove(bbs, move)

	if IsKingAttacked(side, bbs) {
		return "", ErrMoveIntoCheck
	}

	// Set new state

	// push history
	currentHash := gs.computeZobristHash()
	gs.history.Push(move, currentHash)

	// set BitBoards
	gs.BitBoards = bbs

	if side == Black {
		gs.FullmoveNumber += 1
	}
	gs.SideToMove = side.Opposite()

	// CastlingRights update
	gs.CastlingRights &= CastlingRights[from]
	gs.CastlingRights &= CastlingRights[to]

	// HalfmoveClock update
	if move.IsCapture() || ASCIIPieces[move.Piece()].Type() == Pawn {
		gs.HalfmoveClock = 0
	} else {
		gs.HalfmoveClock++
	}

	// EnPassantSquare update
	if move.IsDoublePush() {
		if side == White {
			gs.EnPassantSquare = Square(move.To()) - 8
		} else {
			gs.EnPassantSquare = Square(move.To()) + 8
		}
	} else {
		gs.EnPassantSquare = NoEnPassant
	}

	if IsKingAttacked(side.Opposite(), bbs) {
		if !hasAnyLegalMove(bbs, side.Opposite(), gs.CastlingRights, gs.EnPassantSquare) {
			gs.state = ResultCheckmate
		}
	} else {
		if !hasAnyLegalMove(bbs, side.Opposite(), gs.CastlingRights, gs.EnPassantSquare) {
			gs.state = ResultStalemate
		}
	}

	if gs.isInsufficientMaterial() {
		gs.state = ResultInsufficientMaterial
		return gs.state, nil
	}

	if gs.isThreefoldRepetition() {
		gs.state = ResultThreefoldRepetition
		return gs.state, nil
	}

	if gs.HalfmoveClock >= 150 { // 75 moves per side = 150 ply
		gs.state = ResultDrawBy75Move
		return gs.state, nil
	}

	return gs.state, nil
}

func (gs *GameState) MakeDrawBy50Move() error {
	gs.mx.Lock()
	defer gs.mx.Unlock()

	// chỉ được yêu cầu hòa khi đã đủ 50 nước (100 ply)
	if gs.HalfmoveClock < 100 {
		return fmt.Errorf("cannot claim draw yet: only %d halfmoves", gs.HalfmoveClock)
	}

	// chỉ có thể yêu cầu hòa khi ván chưa kết thúc
	switch gs.state {
	case ResultCheckmate, ResultStalemate, ResultDrawBy75Move, ResultInsufficientMaterial:
		return ErrMatchEnd
	}

	// đặt trạng thái hòa
	gs.state = ResultDrawBy75Move
	return nil
}

func (gs *GameState) CanMakeDrawBy50Move() bool {
	gs.mx.Lock()
	defer gs.mx.Unlock()

	return gs.HalfmoveClock > 100
}

func (gs *GameState) isInsufficientMaterial() bool {
	bb := gs.BitBoards

	// Đếm quân mỗi bên
	whitePieces := bb.WhitePieces
	blackPieces := bb.BlackPieces

	whiteCount := whitePieces.Count()
	blackCount := blackPieces.Count()

	// Nếu chỉ còn vua mỗi bên
	if whiteCount == 1 && blackCount == 1 {
		return true
	}

	// Nếu chỉ còn vua + 1 quân nhẹ (mã hoặc tượng)
	whiteKnights := bb.WhiteKnights.Count()
	whiteBishops := bb.WhiteBishops.Count()
	blackKnights := bb.BlackKnights.Count()
	blackBishops := bb.BlackBishops.Count()

	// Vua + 1 quân nhẹ vs Vua
	if (whiteCount == 2 && (whiteKnights == 1 || whiteBishops == 1) && blackCount == 1) ||
		(blackCount == 2 && (blackKnights == 1 || blackBishops == 1) && whiteCount == 1) {
		return true
	}

	// Vua + tượng vs Vua + tượng cùng màu
	if whiteCount == 2 && blackCount == 2 &&
		whiteBishops == 1 && blackBishops == 1 &&
		sameBishopColor(bb.WhiteBishops, bb.BlackBishops) {
		return true
	}

	return false
}

func (gs *GameState) isThreefoldRepetition() bool {
	currentHash := computeZobristHash(gs.BitBoards, gs.SideToMove, gs.EnPassantSquare, gs.CastlingRights)
	count := 0

	// chỉ cần xét các thế trong phạm vi halfmoveClock nước gần nhất
	limit := max(len(gs.history)-gs.HalfmoveClock, 0)

	for i := len(gs.history) - 1; i >= limit; i-- {
		if gs.history[i].Hash == currentHash {
			count++
			if count >= 2 {
				return true
			}
		}
	}
	return false
}

func (gs *GameState) canForceCheckmate(playerColor Color) bool {
	bb := gs.BitBoards

	var checkingPieces BitBoard
	var targetPieces BitBoard
	var checkingKnights int
	var checkingBishops int
	var targetPiecesCount int // Luôn là 1, nhưng để nhất quán

	if playerColor == White {
		checkingPieces = bb.WhitePieces
		targetPieces = bb.BlackPieces
		checkingKnights = bb.WhiteKnights.Count()
		checkingBishops = bb.WhiteBishops.Count()
	} else {
		checkingPieces = bb.BlackPieces
		targetPieces = bb.WhitePieces
		checkingKnights = bb.BlackKnights.Count()
		checkingBishops = bb.BlackBishops.Count()
	}

	targetPiecesCount = targetPieces.Count()
	checkingPiecesCount := checkingPieces.Count()

	// Trường hợp 1: Có quân nặng/Tốt (Đảm bảo thắng)
	if checkingPiecesCount > (checkingKnights + checkingBishops + 1) { // 1 là Vua
		// Nếu có hơn Vua, Mã, Tượng (ví dụ: Hậu, Xe, Tốt)
		return true
	}

	// Trường hợp 2: K vs K
	if checkingPiecesCount == 1 {
		return false // Chỉ còn Vua
	}

	if targetPiecesCount > 1 {
		return true
	}

	// Trường hợp 3: K+N vs K, hoặc K+B vs K
	if checkingPiecesCount == 2 {
		if checkingKnights == 1 || checkingBishops == 1 {
			return false // K+N vs K, hoặc K+B vs K: KHÔNG THỂ chiếu hết
		}
	}

	// Trường hợp 4: K+N+N vs K
	if checkingPiecesCount == 3 && checkingKnights == 2 {
		return false // K+N+N vs K: KHÔNG THỂ chiếu hết theo luật FIDE cơ bản cho timeout
	}

	// Trường hợp 5: K+B+B vs K (Cần kiểm tra màu ô)
	if checkingPiecesCount == 3 && checkingBishops == 2 {
		// Lấy vị trí 2 Tượng
		var bishop1, bishop2 Square
		if playerColor == Black {
			bishop1 = Square(bb.BlackBishops.LeastSignificantBit())
			bishop2 = Square(bb.BlackBishops.MostSignificantBit())
		} else {
			bishop1 = Square(bb.WhiteBishops.LeastSignificantBit())
			bishop2 = Square(bb.WhiteBishops.MostSignificantBit())
		}

		// Kiểm tra màu ô
		isLight1 := (bishop1.File()+bishop1.Rank())%2 == 0
		isLight2 := (bishop2.File()+bishop2.Rank())%2 == 0

		return isLight1 != isLight2
	}

	return true
}

func (gs *GameState) CanForceCheckmate(playerColor Color) bool {
	gs.mx.Lock()
	defer gs.mx.Unlock()

	return gs.canForceCheckmate(playerColor)
}

// get winner player
func (gs *GameState) Winner() (Color, bool) {
	if gs.state == ResultCheckmate {
		return gs.SideToMove.Opposite(), true
	}
	return None, false
}

// get current game state
func (gs *GameState) State() GameResult {
	return gs.state
}

func (gs *GameState) computeZobristHash() uint64 {
	return computeZobristHash(gs.BitBoards, gs.SideToMove, gs.EnPassantSquare, gs.CastlingRights)
}

// kiểm tra xem 2 tượng có cùng màu ô không
func sameBishopColor(whiteBishop, blackBishop BitBoard) bool {
	if whiteBishop == 0 || blackBishop == 0 {
		return false // Một bên không có tượng thì không cần so sánh
	}

	whiteSq := Square(whiteBishop.LeastSignificantBit())
	blackSq := Square(blackBishop.LeastSignificantBit())

	isLight1 := (whiteSq.File()+whiteSq.Rank())%2 == 0
	isLight2 := (blackSq.File()+blackSq.Rank())%2 == 0

	return isLight1 == isLight2
}

type MoveHistory struct {
	Move Move
	Hash uint64
}

type History []MoveHistory

func (h *History) Clear() {
	*h = (*h)[:0]
}

func (h *History) Push(move Move, hash uint64) {
	*h = append(*h, MoveHistory{Move: move, Hash: hash})
}

func (h *History) Pop() MoveHistory {
	if len(*h) == 0 {
		return MoveHistory{}
	}
	last := (*h)[len(*h)-1]
	*h = (*h)[:len(*h)-1]
	return last
}

func (h *History) CountHash(hash uint64) int {
	count := 0
	for _, entry := range *h {
		if entry.Hash == hash {
			count++
		}
	}
	return count
}

func (ml *History) Copy() History {
	moves := make(History, len(*ml))
	copy(moves, *ml)
	return moves
}
