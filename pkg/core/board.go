package core

const (
	Space       = 2
	BoardWidth  = 10
	BoardHeight = 20 + Space
)

type Board struct {
	Cells []PieceType
}

func NewBoard() *Board {
	return &Board{
		Cells: make([]PieceType, BoardWidth*BoardHeight),
	}
}

func (b *Board) IsInside(p Point) bool {
	return p.X >= 0 && p.X < BoardWidth && p.Y >= 0 && p.Y < BoardHeight
}

func (b *Board) Get(p Point) PieceType {
	if !b.IsInside(p) {
		return PieceNone
	}
	return b.Cells[b.getIndex(p)]
}

func (b *Board) Set(p Point, t PieceType) {
	if !b.IsInside(p) {
		return
	}
	b.Cells[b.getIndex(p)] = t
}

func (b *Board) getIndex(p Point) int {
	return p.Y*BoardWidth + p.X
}

func (b *Board) Clear() {
	for i := range b.Cells {
		b.Cells[i] = PieceNone
	}
}

func (b *Board) HasCollision(p Piece) bool {
	minos := GetRotatedMinos(p.Type, p.Rotation)

	for _, mino := range minos {

		absolute := p.Position.Add(mino)
		if absolute.X < 0 || absolute.X >= BoardWidth || absolute.Y >= BoardHeight {
			return true
		}

		if absolute.Y >= 0 {
			if b.Get(absolute) != PieceNone {
				return true
			}
		}
	}

	return false
}

func GetRotatedMinos(t PieceType, rotation int) []Point {
	minos := GetMinos(t)
	result := make([]Point, len(minos))

	rotation = rotation % 4
	if rotation < 0 {
		rotation += 4
	}

	for i, p := range minos {
		for r := 0; r < rotation; r++ {
			p = p.RotateCW()
		}
		result[i] = p
	}
	return result
}
