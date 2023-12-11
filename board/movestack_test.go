package board

import (
	"testing"
)

func TestMoveStack(t *testing.T) {
	ms := &moveStack{}
	mse1, mse2 := &moveStackEntry{
		move: NewMoveFromStr("e2e4"),
	}, &moveStackEntry{
		move: NewMoveFromStr("e6e5"),
	}
	*ms.top() = *mse1
	ms.push()
	*ms.top() = *mse2
	ms.push()
	if ms.top().move.IsValid() {
		t.Error("unexpected valid move")
	}
	ms.pop()
	if ms.top().move != mse2.move {
		t.Errorf("want %s, got %s", &mse2.move, &ms.top().move)
	}
	ms.pop()
	if ms.top().move != mse1.move {
		t.Errorf("want %s, got %s", &mse1.move, &ms.top().move)
	}
}
