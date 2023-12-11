package board

import "testing"

var initFen = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq -"

func TestNewBoard(t *testing.T) {
	_, err := NewBoard(initFen)
	if err != nil {
		t.Error("Expected board initialization")
	}
}

func TestFen(t *testing.T) {
	board, err := NewBoard(initFen)
	if err != nil {
		t.Error("Expected board initialization")
	}
	if board.Fen() != initFen {
		t.Error("Parsed FEN different from initFen")
	}
}

func TestFen2(t *testing.T) {
	expFen := "r2qr1k1/ppbn1pp1/4bn1p/PN1pp3/1P2P3/3P1N2/2Q1BPPP/R1B2RK1 b - -"
	board, err := NewBoard(expFen)
	if err != nil {
		t.Error("Expected board initialization")
	}
	if board.Fen() != expFen {
		t.Error("Parsed FEN different from expected FEN")
	}
}
