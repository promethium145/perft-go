package movegen

import (
	bd "perft/board"
)

type direction int

const (
	north direction = iota
	south
	east
	west
	northEast
	northWest
	southEast
	southWest
)

func (d direction) String() string {
	switch d {
	case north:
		return "north"
	case south:
		return "south"
	case east:
		return "east"
	case west:
		return "west"
	case northEast:
		return "north-east"
	case northWest:
		return "north-west"
	case southEast:
		return "south-east"
	case southWest:
		return "south-west"
	}
	panic("invalid state")
}

func (d direction) nextSq(sq bd.Square) bd.Square {
	row, col := sq.Row(), sq.Col()
	switch d {
	case north:
		row++
	case south:
		row--
	case east:
		col++
	case west:
		col--
	case northEast:
		row++
		col++
	case northWest:
		row++
		col--
	case southEast:
		row--
		col++
	case southWest:
		row--
		col--
	}
	if !row.IsValid() || !col.IsValid() {
		return bd.InvalidSq
	}
	return bd.MakeSquare(row, col)
}

func (d direction) edgeDist(sq bd.Square) int {
	row, col := int(sq.Row()), int(sq.Col())
	min := func(a int, b int) int {
		if a < b {
			return a
		}
		return b
	}
	switch d {
	case north:
		return 7 - row
	case south:
		return row
	case east:
		return 7 - col
	case west:
		return col
	case northEast:
		return min(7-row, 7-col)
	case northWest:
		return min(7-row, col)
	case southEast:
		return min(row, 7-col)
	case southWest:
		return min(row, col)
	}
	panic("Reached invalid state")
}

// Mask all bits in the given direction starting from sq, excluding sq
// and edge of the board.
func (d direction) maskBits(sq bd.Square) bd.BitBoard {
	var bb bd.BitBoard
	for nextSq := d.nextSq(sq); nextSq.IsValid() && d.nextSq(nextSq).IsValid(); nextSq = d.nextSq(nextSq) {
		bb |= (bd.BitBoard(1) << uint(nextSq))
	}
	return bb
}

// All possible occupancies from the given square along a direction excluding given sq and edge of the board.
func (d direction) occs(sq bd.Square) []bd.BitBoard {
	bbs := make([]bd.BitBoard, 0)

	// Number of squares in this direction excluding sq and edge of board.
	numSq := d.edgeDist(sq) - 1
	if numSq <= 0 {
		return bbs
	}

	// Number of possible piece occupancies in squares along this direction.
	numOcc := uint(1) << uint(numSq)

	// Create bitboard for each occupancy.
	for occ := uint(0); occ < numOcc; occ++ {
		var bb bd.BitBoard
		nextSq := sq
		for bitMask := uint(1); bitMask <= occ; bitMask <<= 1 {
			nextSq = d.nextSq(nextSq)
			if occ&bitMask != 0 {
				bb |= bd.BitBoard(1) << uint(nextSq)
			}
		}
		bbs = append(bbs, bb)
	}
	return bbs
}

// Attack bitboard along a direction from the given square given an occupancy bitboard.
func (d direction) attack(sq bd.Square, occ bd.BitBoard) bd.BitBoard {
	var attackBb bd.BitBoard
	for nextSq := d.nextSq(sq); nextSq.IsValid(); nextSq = d.nextSq(nextSq) {
		attackBb |= (bd.BitBoard(1) << uint(nextSq))
		if occ&(bd.BitBoard(1)<<uint(nextSq)) != 0 {
			break
		}
	}
	return attackBb
}
