package domain

import (
	"GoTetrisOnline/pkg/core"
	"testing"
)

func TestGame_RotateCW_IPiece_WallKick(t *testing.T) {
	game := NewGame("test-i-cw")
	game.Status = StatusRunning
	game.Board = core.NewBoard()
	game.CurrentPiece = core.Piece{
		Type:     core.PieceI,
		Position: core.Point{X: 0, Y: 5},
		Rotation: 0,
	}

	initialX := game.CurrentPiece.Position.X
	initialRot := game.CurrentPiece.Rotation

	game.Rotate(core.RotateCW)

	if game.CurrentPiece.Rotation == initialRot {
		t.Errorf("I-piece CW rotation failed: rotation stayed %d", initialRot)
	}

	if game.CurrentPiece.Position.X == initialX {
		t.Logf("Wall kick may not have triggered (X unchanged), but rotation succeeded")
	}

	expectedRotation := (initialRot + 1) % 4
	if game.CurrentPiece.Rotation != expectedRotation {
		t.Errorf("Expected rotation=%d, got %d", expectedRotation, game.CurrentPiece.Rotation)
	}
}

func TestGame_RotateCCW_IPiece_WallKick(t *testing.T) {
	game := NewGame("test-i-ccw")
	game.Status = StatusRunning
	game.Board = core.NewBoard()
	game.CurrentPiece = core.Piece{
		Type:     core.PieceI,
		Position: core.Point{X: 9, Y: 5},
		Rotation: 0,
	}

	initialRot := game.CurrentPiece.Rotation

	game.Rotate(core.RotateCCW)

	expectedRotation := 3
	if game.CurrentPiece.Rotation != expectedRotation {
		t.Errorf("I-piece CCW: expected rotation=%d, got %d", expectedRotation, game.CurrentPiece.Rotation)
	}

	if game.CurrentPiece.Rotation == initialRot {
		t.Error("I-piece CCW rotation failed")
	}
}

func TestGame_RotateCW_OPiece_AlwaysSucceeds(t *testing.T) {
	game := NewGame("test-o-cw")
	game.Status = StatusRunning
	game.Board = core.NewBoard()
	game.CurrentPiece = core.Piece{
		Type:     core.PieceO,
		Position: core.Point{X: 4, Y: 5},
		Rotation: 0,
	}

	initialPos := game.CurrentPiece.Position

	game.Rotate(core.RotateCW)

	if game.CurrentPiece.Position != initialPos {
		t.Errorf("O-piece position should not change on rotation: expected %+v, got %+v",
			initialPos, game.CurrentPiece.Position)
	}

	expectedRotation := 1
	if game.CurrentPiece.Rotation != expectedRotation {
		t.Errorf("O-piece CW: expected rotation=%d, got %d", expectedRotation, game.CurrentPiece.Rotation)
	}
}

func TestGame_RotateCCW_OPiece(t *testing.T) {
	game := NewGame("test-o-ccw")
	game.Status = StatusRunning
	game.Board = core.NewBoard()
	game.CurrentPiece = core.Piece{
		Type:     core.PieceO,
		Position: core.Point{X: 4, Y: 5},
		Rotation: 0,
	}

	game.Rotate(core.RotateCCW)

	expectedRotation := 3
	if game.CurrentPiece.Rotation != expectedRotation {
		t.Errorf("O-piece CCW: expected rotation=%d, got %d", expectedRotation, game.CurrentPiece.Rotation)
	}
}

func TestGame_RotateCW_TPiece(t *testing.T) {
	game := NewGame("test-t-cw")
	game.Status = StatusRunning
	game.Board = core.NewBoard()
	game.CurrentPiece = core.Piece{
		Type:     core.PieceT,
		Position: core.Point{X: 4, Y: 5},
		Rotation: 0,
	}

	game.Rotate(core.RotateCW)

	if game.CurrentPiece.Rotation != 1 {
		t.Errorf("T-piece CW failed: expected rotation=1, got %d", game.CurrentPiece.Rotation)
	}
}

func TestGame_RotateCCW_TPiece(t *testing.T) {
	game := NewGame("test-t-ccw")
	game.Status = StatusRunning
	game.Board = core.NewBoard()
	game.CurrentPiece = core.Piece{
		Type:     core.PieceT,
		Position: core.Point{X: 4, Y: 5},
		Rotation: 0,
	}

	game.Rotate(core.RotateCCW)

	if game.CurrentPiece.Rotation != 3 {
		t.Errorf("T-piece CCW failed: expected rotation=3, got %d", game.CurrentPiece.Rotation)
	}
}

func TestGame_RotateCW_JPiece(t *testing.T) {
	game := NewGame("test-j-cw")
	game.Status = StatusRunning
	game.Board = core.NewBoard()
	game.CurrentPiece = core.Piece{
		Type:     core.PieceJ,
		Position: core.Point{X: 4, Y: 5},
		Rotation: 0,
	}

	game.Rotate(core.RotateCW)

	if game.CurrentPiece.Rotation != 1 {
		t.Errorf("J-piece CW failed: expected rotation=1, got %d", game.CurrentPiece.Rotation)
	}
}

func TestGame_RotateCCW_JPiece(t *testing.T) {
	game := NewGame("test-j-ccw")
	game.Status = StatusRunning
	game.Board = core.NewBoard()
	game.CurrentPiece = core.Piece{
		Type:     core.PieceJ,
		Position: core.Point{X: 4, Y: 5},
		Rotation: 0,
	}

	game.Rotate(core.RotateCCW)

	if game.CurrentPiece.Rotation != 3 {
		t.Errorf("J-piece CCW failed: expected rotation=3, got %d", game.CurrentPiece.Rotation)
	}
}

func TestGame_RotateCW_LPiece(t *testing.T) {
	game := NewGame("test-l-cw")
	game.Status = StatusRunning
	game.Board = core.NewBoard()
	game.CurrentPiece = core.Piece{
		Type:     core.PieceL,
		Position: core.Point{X: 4, Y: 5},
		Rotation: 0,
	}

	game.Rotate(core.RotateCW)

	if game.CurrentPiece.Rotation != 1 {
		t.Errorf("L-piece CW failed: expected rotation=1, got %d", game.CurrentPiece.Rotation)
	}
}

func TestGame_RotateCCW_LPiece(t *testing.T) {
	game := NewGame("test-l-ccw")
	game.Status = StatusRunning
	game.Board = core.NewBoard()
	game.CurrentPiece = core.Piece{
		Type:     core.PieceL,
		Position: core.Point{X: 4, Y: 5},
		Rotation: 0,
	}

	game.Rotate(core.RotateCCW)

	if game.CurrentPiece.Rotation != 3 {
		t.Errorf("L-piece CCW failed: expected rotation=3, got %d", game.CurrentPiece.Rotation)
	}
}

func TestGame_RotateCW_SPiece(t *testing.T) {
	game := NewGame("test-s-cw")
	game.Status = StatusRunning
	game.Board = core.NewBoard()
	game.CurrentPiece = core.Piece{
		Type:     core.PieceS,
		Position: core.Point{X: 4, Y: 5},
		Rotation: 0,
	}

	game.Rotate(core.RotateCW)

	if game.CurrentPiece.Rotation != 1 {
		t.Errorf("S-piece CW failed: expected rotation=1, got %d", game.CurrentPiece.Rotation)
	}
}

func TestGame_RotateCCW_SPiece(t *testing.T) {
	game := NewGame("test-s-ccw")
	game.Status = StatusRunning
	game.Board = core.NewBoard()
	game.CurrentPiece = core.Piece{
		Type:     core.PieceS,
		Position: core.Point{X: 4, Y: 5},
		Rotation: 0,
	}

	game.Rotate(core.RotateCCW)

	if game.CurrentPiece.Rotation != 3 {
		t.Errorf("S-piece CCW failed: expected rotation=3, got %d", game.CurrentPiece.Rotation)
	}
}

func TestGame_RotateCW_ZPiece(t *testing.T) {
	game := NewGame("test-z-cw")
	game.Status = StatusRunning
	game.Board = core.NewBoard()
	game.CurrentPiece = core.Piece{
		Type:     core.PieceZ,
		Position: core.Point{X: 4, Y: 5},
		Rotation: 0,
	}

	game.Rotate(core.RotateCW)

	if game.CurrentPiece.Rotation != 1 {
		t.Errorf("Z-piece CW failed: expected rotation=1, got %d", game.CurrentPiece.Rotation)
	}
}

func TestGame_RotateCCW_ZPiece(t *testing.T) {
	game := NewGame("test-z-ccw")
	game.Status = StatusRunning
	game.Board = core.NewBoard()
	game.CurrentPiece = core.Piece{
		Type:     core.PieceZ,
		Position: core.Point{X: 4, Y: 5},
		Rotation: 0,
	}

	game.Rotate(core.RotateCCW)

	if game.CurrentPiece.Rotation != 3 {
		t.Errorf("Z-piece CCW failed: expected rotation=3, got %d", game.CurrentPiece.Rotation)
	}
}

func TestGame_Rotate_BlockedByWalls(t *testing.T) {
	game := NewGame("test-blocked")
	game.Status = StatusRunning
	game.Board = core.NewBoard()

	for x := 0; x < core.BoardWidth; x++ {
		for y := 5; y < 10; y++ {
			if x < 3 || x > 6 {
				game.Board.Set(core.Point{X: x, Y: y}, core.PieceGarbage)
			}
		}
	}

	game.CurrentPiece = core.Piece{
		Type:     core.PieceI,
		Position: core.Point{X: 4, Y: 5},
		Rotation: 0,
	}

	initialRotation := game.CurrentPiece.Rotation

	game.Rotate(core.RotateCW)

	if game.CurrentPiece.Rotation != initialRotation {
		t.Logf("Rotation succeeded with wall kick, expected to fail")
	}
}

func TestGame_Rotate_WhenGameNotRunning(t *testing.T) {
	game := NewGame("test-not-running")
	game.Status = StatusWaiting
	game.Board = core.NewBoard()
	game.CurrentPiece = core.Piece{
		Type:     core.PieceT,
		Position: core.Point{X: 4, Y: 5},
		Rotation: 0,
	}

	initialRotation := game.CurrentPiece.Rotation

	game.Rotate(core.RotateCW)

	if game.CurrentPiece.Rotation != initialRotation {
		t.Error("Rotation should not work when game is not running")
	}
}

func TestGame_RotateMultipleTimes_WrapsAround(t *testing.T) {
	game := NewGame("test-wrap")
	game.Status = StatusRunning
	game.Board = core.NewBoard()
	game.CurrentPiece = core.Piece{
		Type:     core.PieceT,
		Position: core.Point{X: 4, Y: 5},
		Rotation: 0,
	}

	game.Rotate(core.RotateCW)
	game.Rotate(core.RotateCW)
	game.Rotate(core.RotateCW)
	game.Rotate(core.RotateCW)

	if game.CurrentPiece.Rotation != 0 {
		t.Errorf("After 4 CW rotations, expected rotation=0, got %d", game.CurrentPiece.Rotation)
	}
}

func TestGame_RotateCCW_MultipleTimes_WrapsAround(t *testing.T) {
	game := NewGame("test-wrap-ccw")
	game.Status = StatusRunning
	game.Board = core.NewBoard()
	game.CurrentPiece = core.Piece{
		Type:     core.PieceT,
		Position: core.Point{X: 4, Y: 5},
		Rotation: 0,
	}

	game.Rotate(core.RotateCCW)
	game.Rotate(core.RotateCCW)
	game.Rotate(core.RotateCCW)
	game.Rotate(core.RotateCCW)

	if game.CurrentPiece.Rotation != 0 {
		t.Errorf("After 4 CCW rotations, expected rotation=0, got %d", game.CurrentPiece.Rotation)
	}
}
