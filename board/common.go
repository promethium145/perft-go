package board

import "fmt"

type Square int8
type Row int8
type Col int8

func (r Row) IsValid() bool {
	return r >= 0 && r <= 7
}

func (c Col) IsValid() bool {
	return c >= 0 && c <= 7
}

const NumSquares = 64
const InvalidSq = -1

func (s Square) Row() Row {
	return Row(s / 8)
}

func (s Square) Col() Col {
	return Col(s % 8)
}

func MakeSquare(row Row, col Col) Square {
	return Square(int(row*8) + int(col%8))
}

func (s Square) IsValid() bool {
	return s >= 0 && s <= 63
}

func (s Square) BitBoard() BitBoard {
	return BitBoard(1) << uint(s)
}

func StrToSquare(str string) Square {
	if len(str) != 2 {
		panic(fmt.Sprintf("Can't convert %s to Square", str))
	}
	col := Col(str[0] - 'a')
	row := Row(str[1] - '1')
	return MakeSquare(row, col)
}

func (s Square) String() string {
	if !s.IsValid() {
		return ""
	}
	row := s.Row()
	col := s.Col()
	return fmt.Sprintf("%c%d", col+'a', row+1)
}

func (r Row) Distance(r2 Row) int {
	if r < r2 {
		return int(r2 - r)
	}
	return int(r - r2)
}

func (c Col) Distance(c2 Col) int {
	if c < c2 {
		return int(c2 - c)
	}
	return int(c - c2)
}
