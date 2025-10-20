package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tommjj/chess_OG/backend/internal/interface/ws"
)

func main() {
	server := http.NewServeMux()

	fs := http.FileServer(http.Dir("./static"))

	hub := ws.NewWSHub()
	eventHandler := ws.NewEventHandler()

	eventHandler.Register("hello", func(ctx *ws.Context, data json.RawMessage) error {
		fmt.Println("Received hello event with data:", ctx.Conn.ID())

		return nil
	})

	handler := ws.NewHandler(hub, eventHandler)

	server.Handle("/ws", handler)

	server.Handle("/", fs)

	http.ListenAndServe(":8080", server)
}
