package main

import (
	"fmt"
	"net/http"

	"github.com/tommjj/chess_OG/backend/internal/interface/ws"
	"github.com/tommjj/chess_OG/backend/internal/web"
)

func main() {
	server := http.NewServeMux()
	fs := http.FileServer(http.FS(web.StaticFS))

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

	handler := ws.NewHandler(ws.NewWSHub(), eventHandler, ws.WithOnConnect(func(ctx *ws.Context) {
		fmt.Println("New connection established with ID:", ctx.Conn.ID())

		go func() {
			<-ctx.ConnCtx().Done()
			fmt.Println("connection close with ID:", ctx.Conn.ID())
		}()
	}))

	server.Handle("/ws", handler)

	server.Handle("/", fs)

	http.ListenAndServe(":8080", server)
}
