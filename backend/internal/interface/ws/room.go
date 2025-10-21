package ws

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/coder/websocket"
)

// room represents a group of connections.
type room struct {
	conns   map[*Connection]struct{}
	connsMu sync.RWMutex
}

func NewRoom() *room {
	return &room{
		conns: make(map[*Connection]struct{}),
	}
}

// Add new conn
func (r *room) Add(conn *Connection) {
	r.connsMu.Lock()
	defer r.connsMu.Unlock()

	r.conns[conn] = struct{}{}
}

// Remove remove conn
func (r *room) Remove(conn *Connection) {
	r.connsMu.Lock()
	defer r.connsMu.Unlock()

	delete(r.conns, conn)
}

// Size get room size
func (r *room) Size() int {
	r.connsMu.RLock()
	defer r.connsMu.RUnlock()
	return len(r.conns)
}

// IsIsEmpty check if room is empty
func (r *room) IsIsEmpty() bool {
	r.connsMu.RLock()
	defer r.connsMu.RUnlock()
	return len(r.conns) == 0
}

// RoomEmitter create a Emit for room
type RoomEmitter struct {
	room *room
}

// Emit a event to all connections in the room
func (r *RoomEmitter) Emit(ctx context.Context, event string, payload any) error {
	mess := Message{
		Event:   event,
		Payload: payload,
	}

	encodeMess, err := mess.Encode()
	if err != nil {
		return fmt.Errorf("failed to encode message: %w", err)
	}

	clientsCopy := *acquireConnSlice()
	defer releaseConnSlice(&clientsCopy)

	r.room.connsMu.RLock()

	roomSize := len(r.room.conns)
	if cap(clientsCopy) < roomSize {
		clientsCopy = make([]*Connection, 0, roomSize)
	}

	for conn := range r.room.conns {
		clientsCopy = append(clientsCopy, conn)
	}

	r.room.connsMu.RUnlock()

	var (
		wg     sync.WaitGroup
		errs   ConnErrors = nil
		errsMu sync.Mutex
	)

	for _, conn := range clientsCopy {
		wg.Add(1)
		go func(c *Connection) {
			defer wg.Done()

			err := c.Write(ctx, websocket.MessageText, encodeMess)
			if err != nil {
				errsMu.Lock()
				if errs == nil {
					errs = make(ConnErrors, 0, 1)
				}

				errs = append(errs, ConnError{
					ConnID: c.ID(),
					Op:     "write",
					Err:    err,
				})

				if errors.Is(err, io.EOF) || websocket.CloseStatus(err) != -1 {
					r.room.Remove(c) // thread-safe remove
				}

				errsMu.Unlock()
			}
		}(conn)
	}
	wg.Wait()

	if len(errs) != 0 {
		return errs
	}
	return nil
}

type RoomOmitter struct {
	room *room
	conn *Connection
}

// Omit a event to all connections in the room except the given conn
func (r *RoomOmitter) Omit(ctx context.Context, event string, payload any) error {
	mess := Message{
		Event:   event,
		Payload: payload,
	}

	encodeMess, err := mess.Encode()
	if err != nil {
		return fmt.Errorf("failed to encode message: %w", err)
	}

	clientsCopy := *acquireConnSlice()
	defer releaseConnSlice(&clientsCopy)

	r.room.connsMu.RLock()

	roomSize := len(r.room.conns)
	if cap(clientsCopy) < roomSize {
		clientsCopy = make([]*Connection, 0, roomSize)
	}

	for conn := range r.room.conns {
		if conn != r.conn {
			clientsCopy = append(clientsCopy, conn)
		}
	}

	r.room.connsMu.RUnlock()

	var (
		wg     sync.WaitGroup
		errs   ConnErrors = nil
		errsMu sync.Mutex
	)

	for _, conn := range clientsCopy {
		wg.Add(1)
		go func(c *Connection) {
			defer wg.Done()

			err := c.Write(ctx, websocket.MessageText, encodeMess)
			if err != nil {
				errsMu.Lock()
				if errs == nil {
					errs = make(ConnErrors, 0, 1)
				}

				errs = append(errs, ConnError{
					ConnID: c.ID(),
					Op:     "write",
					Err:    err,
				})

				if errors.Is(err, io.EOF) || websocket.CloseStatus(err) != -1 {
					r.room.Remove(c) // thread-safe remove
				}

				errsMu.Unlock()
			}
		}(conn)
	}
	wg.Wait()

	if len(errs) != 0 {
		return errs
	}
	return nil
}
