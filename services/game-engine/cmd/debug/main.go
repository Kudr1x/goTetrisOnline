package main

import (
	"GoTetrisOnline/pkg/core"
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	board := core.NewBoard()
	currentPiece := spawnPiece()
	score := 0
	linesTotal := 0

	reader := bufio.NewReader(os.Stdin)

	draw(board, currentPiece, score, linesTotal)

	for {
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "q" {
			break
		}

		nextPiece := currentPiece
		moveApplied := false
		shouldLock := false

		switch input {
		case "a":
			nextPiece.Position.X--
			if !board.HasCollision(nextPiece) {
				currentPiece = nextPiece
				moveApplied = true
			}
		case "d":
			nextPiece.Position.X++
			if !board.HasCollision(nextPiece) {
				currentPiece = nextPiece
				moveApplied = true
			}
		case "w":
			nextPiece.Rotation++
			if !board.HasCollision(nextPiece) {
				currentPiece = nextPiece
				moveApplied = true
			}
		case "s":
			nextPiece.Position.Y++
			if !board.HasCollision(nextPiece) {
				currentPiece = nextPiece
				moveApplied = true
			} else {
				shouldLock = true
			}
		}

		if shouldLock {
			board.LockPiece(currentPiece)

			cleared := board.ClearLines()
			linesTotal += cleared

			switch cleared {
			case 1:
				score += 40
			case 2:
				score += 100
			case 3:
				score += 300
			case 4:
				score += 1200
			}

			currentPiece = spawnPiece()

			if board.HasCollision(currentPiece) {
				draw(board, currentPiece, score, linesTotal)
				return
			}
		}

		draw(board, currentPiece, score, linesTotal)

		if !moveApplied && !shouldLock {
			fmt.Println("block")
		}
	}
}

func spawnPiece() core.Piece {
	return core.Piece{
		Type:     core.PieceT,
		Position: core.Point{X: 4, Y: 1},
		Rotation: 0,
	}
}

func draw(board *core.Board, p core.Piece, score, lines int) {
	fmt.Print("\033[H\033[2J")

	display := make([]string, core.BoardHeight)
	minos := core.GetRotatedMinos(p.Type, p.Rotation)

	for y := 0; y < core.BoardHeight; y++ {
		row := ""
		for x := 0; x < core.BoardWidth; x++ {
			char := " . "

			cell := board.Get(core.Point{X: x, Y: y})
			if cell != core.PieceNone {
				char = "[#]"
			}

			for _, mino := range minos {
				abs := p.Position.Add(mino)
				if abs.X == x && abs.Y == y {
					char = "[@]"
				}
			}
			row += char
		}
		display[y] = fmt.Sprintf("|%2d| |%s|", y, row)
	}

	fmt.Println(strings.Join(display, "\n"))
}
