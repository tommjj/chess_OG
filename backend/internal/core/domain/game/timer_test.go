package game

import (
	"testing"
	"time"
)

func TestTimer(t *testing.T) {
	timer := NewTimer(2, time.Second, White, func(timeoutColor Color) {
		t.Log(timeoutColor.String())
	})

	timer.Start()
	timer.SwitchTurn()

	time.Sleep(time.Second)
	timer.SwitchTurn()
	time.Sleep(time.Second)
	timer.SwitchTurn()

	time.Sleep(time.Second * 3)

	t.Log(timer.GetDuration())

	t.Log(timer.GetWhiteDuration())

	t.Log(timer.GetBlackDuration())

	time.Sleep(time.Second)
}
