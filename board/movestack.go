package board

type moveStackEntry struct {
	move Move

	// Any piece captured by this move.
	capturedPiece Piece

	// Bits (0, 1, 2, 3) represent (WKing, WQueen, BKing, BQueen) side castling availability
	// respectively.
	castle uint8

	// Enpassant target position only updated if last move was a pawn advanced by two squares
	// from its starting position.
	epSq Square

	// Zorbist key of the board position after this move is played.
	zKey uint64
}

type moveStack struct {
	entries  [1000]moveStackEntry
	topIndex int
}

func (s *moveStack) push() {
	s.topIndex++
}

func (s *moveStack) pop() {
	s.topIndex--
}

func (s *moveStack) top() *moveStackEntry {
	return &s.entries[s.topIndex]
}

func (s *moveStack) seek(fromTop int) *moveStackEntry {
	return &s.entries[s.topIndex-fromTop]
}
