package ws

import (
	"context"
	"encoding/json"

	"github.com/coder/websocket"
)

type Context struct {
	context.Context

	Conn *Connection
	Hub  *Hub

	Payload json.RawMessage
}

// Join makes the connection associated with this context join the specified room.
func (c *Context) Join(roomName string) {
	c.Hub.JoinRoom(roomName, c.Conn)
}

// Leave makes the connection associated with this context leave the specified room.
func (c *Context) Leave(roomName string) {
	c.Hub.LeaveRoom(roomName, c.Conn)
}

// LeaveAll makes the connection associated with this context leave all joined rooms.
func (c *Context) LeaveAll() {
	c.Hub.LeaveAllRoom(c.Conn)
}

// ToRoom returns an Emitter that sends events to all connections in the specified room.
func (c *Context) ToRoom(roomName string) Emitter {
	return c.Hub.ToRoom(roomName)
}

// ToRoomOmit returns an Omitter that sends events to all connections in the specified room except the connection associated with this context.
func (c *Context) ToRoomOmit(roomName string) Omitter {
	return c.Hub.ToRoomOmit(roomName, c.Conn)
}

// Emit sends an event with the given payload to the connection associated with this context.
func (c *Context) Emit(ctx context.Context, event string, payload any) error {
	return c.Conn.Emit(ctx, event, payload)
}

// Error sends an error message to the connection associated with this context.
func (c *Context) Error(ctx context.Context, mess string) error {
	return c.Emit(ctx, "error", map[string]string{"message": mess})
}

// CloseWithError sends an error message and then closes the connection associated with this context.
func (c *Context) CloseWithError(ctx context.Context, mess string) error {
	err := c.Error(ctx, mess)
	if err != nil {
		return err
	}

	return c.Close(mess)
}

// Close closes the connection associated with this context with the given reason.
func (c *Context) Close(reason string) error {
	return c.Conn.Close(websocket.StatusNormalClosure, reason)
}

// ID returns the unique identifier of the connection associated with this context.
func (c *Context) ID() ID {
	return c.Conn.ID()
}

// Return connection context
func (c *Context) ConnCtx() context.Context {
	return c.Conn.Ctx()
}

// Set sets a key-value pair in the connection's store.
func (c *Context) Set(key string, value any) {
	c.Conn.Set(key, value)
}

// LoadOrStore loads the value for a key if it exists, or stores and returns the given value otherwise.
func (c *Context) LoadOrStore(key string, value any) (any, bool) {
	return c.Conn.LoadOrStore(key, value)
}

// Get retrieves the value for a key from the connection's store.
func (c *Context) Get(key string) (any, bool) {
	return c.Conn.Get(key)
}

// MustGet retrieves the value for a key from the connection's store, panicking if the key does not exist.
func (c *Context) MustGet(key string) any {
	return c.Conn.MustGet(key)
}

// Delete removes a key-value pair from the connection's store.
func (c *Context) Delete(key string) {
	c.Conn.Delete(key)
}

// Range iterates over all key-value pairs in the connection's store, calling the given function for each pair.
func (c *Context) Range(f func(key any, value any) bool) {
	c.Conn.Range(f)
}

// BindJSON unmarshals the JSON payload of the context into the provided variable.
func (c *Context) BindJSON(v any) error {
	return json.Unmarshal(c.Payload, v)
}
