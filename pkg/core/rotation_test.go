package core

import "testing"

func buildBoard(filledRows int) *Board {
	b := NewBoard()
	for row := 0; row < filledRows; row++ {
		y := BoardHeight - 1 - row
		for x := 0; x < BoardWidth; x++ {
			b.Set(Point{X: x, Y: y}, PieceT)
		}
	}
	return b
}

func TestTryRotate_BasicCW(t *testing.T) {
	b := NewBoard()
	p := Piece{Type: PieceT, Position: Point{X: 5, Y: 5}, Rotation: 0}

	got, ok := TryRotate(b, p, RotateCW)

	if !ok {
		t.Fatal("Wait: fail on empty board")
	}
	if got.Rotation != 1 {
		t.Errorf("rotation: got %d, want 1", got.Rotation)
	}

	if got.Position != p.Position {
		t.Errorf("position change: %v: %v", p.Position, got.Position)
	}
}

func TestTryRotateWallKickLeft(t *testing.T) {
	b := NewBoard()
	p := Piece{Type: PieceI, Position: Point{X: 1, Y: 5}, Rotation: 1}

	got, ok := TryRotate(b, p, RotateCW)

	if !ok {
		t.Fatal("I-piece on left wall must be rotated kick")
	}
	if got.Rotation != 2 {
		t.Errorf("rotation: got %d, want 2", got.Rotation)
	}
	if got.Position.X <= p.Position.X {
		t.Errorf("wait kick on right: X=%d, X=%d",
			p.Position.X, got.Position.X)
	}
}

func TestTryRotate_BlockedAllKicks(t *testing.T) {
	b := NewBoard()

	for y := 0; y < BoardHeight; y++ {
		for x := 0; x < BoardWidth; x++ {
			b.Set(Point{X: x, Y: y}, PieceT)
		}
	}

	center := Point{X: 5, Y: 10}
	for _, offset := range []Point{{0, 1}, {-1, 0}, {0, 0}, {1, 0}} {
		b.Set(Point{X: center.X + offset.X, Y: center.Y + offset.Y}, PieceNone)
	}

	p := Piece{Type: PieceT, Position: center, Rotation: 0}

	got, ok := TryRotate(b, p, RotateCW)

	if ok {
		t.Fatal("wait fatal")
	}
	if got != p {
		t.Errorf("change after failed rotate: %+v: %+v", p, got)
	}
}

func TestTryRotate_OPieceAlwaysSucceeds(t *testing.T) {
	b := buildBoard(BoardHeight)

	p := Piece{Type: PieceO, Position: Point{X: 5, Y: 5}, Rotation: 0}

	got, ok := TryRotate(b, p, RotateCW)

	if !ok {
		t.Fatal("O-piece always can rotate")
	}
	if got.Rotation != 1 {
		t.Errorf("rotation: got %d, want 1", got.Rotation)
	}
}

func TestTryRotate_RotationWrapsAround(t *testing.T) {
	b := NewBoard()
	p := Piece{Type: PieceT, Position: Point{X: 5, Y: 10}, Rotation: 0}

	for i := range 4 {
		var ok bool
		p, ok = TryRotate(b, p, RotateCW)
		if !ok {
			t.Fatalf("roatation %d failed", i+1)
		}
	}

	if p.Rotation != 0 {
		t.Errorf("after 4xCW: rotation=%d, want 0", p.Rotation)
	}
}

func TestTryRotate_TSpin(t *testing.T) {
	b := NewBoard()

	type cellFill struct{ x, y int }
	holes := map[cellFill]bool{
		{3, 18}: true,
		{4, 18}: true,
		{4, 19}: true,
		{4, 20}: true,
	}

	for y := 18; y <= 20; y++ {
		for x := 0; x < BoardWidth; x++ {
			if !holes[cellFill{x, y}] {
				b.Set(Point{X: x, Y: y}, PieceT)
			}
		}
	}

	p := Piece{Type: PieceT, Position: Point{X: 4, Y: 17}, Rotation: 2}

	_, ok := TryRotate(b, p, RotateCW)
	if !ok {
		t.Error("T-Spin failed")
	}
}
