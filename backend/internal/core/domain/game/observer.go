// this is a observer for game session

package game

import (
	"context"
	"sync"
)

// ********************************
// *********** Observer ***********
// ********************************

type Tick struct {
	Tick int
}

func (t *Tick) SetTick(tick int) {
	t.Tick = tick
}

func (t *Tick) GetTick() int {
	return t.Tick
}

type ISetTick interface {
	SetTick(int)
}

type observer[T ISetTick] struct {
	currentID int
	observers map[int]chan T
	tick      int

	ctx context.Context
	mu  sync.Mutex
}

func newObserver[T ISetTick](ctx context.Context) *observer[T] {
	return &observer[T]{
		observers: map[int]chan T{},
		ctx:       ctx,
	}
}

// Register
func (g *observer[T]) Register(ctx context.Context) <-chan T {
	g.mu.Lock()
	defer g.mu.Unlock()

	ch := make(chan T, 5)

	g.currentID++
	id := g.currentID
	g.observers[id] = ch

	go func(observerID int) {
		select {
		case <-ctx.Done():
		case <-g.ctx.Done():
		}
		g.Unregister(observerID)
	}(id)

	return ch
}

// Unregister
func (g *observer[T]) Unregister(id int) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if ch, ok := g.observers[id]; ok {
		close(ch) // Close channel
		delete(g.observers, id)
	}
}

func (g *observer[T]) Publish(event T) {
	go func() {
		g.mu.Lock()
		g.tick++
		event.SetTick(g.tick)
		copy := copyChanValue(g.observers)
		g.mu.Unlock()

		for _, c := range copy {
			select {
			case c <- event:
			default:
			}
		}
	}()
}

func copyChanValue[T ISetTick](m map[int]chan T) []chan T {
	chans := make([]chan T, len(m))

	for i, v := range m {
		chans[i] = v
	}

	return chans
}
