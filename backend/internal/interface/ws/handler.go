package ws

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/google/uuid"
)

type OptionsFunc func(*Handler)

// MiddlewareFunc defines a function to process middleware.
// returns an error to stop upgrading the connection.
type MiddlewareFunc func(conn *Connection, r *http.Request) error

type Handler struct {
	Hub *WSHub

	EventHandler *EventHandler

	// optional fields
	ider func(r *http.Request) (ID, error)

	eventTimeout time.Duration

	middlewares []MiddlewareFunc
}

func NewHandler(hub *WSHub, eventHandler *EventHandler, ops ...OptionsFunc) *Handler {
	h := &Handler{
		Hub:          hub,
		EventHandler: eventHandler,
		ider:         DefaultIDer,

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

func (h *Handler) runReadLoop(ctx context.Context, conn *Connection) {
	sem := make(chan struct{}, 5)

	for {
		var event MessageSchema
		if err := wsjson.Read(ctx, conn.Conn, &event); err != nil {
			return
		}

		if !conn.Allow() {
			continue
		}

		handleFunc, ok := h.EventHandler.Get(event.Event)
		if !ok {
			continue
		}

		eventCtx, cancel := h.newContext(ctx, conn)

		sem <- struct{}{}
		go func() {
			defer func() { <-sem }()
			defer cancel()
			defer releaseContext(eventCtx)

			err := handleFunc(eventCtx, event.Payload)
			if err != nil {
				// emit error event
			}
		}()
	}
}

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
	eventCtx.Hub = h.Hub

	return eventCtx, cancel
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// upgrader the http connection to websocket connection
	wsConn, err := websocket.Accept(w, r, nil)
	if err != nil {
		fmt.Println("Failed to upgrade to websocket:", err)
		// http.Error(w, "Failed to upgrade to websocket:"+err.Error(), http.StatusBadRequest)
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

	h.Hub.AddConn(&conn)
	defer h.Hub.DelConn(&conn)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go h.healthPing(ctx, &conn)
	h.runReadLoop(ctx, &conn)

	wsConn.Close(websocket.StatusNormalClosure, "")
}

func (h *Handler) healthPing(ctx context.Context, conn *Connection) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(30 * time.Second):
			pingCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
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

func DefaultIDer(r *http.Request) (ID, error) {
	return ID(uuid.New()), nil
}

func WithIDer(ider func(r *http.Request) (ID, error)) OptionsFunc {
	return func(h *Handler) {
		h.ider = ider
	}
}

func WithEventTimeout(d time.Duration) OptionsFunc {
	return func(h *Handler) {
		h.eventTimeout = d
	}
}

func WithMiddleware(middleware MiddlewareFunc) OptionsFunc {
	return func(h *Handler) {
		h.middlewares = append(h.middlewares, middleware)
	}
}
