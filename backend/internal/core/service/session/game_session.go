// Game session service package
// this package manages game sessions, including creation, retrieval, and updates.

package session

import "github.com/tommjj/chess_OG/backend/internal/core/domain/game"

type GameSession struct {
	state *game.GameState // Game state
}
