package game

import (
	"testing"
	"time"
)

func TestTimer(t *testing.T) {
	timer := NewTimer(2, 0, White, func(timeoutColor Color) {
		t.Log(timeoutColor.String())
	})

	timer.Start()
	timer.SwitchTurn()

	for timer.SwitchTurn() {
	}

	t.Log(timer.GetDuration())
	t.Log(timer.GetWhiteDuration())
	t.Log(timer.GetBlackDuration())
	time.Sleep(time.Second)
}
