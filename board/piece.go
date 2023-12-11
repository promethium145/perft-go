package board

type PieceType int8
type Piece int8

const (
	PieceTypeNone PieceType = iota
	King
	Queen
	Rook
	Bishop
	Knight
	Pawn
)

const (
	PieceNone Piece = iota
	BKing
	BQueen
	BRook
	BBishop
	BKnight
	BPawn
	WKing
	WQueen
	WRook
	WBishop
	WKnight
	WPawn
	NumPieces int = iota
)

func (piece Piece) Side() Side {
	switch {
	case piece == PieceNone:
		return SideNone
	case piece <= BPawn:
		return SideBlack
	default:
		return SideWhite
	}
}

var chars = []rune{'.', 'k', 'q', 'r', 'b', 'n', 'p', 'K', 'Q', 'R', 'B', 'N', 'P'}

func (piece Piece) AsRune() rune {
	return chars[piece]
}

func (piece Piece) String() string {
	return string(piece.AsRune())
}

var charToPieceMap = map[rune]Piece{
	'.': PieceNone,
	'k': BKing,
	'q': BQueen,
	'r': BRook,
	'b': BBishop,
	'n': BKnight,
	'p': BPawn,
	'K': WKing,
	'Q': WQueen,
	'R': WRook,
	'B': WBishop,
	'N': WKnight,
	'P': WPawn,
}

func CharToPiece(ch rune) Piece {
	return charToPieceMap[ch]
}

func (piece Piece) IsValid() bool {
	return piece != PieceNone
}

func (piece Piece) IsKing() bool {
	return piece == BKing || piece == WKing
}

func (piece Piece) IsQueen() bool {
	return piece == BQueen || piece == WQueen
}

func (piece Piece) IsRook() bool {
	return piece == BRook || piece == WRook
}

func (piece Piece) IsBishop() bool {
	return piece == BBishop || piece == WBishop
}

func (piece Piece) IsKnight() bool {
	return piece == BKnight || piece == WKnight
}

func (piece Piece) IsPawn() bool {
	return piece == BPawn || piece == WPawn
}

func (piece Piece) Type() PieceType {
	if piece <= BPawn {
		return PieceType(piece)
	}
	return PieceType(piece - 6)
}

func (ptype PieceType) Piece(s Side) Piece {
	if s == SideNone {
		return PieceNone
	} else if s.IsB() {
		return Piece(ptype)
	} else if s.IsW() {
		return Piece(ptype + 6)
	}
	panic("bad state")
}
