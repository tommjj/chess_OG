// Package ws implements websocket handler and connection management.
// It provides a way to handle websocket connections, manage rooms, and emit events to connections.
// It uses the github.com/coder/websocket package for websocket handling.

package ws

import (
	"context"
	"net/http"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/google/uuid"
)

// OptionsFunc defines a function to set options on the Handler.
type OptionsFunc func(*Handler)

// MiddlewareFunc defines a function to process middleware.
// returns an error to stop upgrading the connection.
type MiddlewareFunc func(conn *Connection, r *http.Request) error

// Handler handles websocket connections and dispatches events to the appropriate handlers.
type Handler struct {
	hub          *hub
	eventHandler *EventHandler

	// optional fields
	ider           func(r *http.Request) (ID, error) // function to generate connection ID from http request
	originPatterns []string                          // allowed origin patterns for websocket connections
	healthInterval time.Duration                     // interval for health check pings

	// event handling options
	eventTimeout   time.Duration // timeout for each event handler
	eventSemaphore int           // max concurrent event handlers for a connection

	// middlewares to run when a new connection is established
	middlewares  []MiddlewareFunc
	onConnect    HandleFunc
	onDisconnect HandleFunc
}

func NewHandler(hub *hub, eventHandler *EventHandler, ops ...OptionsFunc) *Handler {
	h := &Handler{
		hub:            hub,
		eventHandler:   eventHandler,
		ider:           DefaultIDer,
		eventSemaphore: 5,
		healthInterval: 30 * time.Second,

		middlewares: make([]MiddlewareFunc, 0),
	}

	for _, op := range ops {
		op(h)
	}

	return h
}

func (h *Handler) IDer() func(r *http.Request) (ID, error) {
	return h.ider
}

func (h *Handler) EventTimeout() time.Duration {
	return h.eventTimeout
}

func (h *Handler) Middlewares() []MiddlewareFunc {
	return h.middlewares
}

// Use add a middleware to the handler
func (h *Handler) Use(middleware MiddlewareFunc) {
	h.middlewares = append(h.middlewares, middleware)
}

// runReadLoop reads messages from the websocket connection and dispatches them to the appropriate event handlers.
// It blocks until the connection is closed.
func (h *Handler) runReadLoop(ctx context.Context, conn *Connection) {
	sem := make(chan struct{}, h.eventSemaphore)

	for {
		var event MessageSchema
		if err := wsjson.Read(ctx, conn.Conn, &event); err != nil {
			return
		}

		if !conn.Allow() {
			continue
		}

		handleFunc, ok := h.eventHandler.Get(event.Event)
		if !ok {
			continue
		}

		eventCtx, cancel := h.newContext(ctx, conn)
		eventCtx.Payload = event.Payload

		sem <- struct{}{}
		go func() {
			defer func() { <-sem }()
			defer cancel()
			defer releaseContext(eventCtx)

			handleFunc(eventCtx)
		}()
	}
}

// newContext creates a new Context for handling an event.
func (h *Handler) newContext(ctx context.Context, conn *Connection) (*Context, context.CancelFunc) {
	var baseCtx context.Context = ctx
	var cancel context.CancelFunc

	if h.eventTimeout > 0 {
		baseCtx, cancel = context.WithTimeout(ctx, h.eventTimeout)
	} else {
		baseCtx, cancel = context.WithCancel(ctx)
	}

	eventCtx := acquireContext()
	eventCtx.Context = baseCtx
	eventCtx.Conn = conn
	eventCtx.Hub = h.hub

	return eventCtx, cancel
}

// ServeHTTP handles the websocket connection upgrade and manages the connection lifecycle.
// It runs middlewares, adds the connection to the hub, and starts the read loop.
// this implements http.Handler interface.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// upgrader the http connection to websocket connection
	wsConn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns: h.originPatterns,
	})
	if err != nil {
		return
	}
	defer wsConn.CloseNow()

	connID, err := h.ider(r)
	if err != nil {
		wsConn.Close(websocket.StatusPolicyViolation, "invalid connection id:"+err.Error())
		return
	}

	conn := NewWSConnection(connID, wsConn)

	// run middlewares
	for _, m := range h.middlewares {
		if err := m(&conn, r); err != nil {
			wsConn.Close(websocket.StatusPolicyViolation, "middleware error:"+err.Error())
			return
		}
	}

	// add connection to hub
	h.hub.AddConn(&conn)
	defer h.hub.DelConn(&conn)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// start health ping goroutine
	if h.healthInterval > 0 {
		go h.healthPing(ctx, &conn)
	}

	// call onConnect handler
	if h.onConnect != nil {
		c, cal := h.newContext(ctx, &conn)
		h.onConnect(c)
		cal()
	}

	// start read loop this will block until connection is closed
	h.runReadLoop(ctx, &conn)

	// call onDisconnect handler
	if h.onDisconnect != nil {
		c, cal := h.newContext(ctx, &conn)
		h.onDisconnect(c)
		cal()
	}

	// close websocket connection
	wsConn.Close(websocket.StatusNormalClosure, "")
}

// healthPing sends ping messages to the client every 30 seconds to keep the connection alive.
// If the ping fails, the connection is closed.
func (h *Handler) healthPing(ctx context.Context, conn *Connection) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(30 * time.Second):
			pingCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
			err := conn.Ping(pingCtx)
			cancel()
			if err != nil {
				conn.Close(websocket.StatusNormalClosure, "ping error:"+err.Error())
				return
			}
		}
	}
}

// *******************************************
// *************  Options Funcs **************
// *******************************************

// DefaultIDer generates a random UUID as connection ID.
func DefaultIDer(r *http.Request) (ID, error) {
	return ID(uuid.New()), nil
}

// WithIDer sets the function to generate connection ID from http request.
// If returned error is not nil, the connection is rejected.
func WithIDer(ider func(r *http.Request) (ID, error)) OptionsFunc {
	return func(h *Handler) {
		h.ider = ider
	}
}

// WithEventTimeout sets the timeout for each event handler.
func WithEventTimeout(d time.Duration) OptionsFunc {
	return func(h *Handler) {
		h.eventTimeout = d
	}
}

// WithMiddleware adds a middleware to the handler. middleware runs orderly in the order they are added.
//
// Middlewares run when a new connection is established.
// If a middleware returns an error, the connection is rejected.
func WithMiddleware(middleware MiddlewareFunc) OptionsFunc {
	return func(h *Handler) {
		h.middlewares = append(h.middlewares, middleware)
	}
}

// WithEventSemaphore sets the maximum number of concurrent event handlers for a connection.
func WithEventSemaphore(n int) OptionsFunc {
	return func(h *Handler) {
		h.eventSemaphore = n
	}
}

// WithOriginPatterns sets the allowed origin patterns for websocket connections.
func WithOriginPatterns(patterns []string) OptionsFunc {
	return func(h *Handler) {
		h.originPatterns = patterns
	}
}

// WithHealthInterval sets the interval for health check pings.
// Default is 30 seconds.
// If set to 0, health check pings are disabled.
func WithHealthInterval(d time.Duration) OptionsFunc {
	return func(h *Handler) {
		h.healthInterval = d
	}
}

// WithOnConnect sets the handler function to be called when a connection is established.
// The handler function receives the Context and nil data.
// This calls after the connection is upgraded and before the read loop starts. you can send messages to the connection in this handler.
func WithOnConnect(handleFunc HandleFunc) OptionsFunc {
	return func(h *Handler) {
		h.onConnect = handleFunc
	}
}

// WithOnDisconnect sets the handler function to be called when a connection is disconnected.
// The handler function receives the Context and nil data.
// This calls after the read loop ends and before the connection is closed. you can't send messages to the connection in this handler.
// error returned from the handler is ignored.
func WithOnDisconnect(handleFunc HandleFunc) OptionsFunc {
	return func(h *Handler) {
		h.onDisconnect = handleFunc
	}
}
