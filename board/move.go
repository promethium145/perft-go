package board

import "fmt"

type Move uint16

func NewMove(fromSq, toSq Square, promotedPiece Piece) Move {
	return Move(fromSq)<<10 | Move(toSq)<<4 | Move(promotedPiece)
}

func NewMoveFromStr(str string) Move {
	var piece Piece
	if len(str) == 5 {
		piece = CharToPiece(rune(str[4]))
	} else {
		piece = PieceNone
	}
	return NewMove(StrToSquare(str[0:2]), StrToSquare(str[2:4]), piece)
}

func (move Move) From() Square {
	return Square(move >> 10)
}

func (move Move) To() Square {
	return Square((move & 0x3FF) >> 4)
}

func (move Move) PromotedPiece() Piece {
	return Piece(move & 0xF)
}

func (move Move) IsValid() bool {
	return move != 0
}

func (move Move) IsPromotion() bool {
	return move.PromotedPiece().IsValid()
}

func (move Move) String() string {
	return fmt.Sprintf("%s%s", move.From(), move.To())
}
