package ws

import (
	"context"
	"time"
)

type Context struct {
	context.Context

	Conn *Connection
	Hub  *WSHub
}

func NewContextWithTimeout(ctx context.Context, conn *Connection, hub *WSHub, timeout time.Duration) (*Context, context.CancelFunc) {
	ctxTimeout, cal := context.WithTimeout(ctx, timeout)

	return &Context{
		Context: ctxTimeout,
		Conn:    conn,
		Hub:     hub,
	}, cal
}

func (c *Context) Join(roomName string) {
	c.Hub.JoinRoom(roomName, c.Conn)
}

func (c *Context) Leave(roomName string) {
	c.Hub.LeaveRoom(roomName, c.Conn)
}

func (c *Context) LeaveAll() {
	c.Hub.LeaveAllRoom(c.Conn)
}

// func (c *Context) Emit(event string, data any) error {
// 	// return c.Conn.Emit(c.Context, event, data)
// }
