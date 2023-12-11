package board

import "testing"

func TestNewMove(t *testing.T) {
	fromSq, toSq := Square(1), Square(17)
	m := NewMove(fromSq, toSq, PieceNone)
	if m.From() != fromSq {
		t.Errorf("want %d, got %d", fromSq, m.From())
	}
	if m.To() != toSq {
		t.Errorf("want %d, got %d", toSq, m.To())
	}
	if m.PromotedPiece().IsValid() {
		t.Error("unexpected promotion.")
	}
}

func TestNewMoveFromStr(t *testing.T) {
	m := NewMoveFromStr("e2e4")
	fromSq, toSq := Square(12), Square(28)
	if m.From() != fromSq {
		t.Errorf("want %d, got %d", fromSq, m.From())
	}
	if m.To() != toSq {
		t.Errorf("want %d, got %d", toSq, m.To())
	}
}
