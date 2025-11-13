package game

import (
	chess "github.com/tommjj/chess_OG/chess_core"
)

type EventType int

const (
	MoveMade EventType = iota
	TimeUpdated
	GameEnded
	GameStarted
	GameStopped
)

type Square = chess.Square
type Piece = chess.Piece
type Color = chess.Color
type PieceType = chess.PieceType
type Move = chess.Move

type GameStatus = chess.GameStatus

const (
	ResultOngoing              = GameStatus(chess.ResultOngoing)
	ResultCheckmate            = GameStatus(chess.ResultCheckmate)
	ResultStalemate            = GameStatus(chess.ResultStalemate)
	ResultDrawBy75Move         = GameStatus(chess.ResultDrawBy75Move)
	ResultInsufficientMaterial = GameStatus(chess.ResultInsufficientMaterial)
	ResultThreefoldRepetition  = GameStatus(chess.ResultThreefoldRepetition)
	ResultTimeout              = GameStatus("Result Timeout")
	ResultDrawByTimeClaim      = GameStatus("Result Draw By Time Claim") // Yêu cầu hòa do hết giờ của đối thủ nhưng không đủ vật chất (theo luật FIDE)

	// player
	ResultResignation     = GameStatus("Result Resignation")       // Thua do đầu hàng (Người chơi tự nguyện Quit/Resign)
	ResultDrawByAgreement = GameStatus("Result Draw By Agreement") // Hòa do đồng thuận giữa hai người chơi
	ResultDrawBy50Move    = GameStatus("Result Draw By 50 Move")
	ResultForfeit         = GameStatus("Result Forfeit") // Thua do mất kết nối/hết thời gian kết nối lại (Walkover)

)

type EndReason string

const (
	EndReasonDrawAgreement  EndReason = "draw_agreement" // Hai bên đồng ý hòa
	EndReasonPlayerResigned EndReason = "resigned"       // Một bên đầu hàng
	EndReasonDisconnected   EndReason = "disconnected"   // Một bên thoát giữa chừng
	EndReasonAdminStop      EndReason = "admin_stop"     // Admin dừng trận
	EndReasonError          EndReason = "error"          // Lỗi kỹ thuật hoặc logic
)

var (
	White = Color(chess.White)
	Black = Color(chess.Black)

	Both = Color(chess.Both)
	None = Color(chess.None)
)

// encode pieces
const (
	WPawn   = Piece(chess.WPawn)
	WKnight = Piece(chess.WKnight)
	WBishop = Piece(chess.WBishop)
	WRook   = Piece(chess.WRook)
	WQueen  = Piece(chess.WQueen)
	WKing   = Piece(chess.WKing)

	BPawn   = Piece(chess.BPawn)
	BKnight = Piece(chess.BKnight)
	BBishop = Piece(chess.BBishop)
	BRook   = Piece(chess.BRook)
	BQueen  = Piece(chess.BQueen)
	BKing   = Piece(chess.BKing)

	Empty = Piece(chess.Empty)
)

const (
	Pawn   = PieceType(chess.Pawn)
	Knight = PieceType(chess.Knight)
	Bishop = PieceType(chess.Bishop)
	Rook   = PieceType(chess.Rook)
	Queen  = PieceType(chess.Queen)
	King   = PieceType(chess.King)
)

// Square
const (
	SquareA1 = Square(chess.SquareA1)
	SquareB1 = Square(chess.SquareB1)
	SquareC1 = Square(chess.SquareC1)
	SquareD1 = Square(chess.SquareD1)
	SquareE1 = Square(chess.SquareE1)
	SquareF1 = Square(chess.SquareF1)
	SquareG1 = Square(chess.SquareG1)
	SquareH1 = Square(chess.SquareH1)

	SquareA2 = Square(chess.SquareA2)
	SquareB2 = Square(chess.SquareB2)
	SquareC2 = Square(chess.SquareC2)
	SquareD2 = Square(chess.SquareD2)
	SquareE2 = Square(chess.SquareE2)
	SquareF2 = Square(chess.SquareF2)
	SquareG2 = Square(chess.SquareG2)
	SquareH2 = Square(chess.SquareH2)

	SquareA3 = Square(chess.SquareA3)
	SquareB3 = Square(chess.SquareB3)
	SquareC3 = Square(chess.SquareC3)
	SquareD3 = Square(chess.SquareD3)
	SquareE3 = Square(chess.SquareE3)
	SquareF3 = Square(chess.SquareF3)
	SquareG3 = Square(chess.SquareG3)
	SquareH3 = Square(chess.SquareH3)

	SquareA4 = Square(chess.SquareA4)
	SquareB4 = Square(chess.SquareB4)
	SquareC4 = Square(chess.SquareC4)
	SquareD4 = Square(chess.SquareD4)
	SquareE4 = Square(chess.SquareE4)
	SquareF4 = Square(chess.SquareF4)
	SquareG4 = Square(chess.SquareG4)
	SquareH4 = Square(chess.SquareH4)

	SquareA5 = Square(chess.SquareA5)
	SquareB5 = Square(chess.SquareB5)
	SquareC5 = Square(chess.SquareC5)
	SquareD5 = Square(chess.SquareD5)
	SquareE5 = Square(chess.SquareE5)
	SquareF5 = Square(chess.SquareF5)
	SquareG5 = Square(chess.SquareG5)
	SquareH5 = Square(chess.SquareH5)

	SquareA6 = Square(chess.SquareA6)
	SquareB6 = Square(chess.SquareB6)
	SquareC6 = Square(chess.SquareC6)
	SquareD6 = Square(chess.SquareD6)
	SquareE6 = Square(chess.SquareE6)
	SquareF6 = Square(chess.SquareF6)
	SquareG6 = Square(chess.SquareG6)
	SquareH6 = Square(chess.SquareH6)

	SquareA7 = Square(chess.SquareA7)
	SquareB7 = Square(chess.SquareB7)
	SquareC7 = Square(chess.SquareC7)
	SquareD7 = Square(chess.SquareD7)
	SquareE7 = Square(chess.SquareE7)
	SquareF7 = Square(chess.SquareF7)
	SquareG7 = Square(chess.SquareG7)
	SquareH7 = Square(chess.SquareH7)

	SquareA8 = Square(chess.SquareA8)
	SquareB8 = Square(chess.SquareB8)
	SquareC8 = Square(chess.SquareC8)
	SquareD8 = Square(chess.SquareD8)
	SquareE8 = Square(chess.SquareE8)
	SquareF8 = Square(chess.SquareF8)
	SquareG8 = Square(chess.SquareG8)
	SquareH8 = Square(chess.SquareH8)
)
