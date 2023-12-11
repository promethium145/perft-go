package movegen

import (
	bd "perft/board"
)

type pawnMoveType int8

const (
	NWCapture pawnMoveType = 9
	NECapture              = 7
	OneStep                = 8
	TwoStep                = 16
)

func genPawnMoves(b *bd.Board, ma *bd.MoveArray) {
	s := b.SideToMove()
	opps := s.Opposite()
	oppbb := b.SBitBoard(opps)
	pawnCapturable := oppbb | epBB(s, b)
	pawnbb := b.PBitBoard(bd.Pawn.Piece(s))
	if nwCaptured := sideRelPushNW(s, pawnbb) & pawnCapturable; nwCaptured != 0 {
		addPawnMoves(nwCaptured, NWCapture, s, ma)
	}
	if neCaptured := sideRelPushNE(s, pawnbb) & pawnCapturable; neCaptured != 0 {
		addPawnMoves(neCaptured, NECapture, s, ma)
	}
	emptybb := ^b.BitBoard()
	oneStep := sideRelPushFront(s, pawnbb) & emptybb
	twoStep := sideRelPushFront(s, oneStep) & emptybb & sideRelMaskRow(s, 3)
	addPawnMoves(oneStep, OneStep, s, ma)
	addPawnMoves(twoStep, TwoStep, s, ma)
}

func addPawnMoves(pb bd.BitBoard, pmt pawnMoveType, side bd.Side, ma *bd.MoveArray) {
	add := int(pmt)
	if side.IsW() {
		add = -int(pmt)
	}
	for pb != 0 {
		lsbIdx := uint(pb.Lsb1())
		from := bd.Square(int(lsbIdx) + add)
		if lsbIdx <= 7 || lsbIdx >= 56 {
			addPawnPromotions(from, bd.Square(lsbIdx), side, ma)
		} else {
			ma.Add(bd.NewMove(from, bd.Square(lsbIdx), bd.PieceNone))
		}
		pb ^= (bd.BitBoard(1) << lsbIdx)
	}
}

func addPawnPromotions(from, to bd.Square, side bd.Side, ma *bd.MoveArray) {
	addMove := func(p bd.PieceType) {
		ma.Add(bd.NewMove(from, to, p.Piece(side)))
	}
	addMove(bd.Bishop)
	addMove(bd.Knight)
	addMove(bd.Queen)
	addMove(bd.Rook)
}

func epBB(s bd.Side, b *bd.Board) bd.BitBoard {
	epSq := b.EpTarget()
	if !epSq.IsValid() {
		return bd.BitBoard(0)
	}
	return (bd.BitBoard(1) << uint(epSq)) & sideRelMaskRow(s, 5)
}

func sideRelMaskRow(s bd.Side, row uint) bd.BitBoard {
	if row >= 0 && row <= 7 {
		if s.IsW() {
			return bd.BitBoard(0xFF) << (row * 8)
		}
		return bd.BitBoard(0xFF00000000000000) >> (row * 8)
	}
	panic("invalid row")
}

func sideRelMaskCol(s bd.Side, col uint) bd.BitBoard {
	if col >= 0 && col <= 7 {
		if s.IsW() {
			return bd.BitBoard(0x8080808080808080) >> col
		}
		return bd.BitBoard(0x0101010101010101) << col
	}
	panic("invalid col")
}

func sideRelLeftShift(s bd.Side, bb bd.BitBoard, shift uint) bd.BitBoard {
	if s.IsW() {
		return bb << shift
	}
	return bb >> shift
}

func sideRelRightShift(s bd.Side, bb bd.BitBoard, shift uint) bd.BitBoard {
	return sideRelLeftShift(s.Opposite(), bb, shift)
}

func sideRelPushFront(s bd.Side, bb bd.BitBoard) bd.BitBoard {
	return sideRelLeftShift(s, bb, 8)
}

func sideRelPushNW(s bd.Side, bb bd.BitBoard) bd.BitBoard {
	return sideRelLeftShift(s, bb & ^sideRelMaskCol(s, 0), 9)
}

func sideRelPushNE(s bd.Side, bb bd.BitBoard) bd.BitBoard {
	return sideRelLeftShift(s, bb & ^sideRelMaskCol(s, 7), 7)
}

func pawnCaptures(s bd.Side, pawnbb, pawnCapturable bd.BitBoard) bd.BitBoard {
	return pawnCapturable & (sideRelPushNW(s, pawnbb) | sideRelPushNE(s, pawnbb))
}
