package ws

import (
	"fmt"
	"sync"

	"github.com/coder/websocket"
	"golang.org/x/time/rate"
)

type Connection struct {
	*websocket.Conn

	//
	id    ID
	store sync.Map

	// rooms joined name
	rooms  map[string]struct{}
	roomMu sync.Mutex

	rateLimitBucket *rate.Limiter
}

func NewWSConnection(id ID, conn *websocket.Conn) Connection {
	limiter := rate.NewLimiter(5, 10)

	return Connection{
		Conn:            conn,
		id:              id,
		store:           sync.Map{},
		rateLimitBucket: limiter,
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
