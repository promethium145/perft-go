package board

import (
	"testing"
)

func TestPieceToPieceType(t *testing.T) {
	if PieceNone.Type() != PieceTypeNone {
		t.Error("bad")
	}
	if BKing.Type() != King {
		t.Error("bad")
	}
	if BPawn.Type() != Pawn {
		t.Error("bad")
	}
	if WKing.Type() != King {
		t.Error("bad")
	}
	if WPawn.Type() != Pawn {
		t.Error("bad")
	}
}

func TestPieceTypeToPiece(t *testing.T) {
	if King.Piece(SideWhite) != WKing {
		t.Error("bad")
	}
	if King.Piece(SideBlack) != BKing {
		t.Error("bad")
	}
}
