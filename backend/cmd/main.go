package main

import (
	"fmt"
	"log"

	chess "github.com/tommjj/chess_OG/chess_core"
)

func main() {
	gs := chess.GameState{
		BitBoards: &chess.BitBoards{},
	}

	err := gs.FromFEN(chess.StartingFEN)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(gs.String())

	result, err := gs.MakeMove(chess.White, chess.SquareA2, chess.SquareA4, 0)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result)

	fmt.Println(gs.String())
}
