package core

const (
	RotateCW  = 1
	RotateCCW = -1
)

type KickData [5]Point

func kickIndex(fromRotation, direction int) int {
	if direction > 0 {
		return fromRotation * 2
	}
	return fromRotation*2 + 1
}

var kicksJLSTZ = [8]KickData{
	// 0→CW (0→R)
	0: {{0, 0}, {-1, 0}, {-1, -1}, {0, 2}, {-1, 2}},
	// 0→CCW (0→L)
	1: {{0, 0}, {1, 0}, {1, -1}, {0, 2}, {1, 2}},
	// R→CW (R→2)
	2: {{0, 0}, {1, 0}, {1, 1}, {0, -2}, {1, -2}},
	// R→CCW (R→0)
	3: {{0, 0}, {1, 0}, {1, 1}, {0, -2}, {1, -2}},
	// 2→CW (2→L)
	4: {{0, 0}, {1, 0}, {1, -1}, {0, 2}, {1, 2}},
	// 2→CCW (2→R)
	5: {{0, 0}, {-1, 0}, {-1, -1}, {0, 2}, {-1, 2}},
	// L→CW (L→0)
	6: {{0, 0}, {-1, 0}, {-1, 1}, {0, -2}, {-1, -2}},
	// L→CCW (L→2)
	7: {{0, 0}, {-1, 0}, {-1, 1}, {0, -2}, {-1, -2}},
}

var kicksI = [8]KickData{
	0: {{0, 0}, {-2, 0}, {1, 0}, {-2, 1}, {1, -2}},
	1: {{0, 0}, {-1, 0}, {2, 0}, {-1, -2}, {2, 1}},
	2: {{0, 0}, {-1, 0}, {2, 0}, {-1, -2}, {2, 1}},
	3: {{0, 0}, {2, 0}, {-1, 0}, {2, -1}, {-1, 2}},
	4: {{0, 0}, {2, 0}, {-1, 0}, {2, -1}, {-1, 2}},
	5: {{0, 0}, {1, 0}, {-2, 0}, {1, 2}, {-2, -1}},
	6: {{0, 0}, {1, 0}, {-2, 0}, {1, 2}, {-2, -1}},
	7: {{0, 0}, {-2, 0}, {1, 0}, {-2, 1}, {1, -2}},
}

func getKickTable(t PieceType) *[8]KickData {
	switch t {
	case PieceI:
		return &kicksI
	case PieceO:
		return nil
	default:
		return &kicksJLSTZ
	}
}

func TryRotate(b *Board, p Piece, direction int) (Piece, bool) {
	table := getKickTable(p.Type)
	if table == nil {
		rotated := p
		rotated.Rotation = (p.Rotation + direction%4 + 4) % 4
		return rotated, true
	}

	newRotation := (p.Rotation + direction%4 + 4) % 4

	idx := kickIndex(p.Rotation, direction)
	kicks := table[idx]

	for _, kick := range kicks {
		candidate := p
		candidate.Rotation = newRotation
		candidate.Position = p.Position.Add(kick)

		if !b.HasCollision(candidate) {
			return candidate, true
		}
	}
	return p, false
}
