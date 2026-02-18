package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct{}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Заливаем фон темно-серым
	screen.Fill(color.RGBA{0x2b, 0x2b, 0x2b, 0xff})
	ebitenutil.DebugPrint(screen, "WASM Tetris")
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 480
}

func main() {
	ebiten.SetWindowSize(320, 480)
	ebiten.SetWindowTitle("Go-Tetris")

	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
