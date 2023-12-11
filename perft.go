package main

import (
	"fmt"
	bd "perft/board"
	mg "perft/movegen"
)

func perft(b *bd.Board, depth int) int {
	ma := mg.GenerateMoves(b)
	if depth == 1 {
		return ma.Size()
	}
	nodes := 0
	for i := 0; i < ma.Size(); i++ {
		b.MakeMove(ma.Get(i))
		nodes += perft(b, depth-1)
		b.UnmakeLastMove()
	}
	return nodes
}

func main() {
	board, err := bd.NewBoard("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq -")
	if err != nil {
		panic(err)
	}
	nodes := perft(board, 6)
	fmt.Println(nodes)
}
