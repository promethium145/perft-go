package movegen

import (
	bd "perft/board"
	"testing"
)

func TestRookAttacks(t *testing.T) {
	var b bd.BitBoard = 0xF7F6F3748CA5B610
	sq := bd.StrToSquare("e5")
	var want bd.BitBoard = 17781434093568
	got := Attacks(b, sq, bd.BRook)
	if want != got {
		t.Errorf("bad attacks bitboard for %v, want: %v, got: %v", b, want, got)
	}
}

func TestBishopAttacks(t *testing.T) {
	var b bd.BitBoard = 0xF7F6F3748CA5B610
	sq := bd.StrToSquare("e5")
	var want bd.BitBoard = 1169881047269376
	got := Attacks(b, sq, bd.WBishop)
	if want != got {
		t.Errorf("bad attacks bitboard for %v, want: %v, got: %v", b, want, got)
	}
}

func TestQueenAttacks(t *testing.T) {
	var b bd.BitBoard = 0xF7F6F3748CA5B610
	sq := bd.StrToSquare("e5")
	var want bd.BitBoard = 0x4382c38509000
	got := Attacks(b, sq, bd.BQueen)
	if want != got {
		t.Errorf("bad attacks bitboard for %v, want: %v, got: %v", b, want, got)
	}
	got = Attacks(b, sq, bd.BBishop) | Attacks(b, sq, bd.BRook)
	if want != got {
		t.Errorf("bad attacks bitboard for %v, want: %v, got: %v", b, want, got)
	}
}

func TestKingAttacks(t *testing.T) {
	var b bd.BitBoard = 0xF7F6F3748CA5B610
	sq := bd.StrToSquare("e5")
	var want bd.BitBoard = 61745389371392
	got := Attacks(b, sq, bd.WKing)
	if want != got {
		t.Errorf("bad attacks bitboard for %v, want: %v, got: %v", b, want, got)
	}
}

func TestKnightAttacks(t *testing.T) {
	var b bd.BitBoard = 0xF7F6F3748CA5B610
	sq := bd.StrToSquare("e5")
	var want bd.BitBoard = 11333767002587136
	got := Attacks(b, sq, bd.BKnight)
	if want != got {
		t.Errorf("bad attacks bitboard for %v, want: %v, got: %v", b, want, got)
	}
}
