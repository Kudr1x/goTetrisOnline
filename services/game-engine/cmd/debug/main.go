package main

import (
	"GoTetrisOnline/pkg/core"
	"bufio"
	"fmt"
	"os"
	"strings"
)

type State struct {
	Board        *core.Board
	CurrentPiece core.Piece
	Score        int
	Lines        int
	GameOver     bool
}

func main() {
	state := State{
		Board:        core.NewBoard(),
		CurrentPiece: spawnPiece(),
	}

	reader := bufio.NewReader(os.Stdin)
	draw(state)

	for {
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "q" || state.GameOver {
			break
		}

		target := getTargetPos(state.CurrentPiece, input)

		if !state.Board.HasCollision(target) {
			state.CurrentPiece = target
		} else if input == "s" {
			handleLock(&state)
		}

		draw(state)
	}
}

func getTargetPos(p core.Piece, input string) core.Piece {
	np := p
	switch input {
	case "a":
		np.Position.X--
	case "d":
		np.Position.X++
	case "s":
		np.Position.Y++
	case "w":
		np.Rotation++
	}
	return np
}

func handleLock(s *State) {
	s.Board.LockPiece(s.CurrentPiece)

	cleared := s.Board.ClearLines()
	s.Lines += cleared
	s.Score += calculateScore(cleared)

	s.CurrentPiece = spawnPiece()

	if s.Board.HasCollision(s.CurrentPiece) {
		s.GameOver = true
		fmt.Println("game over")
	}
}

func calculateScore(lines int) int {
	switch lines {
	case 1:
		return 40
	case 2:
		return 100
	case 3:
		return 300
	case 4:
		return 1200
	default:
		return 0
	}
}

func spawnPiece() core.Piece {
	return core.Piece{
		// todo
		Type:     core.PieceT,
		Position: core.Point{X: 4, Y: 1},
		Rotation: 0,
	}
}

func draw(s State) {
	fmt.Print("\033[H\033[2J")
	display := make([]string, core.BoardHeight)
	minos := core.GetRotatedMinos(s.CurrentPiece.Type, s.CurrentPiece.Rotation)

	for y := 0; y < core.BoardHeight; y++ {
		row := ""
		for x := 0; x < core.BoardWidth; x++ {
			char := " . "
			if s.Board.Get(core.Point{X: x, Y: y}) != core.PieceNone {
				char = "[#]"
			}

			for _, mino := range minos {
				abs := s.CurrentPiece.Position.Add(mino)
				if abs.X == x && abs.Y == y {
					char = "[@]"
				}
			}
			row += char
		}
		display[y] = fmt.Sprintf("|%2d| |%s|", y, row)
	}

	fmt.Println(strings.Join(display, "\n"))
	fmt.Printf("Score: %d | Lines: %d\n", s.Score, s.Lines)
}
