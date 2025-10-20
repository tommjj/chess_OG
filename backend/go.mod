module github.com/tommjj/chess_OG/backend

go 1.24.1

replace github.com/tommjj/chess_OG/chess_core => ../chess_core

require (
	github.com/google/uuid v1.6.0
	github.com/tommjj/chess_OG/chess_core v0.0.0-00010101000000-000000000000
)

require (
	github.com/coder/websocket v1.8.14
	golang.org/x/time v0.14.0
)
