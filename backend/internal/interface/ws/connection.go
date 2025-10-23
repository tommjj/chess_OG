package ws

import (
	"context"
	"fmt"
	"sync"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"golang.org/x/time/rate"
)

type Connection struct {
	*websocket.Conn

	ctx context.Context

	id    ID       // unique identifier of the connection
	store sync.Map // key-value store for connection-specific data

	// joined rooms
	rooms  map[string]struct{}
	roomMu sync.Mutex

	rateLimitBucket *rate.Limiter // rate limiter for the connection
}

func NewWSConnection(ctx context.Context, id ID, conn *websocket.Conn) *Connection {
	limiter := rate.NewLimiter(5, 10)

	return &Connection{
		Conn:            conn,
		ctx:             ctx,
		id:              id,
		store:           sync.Map{},
		rateLimitBucket: limiter,

		rooms: make(map[string]struct{}),
	}
}

func (c *Connection) Set(key string, value any) {
	c.store.Store(key, value)
}

func (c *Connection) LoadOrStore(key string, value any) (any, bool) {
	return c.store.LoadOrStore(key, value)
}

func (c *Connection) Get(key string) (any, bool) {
	return c.store.Load(key)
}

func (c *Connection) MustGet(key string) any {
	v, ok := c.store.Load(key)
	if !ok {
		panic(fmt.Sprintf("MustGet: key %s is not exist", key))
	}
	return v
}

func (c *Connection) Delete(key string) {
	c.store.Delete(key)
}

func (c *Connection) Range(f func(key any, value any) bool) {
	c.store.Range(f)
}

func (c *Connection) ID() ID {
	return c.id
}

func (c *Connection) Ctx() context.Context {
	return c.ctx
}

func (c *Connection) addJoinedRoom(name string) {
	c.roomMu.Lock()
	defer c.roomMu.Unlock()

	c.rooms[name] = struct{}{}
}

func (c *Connection) removeJoinedRoom(name string) {
	c.roomMu.Lock()
	defer c.roomMu.Unlock()

	delete(c.rooms, name)
}

func (c *Connection) IsJoinedRoom(name string) bool {
	c.roomMu.Lock()
	defer c.roomMu.Unlock()

	_, ok := c.rooms[name]
	return ok
}

func (c *Connection) JoinedRooms() []string {
	c.roomMu.Lock()
	defer c.roomMu.Unlock()

	names := make([]string, 0, len(c.rooms))
	for name := range c.rooms {
		names = append(names, name)
	}

	return names
}

func (c *Connection) Allow() bool {
	return c.rateLimitBucket.Allow()
}

func (c *Connection) SetLimit(r rate.Limit, b int) {
	c.rateLimitBucket.SetLimit(r)
	c.rateLimitBucket.SetBurst(b)
}

func (c *Connection) Limit() (rate.Limit, int) {
	return c.rateLimitBucket.Limit(), c.rateLimitBucket.Burst()
}

func (c *Connection) Emit(ctx context.Context, event string, payload any) error {
	mess := Message{
		Event:   event,
		Payload: payload,
	}

	return wsjson.Write(ctx, c.Conn, mess)
}

type ConnEmitter struct {
	conn *Connection
}

func (ce *ConnEmitter) Emit(ctx context.Context, event string, payload any) error {
	return ce.conn.Emit(ctx, event, payload)
}
