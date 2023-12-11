package board

import "fmt"

type MoveArray struct {
	moves [256]Move
	n     int
}

func (ma *MoveArray) Add(m Move) {
	ma.moves[ma.n] = m
	ma.n++
}

func (ma *MoveArray) Contains(m Move) bool {
	for i := 0; i < ma.n; i++ {
		if ma.moves[i] == m {
			return true
		}
	}
	return false
}

func (ma *MoveArray) Size() int {
	return ma.n
}

func (ma *MoveArray) Get(i int) Move {
	return ma.moves[i]
}

func (ma *MoveArray) String() string {
	s := ""
	for i := 0; i < ma.n; i++ {
		s += fmt.Sprintf("%s ", ma.Get(i))
	}
	return s
}
