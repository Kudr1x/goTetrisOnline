package core

const (
	PieceNone PieceType = iota
	PieceI
	PieceO
	PieceT
	PieceS
	PieceZ
	PieceJ
	PieceL
	PieceGarbage
)

type Point struct {
	X, Y int
}

func (p Point) RotateCW() Point {
	return Point{
		X: -p.Y,
		Y: p.X,
	}
}

func (p Point) RotateCCW() Point {
	return Point{
		X: p.Y,
		Y: -p.X,
	}
}

func (p Point) Add(other Point) Point {
	return Point{p.X + other.X, p.Y + other.Y}
}

type PieceType uint8

type Piece struct {
	Type     PieceType
	Position Point
	Rotation int
}

func GetMinos(t PieceType) []Point {
	switch t {
	case PieceI:
		return []Point{{-1, 0}, {0, 0}, {1, 0}, {2, 0}}
	case PieceJ:
		return []Point{{-1, 1}, {-1, 0}, {0, 0}, {1, 0}}
	case PieceL:
		return []Point{{1, 1}, {-1, 0}, {0, 0}, {1, 0}}
	case PieceO:
		return []Point{{0, 1}, {1, 1}, {0, 0}, {1, 0}}
	case PieceS:
		return []Point{{0, 1}, {1, 1}, {-1, 0}, {0, 0}}
	case PieceT:
		return []Point{{0, 1}, {-1, 0}, {0, 0}, {1, 0}}
	case PieceZ:
		return []Point{{-1, 1}, {0, 1}, {0, 0}, {1, 0}}
	default:
		return []Point{}
	}
}
