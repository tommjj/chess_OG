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
	for timer.SwitchTurn() {
	}

	t.Log(timer.GetDuration())

	t.Log(timer.WhiteRemaining())

	t.Log(timer.BlackRemaining())

	time.Sleep(time.Second)
}
