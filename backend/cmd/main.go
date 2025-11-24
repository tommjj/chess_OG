package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	chess "github.com/tommjj/chess_OG/chess_core"
)

func main() {
	gs := chess.NewGame()

	err := gs.FromFEN(chess.StartingFEN)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(gs.String())

	fmt.Println(gs.String())

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()

		if line == "history" {
			fmt.Printf("%+v\n", gs.History())
			continue
		}

		if line == "undo" {
			gs.Undo(1)
			fmt.Println(gs.String())
			continue
		}

		from, to, ok := strings.Cut(line, " ")
		if !ok {
			fmt.Println("sai")
			continue
		}

		fromSq, ok := toSquare(from)
		if !ok {
			fmt.Println("sai")
			continue
		}

		toSq, ok := toSquare(to)
		if !ok {
			fmt.Println("sai")
			continue
		}

		result, err := gs.MakeMove(gs.SideToMove, fromSq, toSq, 0)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(result)
		fmt.Println(gs.String())
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Lỗi đọc:", err)
	}
}

func toSquare(p string) (chess.Square, bool) {
	if len(p) != 2 {
		return -1, false
	}
	p = strings.ToLower(p)

	if p[0] < 'a' || p[0] > 'h' {
		return -1, false
	}
	if p[1] < '1' || p[1] > '8' {
		return -1, false
	}

	return chess.Square((p[0] - 'a') + (p[1]-'1')*8), true

}
