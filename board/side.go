package board

type Side int8

const (
	SideNone Side = iota
	SideWhite
	SideBlack
	NumSides int = iota
)

func (side Side) Opposite() Side {
	switch side {
	case SideWhite:
		return SideBlack
	case SideBlack:
		return SideWhite
	default:
		return SideNone
	}
}

func (side Side) IsW() bool {
	return side == SideWhite
}

func (side Side) IsB() bool {
	return side == SideBlack
}

func (side Side) String() string {
	switch side {
	case SideWhite:
		return "w"
	case SideBlack:
		return "b"
	default:
		return ""
	}
}
