package movegen

import (
	bd "perft/board"
	"testing"
)

func TestEdgeDist(t *testing.T) {
	sq := bd.StrToSquare("e3")
	data := []struct {
		dir      direction
		wantDist int
	}{
		{east, 3},
		{north, 5},
		{northEast, 3},
		{northWest, 4},
		{south, 2},
		{southEast, 2},
		{southWest, 2},
		{west, 4},
	}
	for _, d := range data {
		gotDist := d.dir.edgeDist(sq)
		if gotDist != d.wantDist {
			t.Errorf("bad edgeDist, got: %v, want: %v", gotDist, d.wantDist)
		}
	}
}

func TestNextSq(t *testing.T) {
	data := []struct {
		sq      string
		dir     direction
		wantNSq string
	}{
		{"a1", east, "b1"},
		{"a1", north, "a2"},
		{"a1", northEast, "b2"},
		{"a1", northWest, ""},
		{"a1", south, ""},
		{"a1", southEast, ""},
		{"a1", southWest, ""},
		{"a1", west, ""},
		{"e3", east, "f3"},
		{"e3", north, "e4"},
		{"e3", northEast, "f4"},
		{"e3", northWest, "d4"},
		{"e3", south, "e2"},
		{"e3", southEast, "f2"},
		{"e3", southWest, "d2"},
		{"e3", west, "d3"},
		{"h8", east, ""},
		{"h8", north, ""},
		{"h8", northEast, ""},
		{"h8", northWest, ""},
		{"h8", south, "h7"},
		{"h8", southEast, ""},
		{"h8", southWest, "g7"},
		{"h8", west, "g8"},
	}
	for _, d := range data {
		sq := bd.StrToSquare(d.sq)
		gotNSq := d.dir.nextSq(sq).String()
		if gotNSq != d.wantNSq {
			t.Errorf("bad nextSq, got: %v, want: %v", gotNSq, d.wantNSq)
		}
	}
}

func TestMaskBits(t *testing.T) {
	sq := bd.StrToSquare("e3")
	data := []struct {
		dir    direction
		wantbb bd.BitBoard
	}{
		{east, 6291456},
		{north, 4521260801327104},
		{northEast, 275414777856},
		{northWest, 2216337342464},
		{south, 4096},
		{southEast, 8192},
		{southWest, 2048},
		{west, 917504},
	}
	for _, d := range data {
		gotbb := d.dir.maskBits(sq)
		if gotbb != d.wantbb {
			t.Errorf("bad maskBits, got: %v, want: %v", gotbb, d.wantbb)
		}
	}
}

func TestOccs(t *testing.T) {
	sq := bd.StrToSquare("e3")
	wantOccs := map[bd.BitBoard]bool{
		0:                true,
		268435456:        true,
		68719476736:      true,
		68987912192:      true,
		17592186044416:   true,
		17592454479872:   true,
		17660905521152:   true,
		17661173956608:   true,
		4503599627370496: true,
		4503599895805952: true,
		4503668346847232: true,
		4503668615282688: true,
		4521191813414912: true,
		4521192081850368: true,
		4521260532891648: true,
		4521260801327104: true,
	}
	gotOccs := north.occs(sq)
	if len(gotOccs) != len(wantOccs) {
		t.Errorf("bad occs: got: %d, want: %d", len(gotOccs), len(wantOccs))
	}
	for _, o := range gotOccs {
		if _, ok := wantOccs[o]; !ok {
			t.Errorf("bad occs: got %d but did not want", o)
		}
	}
	wantOccs = map[bd.BitBoard]bool{
		0:            true,
		536870912:    true,
		274877906944: true,
		275414777856: true,
	}
	gotOccs = northEast.occs(sq)
	if len(gotOccs) != len(wantOccs) {
		t.Errorf("bad occs: got: %d, want: %d", len(gotOccs), len(wantOccs))
	}
	for _, o := range gotOccs {
		if _, ok := wantOccs[o]; !ok {
			t.Errorf("bad occs: got %v but did not want", o)
		}
	}
}

func TestAttack(t *testing.T) {
	testBoard := "2b5/P1q5/3P4/6r1/2p2p2/b1BR4/5p1R/1K3k2 w - -"
	board, err := bd.NewBoard(testBoard)
	if err != nil {
		t.Errorf("bad board, error: %v", err)
	}
	occ := board.BitBoard()
	data := []struct {
		sq     string
		dir    direction
		wantbb bd.BitBoard
	}{
		{"g7", southWest, 35253226045440},
		{"a3", northEast, 8813306445824},
		{"e3", north, 1157442765408174080},
		{"e3", south, 4112},
	}
	for _, d := range data {
		sq := bd.StrToSquare(d.sq)
		gotbb := d.dir.attack(sq, occ)
		if gotbb != d.wantbb {
			t.Errorf("bad attack: got %v but want %v", gotbb, d.wantbb)
		}
	}
}
