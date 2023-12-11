package movegen

import (
	bd "perft/board"
)

const (
	FILE_A = 0
	FILE_B = 1
	FILE_C = 2
	FILE_D = 3
	FILE_E = 4
	FILE_F = 5
	FILE_G = 6
	FILE_H = 7
)

func GenerateMoves(b *bd.Board) *bd.MoveArray {
	pma := &bd.MoveArray{}
	genPieceMoves(b, bd.Bishop, pma)
	genPieceMoves(b, bd.King, pma)
	genPieceMoves(b, bd.Knight, pma)
	genPieceMoves(b, bd.Pawn, pma)
	genPieceMoves(b, bd.Queen, pma)
	genPieceMoves(b, bd.Rook, pma)

	s := b.SideToMove()
	opps := s.Opposite()
	kingP := bd.King.Piece(s)
	kingSq := b.PBitBoard(kingP).Lsb1()
	if kingSq <= 0 {
		panic("bad board")
	}
	ma := &bd.MoveArray{}
	// Attacks on current playing side.
	oppAttackMap := genAttackBitBoard(opps, b)
	// If king is under check, work on evading.
	if oppAttackMap&b.PBitBoard(kingP) != 0 {
		for i := 0; i < pma.Size(); i++ {
			move := pma.Get(i)
			b.MakeMove(move)
			newOppAttackMap := genAttackBitBoard(opps, b)
			// Include move if legal.
			if newOppAttackMap&b.PBitBoard(kingP) == 0 {
				ma.Add(move)
			}
			b.UnmakeLastMove()
		}
	} else {
		if b.CanCastle(s, bd.King) || b.CanCastle(s, bd.Queen) {
			genCastlingMoves(s, b, oppAttackMap, ma)
		}

		// Check if there are any pinned pieces. This is not
		// an exact set of pins but always a superset.
		potentialPins := Attacks(b.PBitBoard(kingP), kingSq, bd.BQueen) & oppAttackMap & b.SBitBoard(s)
		for i := 0; i < pma.Size(); i++ {
			move := pma.Get(i)

			// Add all non-king, non-pinned piece moves.
			if move.From() != kingSq && (bd.BitBoard(1)<<uint(move.From())&potentialPins) == 0 {
				ma.Add(move)
				continue
			}

			// For king moves and pinned pieces, make sure the moves are valid.
			b.MakeMove(move)
			newOppAttackMap := genAttackBitBoard(opps, b)
			if newOppAttackMap&b.PBitBoard(kingP) == 0 {
				ma.Add(move)
			}
			b.UnmakeLastMove()
		}

	}

	return ma
}

func genCastlingMoves(s bd.Side, b *bd.Board, oppAttackMap bd.BitBoard, ma *bd.MoveArray) {
	opps := s.Opposite()
	row := bd.Row(0)
	if s.IsB() {
		row = bd.Row(7)
	}
	pawnbb := b.PBitBoard(bd.Pawn.Piece(opps))
	pawnCapEmptySquares := sideRelPushNW(opps, pawnbb) | sideRelPushNE(opps, pawnbb) & ^b.BitBoard()
	restrictedSquares := oppAttackMap | pawnCapEmptySquares

	isSqEmpty := func(r bd.Row, c bd.Col) bool {
		return b.PieceAt(bd.MakeSquare(r, c)) == bd.PieceNone
	}

	if b.CanCastle(s, bd.King) {
		if isSqEmpty(row, FILE_F) && isSqEmpty(row, FILE_G) &&
			(restrictedSquares&(bd.MakeSquare(row, FILE_F).BitBoard()|bd.MakeSquare(row, FILE_G).BitBoard())) == 0 {
			ma.Add(bd.NewMove(bd.MakeSquare(row, FILE_E), bd.MakeSquare(row, FILE_G), bd.PieceNone))
		}
	}
	if b.CanCastle(s, bd.Queen) {
		if isSqEmpty(row, FILE_C) && isSqEmpty(row, FILE_D) && isSqEmpty(row, FILE_B) &&
			(restrictedSquares&(bd.MakeSquare(row, FILE_C).BitBoard()|bd.MakeSquare(row, FILE_D).BitBoard())) == 0 {
			ma.Add(bd.NewMove(bd.MakeSquare(row, FILE_E), bd.MakeSquare(row, FILE_C), bd.PieceNone))
		}
	}
}

func genAttackBitBoard(s bd.Side, b *bd.Board) bd.BitBoard {
	opps := s.Opposite()
	oppbb := b.SBitBoard(opps)
	pawnbb := b.PBitBoard(bd.Pawn.Piece(s))
	emptybb := ^b.BitBoard()
	pc := pawnCaptures(s, pawnbb, oppbb|epBB(s, b))
	pone := sideRelPushFront(s, pawnbb) & emptybb
	ptwo := sideRelPushFront(s, pone) & emptybb & sideRelMaskRow(s, 3)

	aa := func(pt bd.PieceType) bd.BitBoard {
		return allAttacks(pt.Piece(s), b)
	}

	return pc | pone | ptwo | aa(bd.Queen) | aa(bd.Rook) | aa(bd.Bishop) | aa(bd.Knight) | aa(bd.King)
}

func allAttacks(p bd.Piece, b *bd.Board) bd.BitBoard {
	occbb := b.BitBoard()
	piecebb := b.PBitBoard(p)
	attackMap := bd.BitBoard(0)
	for piecebb != 0 {
		sq := piecebb.Lsb1()
		attackMap |= Attacks(occbb, sq, p)
		piecebb ^= (bd.BitBoard(1) << uint(sq))
	}
	return attackMap
}

func genPieceMoves(b *bd.Board, pt bd.PieceType, ma *bd.MoveArray) {
	s := b.SideToMove()
	p := pt.Piece(s)
	if p.IsPawn() {
		genPawnMoves(b, ma)
		return
	}
	occb := b.BitBoard()
	selfb := b.SBitBoard(s)
	pieceb := b.PBitBoard(p)
	for pieceb != 0 {
		sq := pieceb.Lsb1()
		attackMap := Attacks(occb, sq, p) & ^selfb
		if attackMap != 0 {
			bbMoves(sq, attackMap, ma)
		}
		pieceb ^= (bd.BitBoard(1) << uint(sq))
	}
}

func bbMoves(sq bd.Square, bb bd.BitBoard, ma *bd.MoveArray) {
	for bb != 0 {
		toSq := bb.Lsb1()
		ma.Add(bd.NewMove(sq, toSq, bd.PieceNone))
		bb ^= (bd.BitBoard(1) << uint(toSq))
	}
}
