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

	currentPiece := core.Piece{
		Type:     core.PieceT,
		Position: core.Point{X: 5, Y: 5},
		Rotation: 0,
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		draw(board, currentPiece)

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		nextPiece := currentPiece

		switch input {
		case "a":
			nextPiece.Position.X--
		case "d":
			nextPiece.Position.X++
		case "s":
			nextPiece.Position.Y++
		case "w":
			nextPiece.Rotation++
		case "q":
			return
		}

		if !board.HasCollision(nextPiece) {
			currentPiece = nextPiece
		} else {
			fmt.Println("Err: collision.")
		}
	}
}

func draw(board *core.Board, p core.Piece) {
	fmt.Print("\033[H\033[2J")

	display := make([]string, core.BoardHeight)

	minos := core.GetRotatedMinos(p.Type, p.Rotation)

	for y := 0; y < core.BoardHeight; y++ {
		row := ""
		for x := 0; x < core.BoardWidth; x++ {
			char := " . "

			if board.Get(core.Point{X: x, Y: y}) != core.PieceNone {
				char = "[#]"
			}

			for _, mino := range minos {
				absolute := p.Position.Add(mino)
				if absolute.X == x && absolute.Y == y {
					char = "[@]"
				}
			}
			row += char
		}
		display[y] = fmt.Sprintf("|%s|", row)
	}

	fmt.Println(strings.Join(display, "\n"))
	fmt.Println("================================")
	fmt.Printf("Pos: (%d, %d), Rot: %d\n", p.Position.X, p.Position.Y, p.Rotation)
}
