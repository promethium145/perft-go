package board

import (
	"math/rand"
	"time"
)

var zobrist [NumSquares][NumPieces]uint64
var ep [NumSquares]uint64
var turn uint64
var castling [16]uint64

func init() {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < NumSquares; i++ {
		for j := 0; j < NumPieces; j++ {
			zobrist[i][j] = rand.Uint64()
		}
	}
	for i := 0; i < NumSquares; i++ {
		ep[i] = rand.Uint64()
	}
	turn = rand.Uint64()
	for i := 0; i < 16; i++ {
		castling[i] = rand.Uint64()
	}
}

func zGet(piece Piece, sq Square) uint64 {
	return zobrist[sq][int(piece)]
}

func zEp(sq Square) uint64 {
	return ep[sq]
}

func zTurn() uint64 {
	return turn
}

func zCastling(castle uint8) uint64 {
	return castling[castle]
}
