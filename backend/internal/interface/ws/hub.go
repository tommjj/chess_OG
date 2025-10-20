package ws

import (
	"context"
	"sync"

	"github.com/google/uuid"
)

type ID uuid.UUID

// Emiter là interface định nghĩa một phương thức Emit.
type Emiter interface {
	// Emit gửi dữ liệu với tên sự kiện, đồng thời kiểm tra context để hủy bỏ (cancel)
	// hoặc giới hạn thời gian (timeout) nếu cần.
	Emit(ctx context.Context, event string, data any) error
}

type WSHub struct {
	conns   map[ID]*Connection
	connsMu sync.RWMutex

	rooms   map[string]*Room
	roomsMu sync.RWMutex
}

func NewWSHub() *WSHub {
	return &WSHub{
		conns: make(map[ID]*Connection),
		rooms: make(map[string]*Room),
	}
}

func (h *WSHub) addConn(conn *Connection) {
	h.conns[conn.ID()] = conn
}

func (h *WSHub) AddConn(conn *Connection) {
	h.connsMu.Lock()
	defer h.connsMu.Unlock()

	h.addConn(conn)
}

// Remove conn by id
func (h *WSHub) delConn(id ID) {
	delete(h.conns, id)
}

func (h *WSHub) leaveAllRoom(conn *Connection) {
	roomNames := conn.JoinedRooms()

	h.connsMu.Lock()
	defer h.connsMu.Unlock()

	for _, name := range roomNames {
		h.leaveRoom(name, conn)
	}
}

func (h *WSHub) LeaveAllRoom(conn *Connection) {
	h.leaveAllRoom(conn)
}

// Remove conn
// leave all rooms before remove
func (h *WSHub) DelConn(conn *Connection) {
	h.connsMu.Lock()
	defer h.connsMu.Unlock()

	h.leaveAllRoom(conn)

	h.delConn(conn.ID())
}

func (h *WSHub) joinRoom(name string, conn *Connection) {
	room, ok := h.rooms[name]
	if !ok {
		room = NewRoom()
	}

	room.Add(conn)
	h.rooms[name] = room
	conn.addJoinedRoom(name)
}

func (h *WSHub) JoinRoom(name string, conn *Connection) {
	h.roomsMu.Lock()
	defer h.roomsMu.Unlock()

	h.joinRoom(name, conn)
}

func (h *WSHub) leaveRoom(name string, conn *Connection) {
	room, ok := h.rooms[name]
	if !ok {
		return
	}

	room.Remove(conn)
	if room.IsIsEmpty() {
		delete(h.rooms, name)
	}
}

func (h *WSHub) LeaveRoom(name string, conn *Connection) {
	h.roomsMu.Lock()
	defer h.roomsMu.Unlock()

	h.leaveRoom(name, conn)
}
