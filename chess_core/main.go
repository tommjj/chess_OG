package main

import (
	"fmt"
	"log"
)

var fens = []string{
	"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
	"r1bq1rk1/ppp2ppp/2n2n2/3pp3/3PP3/2N2N2/PPP2PPP/R1BQ1RK1 w - - 0 1",
	"r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1",
	"rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1",
	"8/P7/8/8/8/8/8/k6K w - - 0 1",
	"k7/8/8/8/8/8/7p/6K1 b - - 0 1",
	"r4rk1/1pp1qppp/p1np1n2/4p3/2B1P3/2N2Q2/PPPP1PPP/R1B2RK1 w - - 0 1",
	"8/8/8/8/8/7P/8/6kK w - - 0 1",
	"8/8/8/8/8/8/8/4K2k w - - 0 1",
	"rnbq1bnr/pppp1ppp/8/4p3/3PP3/5N2/PPP2PPP/RNBQKB1R w KQ - 0 1",
	"r3k2r/pppq1ppp/2n1pn2/3p4/3P4/2N1PN2/PPP2PPP/R1BQ1RK1 w - - 0 1",
}

// Example usage
func main() {
	gs := &GameState{
		BitBoards: &BitBoards{},
	}
	err := gs.FromFEN("6k1/R7/1R6/8/8/8/8/K7 w - - 0 1")
	if err != nil {
		fmt.Println("Error parsing FEN:", err)
		return
	}
	fmt.Print(gs)

	r, err := gs.MakeMove(White, SquareB6, SquareB8, 0)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(r)

	fmt.Print(gs)
}
