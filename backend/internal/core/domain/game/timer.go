package game

import (
	"sync"
	"time"
)

var NullTime = time.Time{}

// timer represents a chess game timer.
type timer struct {
	InitialTimeSeconds int // Initial time in seconds for each player.
	IncreaseDuration   time.Duration

	BlackTime time.Duration // Remaining time for Black.
	WhiteTime time.Duration // Remaining time for White.

	CurrentTurn Color // Color of the player whose turn it is.

	LastUpdate time.Time // Timestamp of the last update; -1 if the timer is stopped.

	timeoutCallback func(timeoutColor Color) // Timeout callback
	activeTimer     *time.Timer              // Internal timer for countdown.

	duration time.Duration // time of game

	mx sync.Mutex
}

// NewTimer creates a new Timer with the specified initial time for each player.
// Writes initial time to both players' clocks
func NewTimer(initialTimeSeconds int, increaseDuration time.Duration, turn Color, timeoutCallback func(timeoutColor Color)) *timer {
	return &timer{
		InitialTimeSeconds: initialTimeSeconds,
		IncreaseDuration:   increaseDuration,
		BlackTime:          time.Duration(initialTimeSeconds) * time.Second,
		WhiteTime:          time.Duration(initialTimeSeconds) * time.Second,
		CurrentTurn:        turn,
		LastUpdate:         NullTime,
		timeoutCallback:    timeoutCallback,
	}
}

func (t *timer) HasStarted() bool {
	t.mx.Lock()
	defer t.mx.Unlock()

	return t.duration == 0 && t.LastUpdate.Equal(NullTime)
}

func (t *timer) isStopped() bool { return t.LastUpdate.Equal(NullTime) }
func (t *timer) IsStopped() bool {
	t.mx.Lock()
	defer t.mx.Unlock()

	return t.isStopped()
}

func (t *timer) isRunning() bool { return t.LastUpdate != NullTime }
func (t *timer) IsRunning() bool {
	t.mx.Lock()
	defer t.mx.Unlock()

	return t.isRunning()
}

func (t *timer) updateTime() {
	if t.LastUpdate.Equal(NullTime) {
		return
	}

	now := time.Now()
	elapsed := now.Sub(t.LastUpdate)

	if t.CurrentTurn == White {
		t.WhiteTime -= elapsed
		if t.WhiteTime < 0 {
			t.WhiteTime = 0
		}
	} else {
		t.BlackTime -= elapsed
		if t.BlackTime < 0 {
			t.BlackTime = 0
		}
	}
	t.duration += elapsed
	t.LastUpdate = now
}

func (t *timer) Stop() bool {
	t.mx.Lock()
	defer t.mx.Unlock()

	if t.isStopped() {
		return false
	}

	if t.activeTimer != nil {
		t.activeTimer.Stop()
	}
	t.activeTimer = nil

	t.updateTime()
	t.LastUpdate = NullTime

	return true
}

func (t *timer) Start() bool {
	t.mx.Lock()
	defer t.mx.Unlock()

	if t.isRunning() {
		return false
	}

	t.LastUpdate = time.Now()
	t.setTimeout()
	return true
}

func (t *timer) GetWhiteDuration() time.Duration {
	t.mx.Lock()
	defer t.mx.Unlock()

	if t.isRunning() && t.CurrentTurn == White {
		elapsed := time.Since(t.LastUpdate)
		return t.WhiteTime - elapsed
	}

	return t.WhiteTime
}

func (t *timer) GetBlackDuration() time.Duration {
	t.mx.Lock()
	defer t.mx.Unlock()

	if t.isRunning() && t.CurrentTurn == Black {
		elapsed := time.Since(t.LastUpdate)
		return t.BlackTime - elapsed
	}

	return t.BlackTime
}

func (t *timer) GetTimes(color Color) time.Duration {
	switch color {
	case Black:
		return t.GetBlackDuration()
	case White:
		return t.GetWhiteDuration()
	default:
		return -1
	}
}

func (t *timer) GetCurrentPlayerTime() time.Duration {
	switch t.CurrentTurn {
	case Black:
		return t.GetBlackDuration()
	case White:
		return t.GetWhiteDuration()
	default:
		return -1
	}
}

func (t *timer) SwitchTurn() bool {
	t.mx.Lock()
	defer t.mx.Unlock()

	if t.isStopped() {
		return false
	}

	t.updateTime()
	if t.BlackTime == 0 || t.WhiteTime == 0 { // timeout
		t.LastUpdate = NullTime
		return false
	}

	if t.CurrentTurn == White {
		t.WhiteTime += t.IncreaseDuration
	} else {
		t.BlackTime += t.IncreaseDuration
	}

	t.CurrentTurn = t.CurrentTurn.Opposite()

	t.setTimeout()
	return true
}

func (t *timer) HasFlagged() bool {
	t.mx.Lock()
	defer t.mx.Unlock()

	if t.isStopped() {
		return t.WhiteTime == 0 || t.BlackTime == 0
	}

	elapsed := time.Since(t.LastUpdate)
	if t.CurrentTurn == Black {
		return t.WhiteTime == 0 || t.BlackTime-elapsed <= 0
	} else {
		return t.WhiteTime-elapsed <= 0 || t.BlackTime == 0
	}
}

func (t *timer) GetWinnerOnFlag() (Color, bool) {
	if t.HasFlagged() {
		return t.CurrentTurn, true
	}
	return None, false
}

// get current duration
func (t *timer) GetDuration() time.Duration {
	t.mx.Lock()
	defer t.mx.Unlock()

	if t.isRunning() {
		elapsed := time.Since(t.LastUpdate)
		return t.duration + elapsed
	}

	return t.duration
}

func (t *timer) handleTimeout() {
	t.mx.Lock()

	t.updateTime()
	t.LastUpdate = NullTime

	var timeoutColor Color

	switch {
	case t.BlackTime == 0:
		timeoutColor = Black
	case t.WhiteTime == 0:
		timeoutColor = White
	default: // do nothing
		t.mx.Unlock()
		return
	}
	t.mx.Unlock()

	if t.timeoutCallback != nil {
		t.timeoutCallback(timeoutColor)
	}
}

func (t *timer) clearTimeout() {
	if t.activeTimer != nil {
		t.activeTimer.Stop()
		t.activeTimer = nil
	}
}

func (t *timer) setTimeout() {
	t.clearTimeout()
	if t.CurrentTurn == White {
		t.activeTimer = time.AfterFunc(t.WhiteTime+time.Millisecond*50, t.handleTimeout) // add more 50ms
	} else {
		t.activeTimer = time.AfterFunc(t.BlackTime+time.Millisecond*50, t.handleTimeout) // add more 50ms
	}
}
