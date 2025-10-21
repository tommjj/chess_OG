package ws

import (
	"context"
	"sync"

	"github.com/google/uuid"
)

type ID uuid.UUID

func (id ID) String() string {
	return uuid.UUID(id).String()
}

// Emitter create a Emit for connections
type Emitter interface {
	// Emit a event to all connections
	Emit(ctx context.Context, event string, data any) error
}

type NilEmitter struct{}

func (n *NilEmitter) Emit(ctx context.Context, event string, data any) error { return nil }

// Omitter create a Emit that omit a connection
type Omitter interface {
	// Omit a event to all connections except the omit connection
	Omit(ctx context.Context, event string, data any) error
}

type NilOmitter struct{}

func (n *NilOmitter) Omit(ctx context.Context, event string, data any) error { return nil }

// hub manages websocket connections and rooms.
type hub struct {
	conns   map[ID]*Connection
	connsMu sync.RWMutex

	rooms   map[string]*room
	roomsMu sync.RWMutex
}

func NewWSHub() *hub {
	return &hub{
		conns: make(map[ID]*Connection),
		rooms: make(map[string]*room),
	}
}

// addConn add new conn
func (h *hub) addConn(conn *Connection) {
	h.conns[conn.ID()] = conn
}

// AddConn add new conn
func (h *hub) AddConn(conn *Connection) {
	h.connsMu.Lock()
	defer h.connsMu.Unlock()

	h.addConn(conn)
}

// Remove conn by id
func (h *hub) delConn(id ID) {
	delete(h.conns, id)
}

// leaveAllRoom makes the connection leave all joined rooms.
func (h *hub) leaveAllRoom(conn *Connection) {
	roomNames := conn.JoinedRooms()

	h.roomsMu.Lock()
	defer h.roomsMu.Unlock()

	for _, name := range roomNames {
		h.leaveRoom(name, conn)
	}
}

// LeaveAllRoom makes the connection leave all joined rooms.
func (h *hub) LeaveAllRoom(conn *Connection) {
	h.leaveAllRoom(conn)
}

// Remove conn
// leave all rooms before remove
func (h *hub) DelConn(conn *Connection) {
	h.connsMu.Lock()
	defer h.connsMu.Unlock()

	h.leaveAllRoom(conn)

	h.delConn(conn.ID())
}

// joinRoom makes the connection join the specified room.
func (h *hub) joinRoom(name string, conn *Connection) {
	room, ok := h.rooms[name]
	if !ok {
		room = NewRoom()
	}

	room.Add(conn)
	h.rooms[name] = room
	conn.addJoinedRoom(name)
}

// JoinRoom makes the connection join the specified room.
func (h *hub) JoinRoom(name string, conn *Connection) {
	h.roomsMu.Lock()
	defer h.roomsMu.Unlock()

	h.joinRoom(name, conn)
}

// leaveRoom makes the connection leave the specified room.
func (h *hub) leaveRoom(name string, conn *Connection) {
	room, ok := h.rooms[name]
	if !ok {
		return
	}

	room.Remove(conn)
	conn.removeJoinedRoom(name)
	if room.IsIsEmpty() {
		delete(h.rooms, name)
	}
}

// LeaveRoom makes the connection leave the specified room.
func (h *hub) LeaveRoom(name string, conn *Connection) {
	h.roomsMu.Lock()
	defer h.roomsMu.Unlock()

	h.leaveRoom(name, conn)
}

// ToRoom returns an Emitter that sends events to all connections in the specified room.
func (h *hub) ToRoom(name string) Emitter {
	h.roomsMu.RLock()
	defer h.roomsMu.RUnlock()

	room, ok := h.rooms[name]
	if !ok {
		return &NilEmitter{}
	}

	return &RoomEmitter{
		room: room,
	}
}

// ToRoomOmit returns an Omitter that sends events to all connections in the specified room except the specified connection.
func (h *hub) ToRoomOmit(name string, omitConn *Connection) Omitter {
	h.roomsMu.RLock()
	defer h.roomsMu.RUnlock()

	room, ok := h.rooms[name]
	if !ok {
		return &NilOmitter{}
	}

	return &RoomOmitter{
		room: room,
		conn: omitConn,
	}
}

func (h *hub) ToConn(id ID) Emitter {
	h.connsMu.RLock()
	defer h.connsMu.RUnlock()

	conn, ok := h.conns[id]
	if !ok {
		return &NilEmitter{}
	}

	return &ConnEmitter{
		conn: conn,
	}
}

func (h *hub) Size() int {
	h.connsMu.RLock()
	defer h.connsMu.RUnlock()
	return len(h.conns)
}

func (h *hub) RoomsSize() int {
	h.roomsMu.RLock()
	defer h.roomsMu.RUnlock()
	return len(h.rooms)
}
