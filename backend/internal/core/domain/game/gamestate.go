package game

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	chess "github.com/tommjj/chess_OG/chess_core"
)

// EventType represents the type of game event
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

// GameState represents the state of a chess game.
// It encapsulates the chess game logic, timer, and game status.
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

// NewGame creates a new GameState instance. you need call Start() to start the game timer.
//
// fen: initial fen string
// timeSeconds: time per side in seconds
// endCallBack: callback function when game ends
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

// Start the game timer. Returns true if the timer was started, false if it was already started.
// you should call this method when the game actually starts (e.g., after the first move).
func (g *GameState) Start() bool {
	g.mu.Lock()
	defer g.mu.Unlock()

	if !g.timer.HasStarted() {
		g.timer.Start()
		return true
	}
	return false
}

// Pause the game timer. Returns true if the timer was paused, false if it was already paused.
func (g *GameState) Pause() bool {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.timer.IsRunning() {
		g.timer.Stop()
		return true
	}
	return false
}

// Resume the game timer. Returns true if the timer was resumed, false if it was already running.
func (g *GameState) Resume() bool {
	g.mu.Lock()
	defer g.mu.Unlock()

	if !g.timer.IsRunning() && g.timer.HasStarted() {
		g.timer.Start()
		return true
	}
	return false
}

// IsRunning returns true if the game timer is running.
func (g *GameState) IsRunning() bool {
	g.mu.Lock()
	defer g.mu.Unlock()

	return g.timer.IsRunning()
}

// HasStarted returns true if the game timer has started.
func (g *GameState) HasStarted() bool {
	g.mu.Lock()
	defer g.mu.Unlock()

	return g.timer.HasStarted()
}

// EndByLeaveGame this method ends the game when a player leaves.
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

// EndByForfeit this method ends the game when a player forfeits.
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

// MakeMove makes a move on the game state.
//
//	side: the side making the move
//	from: the square the piece is moving from
//	to: the square the piece is moving to
//	promo: the piece type to promote to (if applicable) 0 if no promotion
//
//	returns the game status after the move and an error if the move is invalid or if the game has ended.
func (g *GameState) MakeMove(side Color, from Square, to Square, promo PieceType) (GameStatus, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.status != ResultOngoing { // match end
		return g.status, ErrMatchEnd
	}

	if !g.timer.IsRunning() { // game is not running
		if g.timer.HasStarted() {
			return "", ErrGamePaused
		} else {
			return "", ErrGameNotStarted
		}
	}

	remaining := g.timer.Remaining(side)
	const safeThreshold = 30 * time.Millisecond
	if remaining <= safeThreshold {
		return "", ErrTimeout
	}

	result, err := g.state.MakeMove(side, from, to, promo)
	if err != nil {
		return result, err
	}
	g.status = result

	// handle match end
	if result != ResultOngoing {
		g.timer.Stop()
		if result == ResultCheckmate {
			g.winner = side
		} else {
			g.winner = Both
		}
		go g.handleMatchEnd()
	} else {
		ok := g.timer.SwitchTurn()
		if !ok { // rollback move
			g.state.Undo(1)           // undo last move
			return result, ErrTimeout // cant switch turn, timeout
		}
	}
	g.currentFen = g.state.ToFEN()

	return result, err
}

func (g *GameState) handleMatchEnd() {
	g.mu.Lock()

	if g.endCallBack != nil {
		return
	} // return if callback is nil

	result := GameResult{
		Winner:    g.winner,
		Result:    g.status,
		Duration:  g.timer.GetDuration(),
		BlackTime: g.timer.BlackRemaining(),
		WhiteTime: g.timer.BlackRemaining(),
		StartFen:  g.state.StartFen(),
		FinalFen:  g.state.ToFEN(),
	}

	history := g.state.History()
	moves := make([]Move, len(history))

	for i, v := range history {
		moves[i] = v.Move
	}
	result.Moves = moves

	g.mu.Unlock()

	g.endCallBack(result)
}

// timeoutColor
//
//	Black - Black times out
//	White - White times out
//	Both - Draw by timeout
func (g *GameState) handleTimeout(timeoutColor Color) {
	g.mu.Lock()

	if g.status != ResultOngoing { // Do nothing if game ended
		return
	}

	// handle result
	oppColor := timeoutColor.Opposite()
	canForceCheckmate := g.state.CanForceCheckmate(oppColor)

	if canForceCheckmate {
		g.winner = oppColor
		g.status = ResultTimeout
	} else {
		g.winner = Both
		g.status = ResultTimeout
	}

	g.mu.Unlock()
	g.handleMatchEnd()
}

// MakeDraw makes the game a draw either by agreement or by 50-move rule.
func (g *GameState) MakeDraw() error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.status != ResultOngoing { // Do nothing if game ended
		return ErrMatchEnd
	}

	g.timer.Stop()

	drawResult := ResultDrawByAgreement
	if g.state.CanDrawBy50Move() {
		drawResult = ResultDrawBy50Move
	}

	g.status = drawResult
	g.winner = Both

	go g.handleMatchEnd()
	return nil
}
