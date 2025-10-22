package main

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"

	"github.com/tommjj/chess_OG/backend/internal/interface/ws"
	"github.com/tommjj/chess_OG/backend/internal/web"
)

func main() {
	sub, err := fs.Sub(web.StaticFiles, ".")
	if err != nil {
		log.Fatalf("failed to get sub FS: %v", err)
	}

	server := http.NewServeMux()

	httpFS := http.FS(sub)
	fs := http.FileServer(httpFS)

	hub := ws.NewWSHub()
	eventHandler := ws.NewEventHandler()

	eventHandler.Register("hello", func(ctx *ws.Context) {
		fmt.Println("Received hello event with data:", ctx.Conn.ID())

		fmt.Println("conn size", ctx.Hub.Size())
		fmt.Println("room size", ctx.Hub.RoomsSize())

	})

	eventHandler.Register("join", func(ctx *ws.Context) {
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

	eventHandler.Register("message", func(ctx *ws.Context) {
		var msg Message
		err := ctx.BindJSON(&msg)
		if err != nil {
			fmt.Println("Error binding JSON:", err)
			return
		}

		ctx.ToRoom(msg.To).Emit(ctx, "message", msg.Mess)

		fmt.Println("Received message from", ctx.Conn.ID().String(), ":", msg.Mess)
	})

	handler := ws.NewHandler(hub, eventHandler, ws.WithOnConnect(func(ctx *ws.Context) {
		fmt.Println("New connection established with ID:", ctx.Conn.ID())
	}))

	server.Handle("/ws", handler)

	server.Handle("/", fs)

	http.ListenAndServe(":8080", server)
}
