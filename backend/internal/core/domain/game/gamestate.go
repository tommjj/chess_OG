package game

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	chess "github.com/tommjj/chess_OG/chess_core"
)

type GameEvent struct {
	EventType EventType // Event type

	Move         Move
	SequenceTick int
	MoveColor    Color

	BlackTime time.Duration
	WhiteTime time.Duration

	NewStatus GameStatus
	Winner    Color

	Timestamp time.Time
}

// Game Result
type GameResult struct {
	// Current winner
	//  Black - Black wins
	//  White - White wins
	//  Both - Draw
	//  None - Game is ongoing
	Winner Color
	Result GameStatus

	// time
	Duration  time.Duration
	BlackTime time.Duration
	WhiteTime time.Duration

	Moves    []Move
	StartFen string
	FinalFen string
}

type gameState = chess.GameState

type GameState struct {
	state      *gameState
	currentFen string

	timer *timer

	status GameStatus
	// Current winner
	//  Black - Black wins
	//  White - White wins
	//  Both - Draw
	//  None - Game is ongoing
	winner Color

	endCallBack func(result GameResult)

	mu sync.Mutex
}

func NewGame(ctx context.Context, fen string, timeSeconds int, endCallBack func(result GameResult)) (*GameState, error) {
	board := chess.NewGame()
	err := board.FromFEN(fen)
	if err != nil {
		return nil, err
	}

	s := &GameState{
		currentFen:  fen,
		state:       board,
		status:      chess.ResultOngoing,
		winner:      None,
		endCallBack: endCallBack,
	}

	s.timer = NewTimer(timeSeconds, time.Second, board.SideToMove, s.handleTimeout)

	return s, nil
}

func (g *GameState) Start() bool {
	g.mu.Lock()
	defer g.mu.Unlock()

	if !g.timer.HasStarted() {
		g.timer.Start()
		return true
	}
	return false
}

func (g *GameState) Pause() bool {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.timer.IsRunning() {
		g.timer.Stop()
		return true
	}
	return false
}

func (g *GameState) Resume() bool {
	g.mu.Lock()
	defer g.mu.Unlock()

	if !g.timer.IsRunning() && g.timer.HasStarted() {
		g.timer.Start()
		return true
	}
	return false
}

func (g *GameState) IsRunning() bool {
	g.mu.Lock()
	defer g.mu.Unlock()

	return g.timer.IsRunning()
}

func (g *GameState) HasStarted() bool {
	g.mu.Lock()
	defer g.mu.Unlock()

	return g.timer.HasStarted()
}

// color
//
//	Black - Black wins
//	White - White wins
func (g *GameState) EndByLeaveGame(color Color) error {
	if color != Black && color != White {
		return errors.New("EndByLeaveGame: invalid color")
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.status != ResultOngoing {
		return fmt.Errorf("game already ended with result: %s", g.state)
	}

	g.timer.Stop()
	if g.state.CanForceCheckmate(color.Opposite()) {
		g.winner = color.Opposite()
	} else {
		g.winner = Both
	}
	g.status = ResultResignation

	return nil
}

// color
//
//	Black - Black wins
//	White - White wins
//	Both - Draw
func (g *GameState) EndByForfeit(color Color) error {
	if color == None {
		return errors.New("EndByForfeit: invalid color")
	}

	g.mu.Lock()
	defer g.mu.Unlock()
	if g.status != ResultOngoing {
		return fmt.Errorf("game already ended with result: %s", g.state)
	}

	g.timer.Stop()
	g.winner = color
	g.status = ResultForfeit
	return nil
}

func (g *GameState) MakeMove(side Color, from Square, to Square, promo PieceType) (GameStatus, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.status != ResultOngoing {
		return g.status, ErrMatchEnd
	}

	if !g.timer.IsRunning() { // game is not run
		if g.timer.HasStarted() {
			return "", ErrGamePaused
		} else {
			return "", ErrGameNotStarted
		}
	}

	result, err := g.state.MakeMove(side, from, to, promo)
	if err != nil {
		return result, err
	}
	g.status = result

	// handle match end
	if result != ResultOngoing {
		g.timer.Stop()
		// call back
	} else {
		g.timer.SwitchTurn()

	}
	g.currentFen = g.state.ToFEN()

	return result, err
}

func (g *GameState) handleMatchEnd() {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.endCallBack != nil {
		return
	} // return if callback is nil

	// result := GameResult{
	// 	Winner:    g.winner,
	// 	Result:    g.status,
	// 	Duration:  g.timer.GetDuration(),
	// 	BlackTime: g.timer.getBlackTime(),
	// 	WhiteTime: g.timer.getWhiteTime(),
	// }
}

func (g *GameState) handleTimeout(timeoutColor Color) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.status != ResultOngoing { // Do nothing if game ended
		return
	}

	// handle result

	go g.handleMatchEnd()
}

func (g *GameState) MakeDraw() error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.status != ResultOngoing { // Do nothing if game ended
		return ErrMatchEnd
	}

	return nil
}
