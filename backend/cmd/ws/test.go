package main

import (
	"fmt"
	"net/http"

	"github.com/tommjj/chess_OG/backend/internal/core/domain/game"
	"github.com/tommjj/chess_OG/backend/internal/interface/ws"
	"github.com/tommjj/chess_OG/backend/internal/web"
)

func main() {
	server := http.NewServeMux()
	clientFS := http.FileServer(http.FS(web.ClientFS))

	e := ws.NewEventHandler()

	var state *game.GameState
	_ = state

	e.Register("new", func(ctx *ws.Context) {
		var err error
		state, err = game.NewGame("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", 300000, func(result game.GameResult) {
			fmt.Println("Game ended with result:", result)
		})
		if err != nil {
			fmt.Println("Error creating new game:", err)
			return
		}

		ctx.Emit(ctx, "new", "Game created successfully")
	})

	e.Register("hello", func(ctx *ws.Context) {
		ctx.Emit(ctx, "hello", "Hello from server!")
	})

	e.Register("join", func(ctx *ws.Context) {
		var room string
		err := ctx.BindJSON(&room)
		if err != nil {
			fmt.Println("Error binding JSON:", err)
			return
		}
		ctx.Join(room)

		fmt.Println("Connection", ctx.Conn.ID().String(), "joined room:", room)
	})

	type Message struct {
		To   string `json:"to"`
		Mess string `json:"mess"`
	}

	e.Register("message", func(ctx *ws.Context) {
		var msg Message
		err := ctx.BindJSON(&msg)
		if err != nil {
			fmt.Println("Error binding JSON:", err)
			return
		}

		ctx.ToRoom(msg.To).Emit(ctx, "message", msg.Mess)

		fmt.Println("Received message from", ctx.Conn.ID().String(), ":", msg.Mess)
	})

	e.Register("disconnection", func(ctx *ws.Context) {
		ctx.Close("disconnection event")
	})

	handler := ws.NewHandler(ws.NewWSHub(), e, ws.WithOnConnect(func(ctx *ws.Context) {
		fmt.Println("New connection established with ID:", ctx.Conn.ID())

		copyCtx := ctx.Clone()

		go func() {
			<-copyCtx.ConnCtx().Done()
			fmt.Println("Connection close with ID:", ctx.Conn.ID())

		}()
	}), ws.WithOriginPatterns([]string{"localhost:5173", "127.0.0.1:5173", "127.0.0.1:8080", "localhost:8080"}))

	server.Handle("/ws", handler)

	server.Handle("GET /assets/", clientFS)

	server.HandleFunc("/", justAllowMethod(http.MethodGet, func(w http.ResponseWriter, r *http.Request) {
		// just serve index.html for all other routes
		http.ServeFileFS(w, r, web.ClientFS, "index.html")
	}))

	http.ListenAndServe(":8080", server)
}

func justAllowMethod(method string, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			http.NotFound(w, r)
			return
		}

		handler(w, r)
	}
}
