// Game session service package
// this package manages game sessions, including creation, retrieval, and updates.

package session

import (
	"time"

	"github.com/google/uuid"
	"github.com/tommjj/chess_OG/backend/internal/core/domain/game"
)

type Player struct {
	ID       string // Player ID
	Username string // Player username
	Avatar   string // Player avatar URL
}

type GameState struct {
	Fen    string          // FEN representation of the game state
	Status game.GameStatus // Current game status

	Black          Player        // Black player
	BlackRemaining time.Duration // Black player's remaining time

	White          Player        // White player
	WhiteRemaining time.Duration // White player's remaining time

	Moved *game.Move // Move that was just made
}

// GameSession struct to manage a chess game session
type GameSession struct {
	mode  game.GameMode   // Game mode (e.g., Blitz, Rapid)
	state *game.GameState // Game state

	white *Player // White player
	black *Player // Black player

	drawOfferedBy *Player // Player who offered a draw, nil if no offer

	spectators []uuid.UUID // List of spectator IDs

}

func NewGameSession(mode game.GameMode, white *Player, black *Player) *GameSession {
	return &GameSession{
		mode: mode,
		// state:      game

		white:         white,
		black:         black,
		drawOfferedBy: nil,
		spectators:    []uuid.UUID{},
	}
}

func (gs *GameSession) GetMode() game.GameMode {
	return gs.mode
}

func (gs *GameSession) GetState() *game.GameState {
	return gs.state
}

// func (gs *GameSession) AddSpectator(spectator Player) {
// 	gs.Spectators = append(gs.Spectators, spectator)
// }
