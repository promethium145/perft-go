package board

import "fmt"
import "math/bits"

type BitBoard uint64

func (bb BitBoard) Lsb1() Square {
	return Square(bits.TrailingZeros64(uint64(bb)))
}

func (bb BitBoard) PopCount() int {
	return bits.OnesCount64(uint64(bb))
}

func (b BitBoard) String() string {
	s := fmt.Sprintf("%d\n", b)
	for i := 7; i >= 0; i-- {
		for j := 0; j < 8; j++ {
			if ((1 << uint(i*8+j)) & b) != 0 {
				s += "#"
			} else {
				s += "."
			}
		}
		s += "\n"
	}
	return s
}
