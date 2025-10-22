package game

import (
	chess "github.com/tommjj/chess_OG/chess_core"
)

type gameState = chess.GameState

type ChessGame struct {
	state gameState

	timer *Timer

	status GameResult

	DoneCh <-chan struct{}
}
