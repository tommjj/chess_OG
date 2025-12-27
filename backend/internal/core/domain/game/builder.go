// builder game state builder

package game

import "time"

type GameMode string

// Game mode constants
// Bt - Bullet
// Bz - Blitz
// Rd - Rapid
// Cl - Classical
const (
	ModeBt1m0s GameMode = "bt_1m_0s"
	ModeBt2m1s GameMode = "bt_2m_1s"

	ModeBz3m0s GameMode = "bz_3m_0s"
	ModeBz3m2s GameMode = "bz_3m_2s"
	ModeBz5m0s GameMode = "bz_5m_0s"
	ModeBz5m5s GameMode = "bz_5m_5s"

	ModeRd10m0s  GameMode = "rd_10m_0s"
	ModeRd15m10s GameMode = "rd_15m_10s"

	ModeCl30m0s GameMode = "cl_30m_0s"
	ModeCl60m0s GameMode = "cl_60m_0s"
)

type timeControl struct {
	minutes          int
	incrementSeconds int
}

var modeTimeControlMap = map[GameMode]timeControl{
	ModeBt1m0s: {1, 0},
	ModeBt2m1s: {2, 1},

	ModeBz3m0s: {3, 0},
	ModeBz3m2s: {3, 2},
	ModeBz5m0s: {5, 0},
	ModeBz5m5s: {5, 5},

	ModeRd10m0s:  {10, 0},
	ModeRd15m10s: {15, 10},

	ModeCl30m0s: {30, 0},
	ModeCl60m0s: {60, 0},
}

func InvalidGameMode(mode GameMode) bool {
	_, ok := modeTimeControlMap[mode]
	return !ok
}

func BuildGameTimeControl(mode GameMode) (minutes int, incrementSeconds int) {
	if tc, ok := modeTimeControlMap[mode]; ok {
		return tc.minutes, tc.incrementSeconds
	}
	return
}

const initialFEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

// BuildGameState builds a new GameState based on the given GameMode
func BuildGameState(mode GameMode, endCallBack func(result GameResult)) (*GameState, error) {
	minutes, incrementSeconds := BuildGameTimeControl(mode)

	if minutes == 0 {
		return nil, ErrInvalidGameMode
	}

	return NewGame(initialFEN, minutes*60, time.Duration(incrementSeconds)*time.Second, endCallBack)
}
