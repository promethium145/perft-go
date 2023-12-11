package board

import (
	"testing"
)

func TestSquareStr(t *testing.T) {
	sq := MakeSquare(Row(1), Col(2))
	if sq.String() != "c2" {
		t.Errorf("want c2, got %s", string(sq))
	}
}
