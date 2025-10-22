package game

import "time"

var NullTime = time.Time{}

// Timer represents a chess game timer.
type Timer struct {
	InitialTimeSeconds int // Initial time in seconds for each player.

	BlackTime time.Duration // Remaining time for Black.
	WhiteTime time.Duration // Remaining time for White.

	CurrentTurn Color // Color of the player whose turn it is.

	LastUpdate time.Time // Timestamp of the last update; -1 if the timer is stopped.

	TimeOutSignal chan Color  // Channel to signal when a player runs out of time.
	ActiveTimer   *time.Timer // Internal timer for countdown.
}

// NewTimer creates a new Timer with the specified initial time for each player.
// Writes initial time to both players' clocks and sets the current turn to White.
func NewTimer(initialTimeSeconds int) *Timer {
	return &Timer{
		InitialTimeSeconds: initialTimeSeconds,
		BlackTime:          time.Duration(initialTimeSeconds) * time.Second,
		WhiteTime:          time.Duration(initialTimeSeconds) * time.Second,
		CurrentTurn:        White,
		LastUpdate:         NullTime,

		TimeOutSignal: make(chan Color, 1),
	}
}

// Clone creates a deep copy of the Timer.
func (t *Timer) Clone() *Timer {
	return &Timer{
		InitialTimeSeconds: t.InitialTimeSeconds,
		BlackTime:          t.BlackTime,
		WhiteTime:          t.WhiteTime,
		CurrentTurn:        t.CurrentTurn,
		LastUpdate:         t.LastUpdate,
	}
}

func (t *Timer) IsStopped() bool {
	return t.LastUpdate.IsZero()
}

func (t *Timer) IsRunning() bool {
	return !t.LastUpdate.IsZero()
}

func (t *Timer) Stop() {
	if t.IsRunning() {
		if t.ActiveTimer != nil {
			t.ActiveTimer.Stop()
		}
		t.ActiveTimer = nil

		now := time.Now()
		elapsed := now.Sub(t.LastUpdate)

		if t.CurrentTurn == White {
			t.WhiteTime -= elapsed
		} else {
			t.BlackTime -= elapsed
		}

		t.LastUpdate = NullTime
	}
}

func (t *Timer) Start() {
	if t.IsStopped() {
		if t.ActiveTimer != nil {
			t.ActiveTimer.Stop()
		}

		if t.CurrentTurn == White {
			t.ActiveTimer = time.AfterFunc(t.BlackTime, func() {
				t.TimeOutSignal <- White
			})
		} else {
			t.ActiveTimer = time.AfterFunc(t.WhiteTime, func() {
				t.TimeOutSignal <- Black
			})
		}

		t.LastUpdate = time.Now()
	}
}

func (t *Timer) getWhiteTime() time.Duration {
	if t.IsRunning() && t.CurrentTurn == White {
		elapsed := time.Since(t.LastUpdate)
		return t.WhiteTime - elapsed
	}
	return t.WhiteTime
}

func (t *Timer) getBlackTime() time.Duration {
	if t.IsRunning() && t.CurrentTurn == Black {
		elapsed := time.Since(t.LastUpdate)
		return t.BlackTime - elapsed
	}
	return t.BlackTime
}

func (t *Timer) GetCurrentPlayerTime() time.Duration {
	if t.CurrentTurn == White {
		return t.getWhiteTime()
	}
	return t.getBlackTime()
}

func (t *Timer) SwitchTurn() {
	if t.IsStopped() {
		return
	}

	// clear timer
	if t.ActiveTimer != nil {
		t.ActiveTimer.Stop()
	}

	now := time.Now()
	elapsed := now.Sub(t.LastUpdate)

	if t.CurrentTurn == White {
		t.WhiteTime -= elapsed
		t.CurrentTurn = Black

		t.ActiveTimer = time.AfterFunc(t.BlackTime, func() {
			t.TimeOutSignal <- Black
		})
	} else {
		t.BlackTime -= elapsed
		t.CurrentTurn = White

		t.ActiveTimer = time.AfterFunc(t.WhiteTime, func() {
			t.TimeOutSignal <- White
		})
	}

	t.LastUpdate = now
}

func (t *Timer) HasFlagged() bool {
	return t.getWhiteTime() <= 0 || t.getBlackTime() <= 0
}

func (t *Timer) GetWinnerOnFlag() Color {
	if t.getWhiteTime() <= 0 {
		return Black
	} else if t.getBlackTime() <= 0 {
		return White
	}
	return Both
}

func (t *Timer) GetTimes(color Color) time.Duration {
	if color == White {
		return t.getWhiteTime()
	}
	return t.getBlackTime()
}
