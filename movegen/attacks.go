package movegen

import (
	bd "perft/board"
)

func Attacks(bb bd.BitBoard, sq bd.Square, p bd.Piece) bd.BitBoard {
	switch p.Type() {
	case bd.Bishop:
		return bishopAttacks(bb, sq)
	case bd.King:
		return kingAttacks[sq]
	case bd.Knight:
		return knightAttacks[sq]
	case bd.Queen:
		return queenAttacks(bb, sq)
	case bd.Rook:
		return rookAttacks(bb, sq)
	}
	panic("unknown piece")
}

func rookAttacks(bb bd.BitBoard, sq bd.Square) bd.BitBoard {
	occ := bb & rookTables.masks[sq]
	attackIndx := ((occ * rookTables.magics[sq]) >> (64 - rookTables.shifts[sq])) + bd.BitBoard(rookTables.offsets[sq])
	return rookTables.attacks[attackIndx]
}

func bishopAttacks(bb bd.BitBoard, sq bd.Square) bd.BitBoard {
	occ := bb & bishopTables.masks[sq]
	attackIndx := ((occ * bishopTables.magics[sq]) >> (64 - bishopTables.shifts[sq])) + bd.BitBoard(bishopTables.offsets[sq])
	return bishopTables.attacks[attackIndx]
}

func queenAttacks(bb bd.BitBoard, sq bd.Square) bd.BitBoard {
	return rookAttacks(bb, sq) | bishopAttacks(bb, sq)
}
