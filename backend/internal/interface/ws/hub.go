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

// Hub manages websocket connections and rooms.
type Hub struct {
	conns   map[ID]*Connection
	connsMu sync.RWMutex

	rooms   map[string]*Room
	roomsMu sync.RWMutex
}

func NewWSHub() *Hub {
	return &Hub{
		conns: make(map[ID]*Connection),
		rooms: make(map[string]*Room),
	}
}

// addConn add new conn
func (h *Hub) addConn(conn *Connection) {
	h.conns[conn.ID()] = conn
}

// AddConn add new conn
func (h *Hub) AddConn(conn *Connection) {
	h.connsMu.Lock()
	defer h.connsMu.Unlock()

	h.addConn(conn)
}

// Remove conn by id
func (h *Hub) delConn(id ID) {
	delete(h.conns, id)
}

// leaveAllRoom makes the connection leave all joined rooms.
func (h *Hub) leaveAllRoom(conn *Connection) {
	roomNames := conn.JoinedRooms()

	h.roomsMu.Lock()
	defer h.roomsMu.Unlock()

	for _, name := range roomNames {
		h.leaveRoom(name, conn)
	}
}

// LeaveAllRoom makes the connection leave all joined rooms.
func (h *Hub) LeaveAllRoom(conn *Connection) {
	h.leaveAllRoom(conn)
}

// Remove conn
// leave all rooms before remove
func (h *Hub) DelConn(conn *Connection) {
	h.connsMu.Lock()
	defer h.connsMu.Unlock()

	h.leaveAllRoom(conn)

	h.delConn(conn.ID())
}

// joinRoom makes the connection join the specified room.
func (h *Hub) joinRoom(name string, conn *Connection) {
	room, ok := h.rooms[name]
	if !ok {
		room = NewRoom()
	}

	room.Add(conn)
	h.rooms[name] = room
	conn.addJoinedRoom(name)
}

// JoinRoom makes the connection join the specified room.
func (h *Hub) JoinRoom(name string, conn *Connection) {
	h.roomsMu.Lock()
	defer h.roomsMu.Unlock()

	h.joinRoom(name, conn)
}

// leaveRoom makes the connection leave the specified room.
func (h *Hub) leaveRoom(name string, conn *Connection) {
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
func (h *Hub) LeaveRoom(name string, conn *Connection) {
	h.roomsMu.Lock()
	defer h.roomsMu.Unlock()

	h.leaveRoom(name, conn)
}

// ToRoom returns an Emitter that sends events to all connections in the specified room.
func (h *Hub) ToRoom(name string) Emitter {
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
func (h *Hub) ToRoomOmit(name string, omitConn *Connection) Omitter {
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

func (h *Hub) ToConn(id ID) Emitter {
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

func (h *Hub) Size() int {
	h.connsMu.RLock()
	defer h.connsMu.RUnlock()
	return len(h.conns)
}

func (h *Hub) RoomsSize() int {
	h.roomsMu.RLock()
	defer h.roomsMu.RUnlock()
	return len(h.rooms)
}
