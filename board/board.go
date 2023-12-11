package board

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type Board struct {
	array          [NumSquares]Piece
	bitboardSides  [NumSides]BitBoard
	bitboardPieces [NumPieces]BitBoard
	sideToMove     Side
	moveStack      moveStack
}

func NewBoard(fen string) (*Board, error) {
	var b Board

	fenParts := strings.Split(fen, " ")
	if len(fenParts) != 4 {
		return nil, fmt.Errorf("Invalid fen: %s", fen)
	}

	boardParts := strings.Split(fenParts[0], "/")
	if len(boardParts) != 8 {
		return nil, fmt.Errorf("Invalid fen: %s", fen)
	}

	row, col := Row(7), Col(0)
	for _, boardPart := range boardParts {
		for _, r := range boardPart {
			if unicode.IsNumber(r) {
				v, err := strconv.Atoi(string(r))
				if err != nil {
					return nil, err
				}
				col += Col(v)
			} else {
				b.array[MakeSquare(row, col)] = CharToPiece(r)
				col++
			}
		}
		col = 0
		row--
	}

	switch fenParts[1] {
	case "w":
		b.sideToMove = SideWhite
	case "b":
		b.sideToMove = SideBlack
	default:
		return nil, fmt.Errorf("Expected w or b. Found: %s", fenParts[1])
	}

	castle := &b.moveStack.top().castle
	for _, r := range fenParts[2] {
		switch r {
		case 'K':
			*castle |= 0x1
		case 'Q':
			*castle |= 0x2
		case 'k':
			*castle |= 0x4
		case 'q':
			*castle |= 0x8
		}
	}

	if epStr := fenParts[3]; epStr == "-" {
		b.moveStack.top().epSq = InvalidSq
	} else {
		b.moveStack.top().epSq = StrToSquare(epStr)
	}

	for sq := Square(0); sq < NumSquares; sq++ {
		piece := b.array[sq]
		if piece == PieceNone {
			continue
		}
		bb := BitBoard(1 << uint(sq))
		b.bitboardPieces[piece] |= bb
		b.bitboardSides[piece.Side()] |= bb
	}

	return &b, nil
}

func (b *Board) SideToMove() Side {
	return b.sideToMove
}

func (b *Board) MakeMove(move Move) {
	b.moveStack.push()

	from, to := move.From(), move.To()
	fromRow, fromCol := from.Row(), from.Col()
	toRow, toCol := to.Row(), to.Col()
	srcPiece, dstPiece := b.array[from], b.array[to]

	top, prev := b.moveStack.top(), b.moveStack.seek(1)
	*top = moveStackEntry{
		move:          move,
		capturedPiece: dstPiece,
		epSq:          prev.epSq,
		zKey:          prev.zKey,
		castle:        prev.castle}

	// If previous move was a 2 space pawn move, update zobrist key to
	// distinguish between two boards with same positions but different
	// enpassant capture opportunities.
	if top.epSq.IsValid() {
		top.zKey ^= zEp(top.epSq)
	}

	// Remove pieces at source and destination (if any).
	b.removePiece(from)
	if dstPiece != PieceNone {
		b.removePiece(to)
	}

	// Handle pawn.
	if srcPiece.IsPawn() {
		top.epSq = InvalidSq

		if toRow.Distance(fromRow) == 2 { // 2 space pawn move
			top.epSq = MakeSquare((toRow+fromRow)>>1, toCol)
			top.zKey ^= zEp(top.epSq)
		} else if fromCol != toCol && !dstPiece.IsValid() { // enpassant capture
			b.removePiece(MakeSquare(fromRow, toCol))
		}

		if piece := move.PromotedPiece(); piece.IsValid() {
			b.placePiece(to, piece)
		} else {
			b.placePiece(to, srcPiece)
		}

		b.flipSideToMove()
		return
	}

	// Handle castling.
	castleKing, castleQueen := uint8(0), uint8(0)
	switch b.sideToMove {
	case SideBlack:
		castleKing, castleQueen = top.castle&0x4, top.castle&0x8
	case SideWhite:
		castleKing, castleQueen = top.castle&0x1, top.castle&0x2
	}

	castleBits := castleKing | castleQueen
	if castleBits != 0 {
		// Update castling bits if moving piece is King or Rook.
		switch srcPiece {
		case BKing, WKing:
			// When the king is moved by more than 1 space, it is for castling so
			// move the rook to appropriate square.
			if toCol.Distance(fromCol) > 1 {
				rookCol := Col(0)
				if toCol > fromCol {
					rookCol = Col(7)
				}
				rookSq := MakeSquare(fromRow, rookCol)
				rook := b.array[rookSq]
				b.removePiece(rookSq)
				b.placePiece(MakeSquare(fromRow, (toCol+fromCol)>>1), rook)
			}
			top.castle &= ^castleBits
			top.zKey ^= zCastling(top.castle)

		case BRook, WRook:
			// Depending on which rook moved, mask the relevant castle bit.
			if fromCol == 0 && castleQueen != 0 {
				top.castle &= ^castleQueen
				top.zKey ^= zCastling(top.castle)
			} else if fromCol == 0 && castleKing != 0 {
				top.castle &= ^castleKing
				top.zKey ^= zCastling(top.castle)
			}
		}
	}

	// If destination piece is rook on it's home index, update opponent castling rights.
	if dstPiece.IsRook() {
		switch to {
		case 0:
			top.castle &= ^uint8(0x2)
		case 7:
			top.castle &= ^uint8(0x1)
		case 56:
			top.castle &= ^uint8(0x8)
		case 63:
			top.castle &= ^uint8(0x4)
		}
	}

	top.epSq = InvalidSq
	b.placePiece(to, srcPiece)
	b.flipSideToMove()
}

func (b *Board) UnmakeLastMove() bool {
	if b.moveStack.topIndex <= 0 {
		return false
	}
	b.flipSideToMove()

	top := b.moveStack.top()
	move := top.move
	from, to := move.From(), move.To()
	fromRow, fromCol, toCol := from.Row(), from.Col(), to.Col()
	dstPiece := b.array[to]

	if piece := move.PromotedPiece(); piece.IsValid() {
		b.removePiece(to)
		if b.sideToMove == SideBlack {
			b.placePiece(from, BPawn)
		} else {
			b.placePiece(from, WPawn)
		}
		b.moveStack.pop()
		return true
	}

	switch dstPiece {
	case BPawn:
		if toCol != fromCol && !top.capturedPiece.IsValid() {
			b.placePiece(MakeSquare(fromRow, toCol), WPawn)
		}
	case WPawn:
		if toCol != fromCol && !top.capturedPiece.IsValid() {
			b.placePiece(MakeSquare(fromRow, toCol), BPawn)
		}
	case WKing, BKing:
		if toCol.Distance(fromCol) == 2 {
			rookOldSq := Square(0)
			if toCol > fromCol {
				rookOldSq = 7
			}
			rookCurSq := MakeSquare(fromRow, (toCol+fromCol)>>1)
			rook := b.array[rookCurSq]
			b.removePiece(rookCurSq)
			b.placePiece(rookOldSq, rook)
		}
	}

	b.removePiece(to)
	b.placePiece(from, dstPiece)
	if top.capturedPiece.IsValid() {
		b.placePiece(to, top.capturedPiece)
	}

	b.moveStack.pop()

	return true
}

func (b *Board) BitBoard() BitBoard {
	return b.bitboardSides[SideWhite] | b.bitboardSides[SideBlack]
}

func (b *Board) SBitBoard(s Side) BitBoard {
	return b.bitboardSides[s]
}

func (b *Board) EpTarget() Square {
	return b.moveStack.top().epSq
}

func (b *Board) PBitBoard(p Piece) BitBoard {
	return b.bitboardPieces[p]
}

func (b *Board) CanCastle(s Side, pt PieceType) bool {
	if s == SideNone {
		panic("invalid side")
	}
	shift := 0
	if s.IsB() {
		shift += 2
	}
	if pt == Queen {
		shift += 1
	}
	return b.moveStack.top().castle&(1<<uint(shift)) != 0
}

func (b *Board) PieceAt(sq Square) Piece {
	return b.array[sq]
}

func (b *Board) Fen() string {
	var fen string

	for row := Row(7); row >= 0; row-- {
		emptyCount := 0
		for col := Col(0); col <= 7; col++ {
			piece := b.array[MakeSquare(row, col)]
			if !piece.IsValid() {
				emptyCount++
				continue
			}
			if emptyCount != 0 {
				fen += fmt.Sprintf("%d", emptyCount)
				emptyCount = 0
			}
			fen += piece.String()
		}
		if emptyCount != 0 {
			fen += fmt.Sprintf("%d", emptyCount)
		}
		if row != 0 {
			fen += "/"
		}
	}

	fen += " " + b.sideToMove.String() + " "

	castle := b.moveStack.top().castle
	if castle != 0 {
		if castle&0x1 != 0 {
			fen += "K"
		}
		if castle&0x2 != 0 {
			fen += "Q"
		}
		if castle&0x4 != 0 {
			fen += "k"
		}
		if castle&0x8 != 0 {
			fen += "q"
		}
	} else {
		fen += "-"
	}
	fen += " "

	epSq := b.moveStack.top().epSq
	if epSq.IsValid() {
		fen += epSq.String()
	} else {
		fen += "-"
	}

	return fen
}

func (b *Board) flipSideToMove() {
	b.sideToMove = b.sideToMove.Opposite()
	b.moveStack.top().zKey ^= zTurn()
}

func (b *Board) removePiece(sq Square) {
	piece := b.array[sq]
	b.array[sq] = PieceNone
	mask := ^BitBoard(uint64(1) << uint(sq))
	b.bitboardSides[piece.Side()] &= mask
	b.bitboardPieces[piece] &= mask
	b.moveStack.top().zKey ^= zGet(piece, sq)
}

func (b *Board) placePiece(sq Square, piece Piece) {
	b.array[sq] = piece
	mask := BitBoard(uint64(1) << uint(sq))
	b.bitboardSides[piece.Side()] |= mask
	b.bitboardPieces[piece] |= mask
	b.moveStack.top().zKey ^= zGet(piece, sq)
}

func (b *Board) String() string {
	var s string
	for i := Row(7); i >= 0; i-- {
		for j := Col(0); j < 8; j++ {
			s += fmt.Sprintf("%v ", b.array[MakeSquare(i, j)])
		}
		s += "\n"
	}

	s += "Next side to move "
	switch b.sideToMove {
	case SideWhite:
		s += "white"
	case SideBlack:
		s += "black"
	default:
		s += "unknown!"
	}
	s += "\n"

	return s
}
