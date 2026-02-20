package renderer

import (
	pb "GoTetrisOnline/api/proto/game/v1"
	"GoTetrisOnline/pkg/core"
	"image/color"
)

type CellType int

const (
	CellEmpty CellType = iota
	CellPiece
	CellGhost
	CellFixed
)

type Cell struct {
	Type      CellType
	PieceType core.PieceType
}

type GameView struct {
	Board     [][]Cell
	NextPiece core.PieceType
	Score     int32
	Level     int32
	Width     int
	Height    int
}

func StateToView(state *pb.StateUpdate) *GameView {
	view := &GameView{
		Score:  state.Score,
		Level:  state.Level,
		Width:  core.BoardWidth,
		Height: core.BoardHeight - core.Space,
	}

	view.Board = make([][]Cell, view.Height)
	for i := range view.Board {
		view.Board[i] = make([]Cell, view.Width)
	}

	currentPiece := core.Piece{
		Type: core.PieceType(state.CurrentPiece.Type), //nolint:gosec
		Position: core.Point{
			X: int(state.CurrentPiece.X),
			Y: int(state.CurrentPiece.Y),
		},
		Rotation: int(state.CurrentPiece.Rotation),
	}

	minos := core.GetRotatedMinos(currentPiece.Type, currentPiece.Rotation)

	for y := core.Space; y < core.BoardHeight; y++ {
		for x := 0; x < core.BoardWidth; x++ {
			viewY := y - core.Space
			cell := Cell{Type: CellEmpty}

			idx := y*core.BoardWidth + x
			if idx < len(state.Grid) && state.Grid[idx] != 0 {
				cell.Type = CellFixed
				cell.PieceType = core.PieceType(state.Grid[idx]) //nolint:gosec
			}

			for _, mino := range minos {
				absX := currentPiece.Position.X + mino.X
				absY := currentPiece.Position.Y + mino.Y

				if absX == x && absY == y {
					cell.Type = CellPiece
					cell.PieceType = currentPiece.Type
					break
				}
			}

			view.Board[viewY][x] = cell
		}
	}

	if len(state.NextPieces) > 0 {
		view.NextPiece = core.PieceType(state.NextPieces[0]) //nolint:gosec
	}

	return view
}

func GetPieceColor(t core.PieceType) color.RGBA {
	switch t {
	case core.PieceI:
		return color.RGBA{0, 255, 255, 255} // Cyan
	case core.PieceO:
		return color.RGBA{255, 255, 0, 255} // Yellow
	case core.PieceT:
		return color.RGBA{160, 32, 240, 255} // Purple
	case core.PieceS:
		return color.RGBA{0, 255, 0, 255} // Green
	case core.PieceZ:
		return color.RGBA{255, 0, 0, 255} // Red
	case core.PieceJ:
		return color.RGBA{0, 0, 255, 255} // Blue
	case core.PieceL:
		return color.RGBA{255, 165, 0, 255} // Orange
	default:
		return color.RGBA{128, 128, 128, 255} // Gray
	}
}

type NextPieceGrid struct {
	Grid [][]bool
	Size int
}

func RenderNextPieceGrid(t core.PieceType) *NextPieceGrid {
	minos := core.GetRotatedMinos(t, 0)

	minX, maxX := 0, 0
	minY, maxY := 0, 0
	for i, m := range minos {
		if i == 0 {
			minX, maxX = m.X, m.X
			minY, maxY = m.Y, m.Y
		} else {
			if m.X < minX {
				minX = m.X
			}
			if m.X > maxX {
				maxX = m.X
			}
			if m.Y < minY {
				minY = m.Y
			}
			if m.Y > maxY {
				maxY = m.Y
			}
		}
	}

	width := maxX - minX + 1
	height := maxY - minY + 1
	gridSize := 4
	offsetX := (gridSize - width) / 2
	offsetY := (gridSize - height) / 2

	grid := make([][]bool, gridSize)
	for i := range grid {
		grid[i] = make([]bool, gridSize)
	}

	for _, m := range minos {
		x := m.X - minX + offsetX
		y := m.Y - minY + offsetY
		if x >= 0 && x < gridSize && y >= 0 && y < gridSize {
			grid[y][x] = true
		}
	}

	return &NextPieceGrid{
		Grid: grid,
		Size: gridSize,
	}
}
