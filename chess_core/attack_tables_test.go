package chess_core

import (
	"testing"
)

var testFENs = []string{
	"r6k/8/8/8/8/8/8/K7 w KQkq - 0 1",
	"k7/8/8/8/8/8/8/R3K2R w KQ - 0 1",
	"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
	"4k3/8/8/8/8/8/8/4K3 w - - 0 1",
}

func BenchmarkIsKingAttacked(b *testing.B) {

	for _, fen := range testFENs {
		bbs := &BitBoards{}
		bbs.FromFEN(fen)

		b.Run(fen, func(b *testing.B) {
			initAttackTables()

			b.ResetTimer()
			for b.Loop() {
				IsKingAttacked(White, bbs)
				IsKingAttacked(Black, bbs)
			}
		})
	}
}
