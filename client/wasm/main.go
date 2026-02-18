package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/coder/websocket"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const wsURL = "ws://localhost:8081/ws"

type Game struct {
	conn      *websocket.Conn
	msg       string
	connected bool
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	status := "Status: Disconnected"
	if g.connected {
		status = "Status: Connected to Gateway!"
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf("%s\nLast Msg: %s", status, g.msg))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	g := &Game{}

	go g.connectToGateway()

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Tetris WASM Client")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

func (g *Game) connectToGateway() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	fmt.Printf("Connecting to %s...\n", wsURL)

	c, _, err := websocket.Dial(ctx, wsURL, &websocket.DialOptions{
		CompressionMode: websocket.CompressionDisabled,
	})
	if err != nil {
		g.msg = "Error connecting: " + err.Error()
		log.Println("Connection error:", err)
		return
	}

	g.conn = c
	g.connected = true
	g.msg = "Handshake success!"
	log.Println("Connected to WebSocket!")

	go g.readLoop()
}

func (g *Game) readLoop() {
	for {
		_, data, err := g.conn.Read(context.Background())
		if err != nil {
			g.msg = "Read error: " + err.Error()
			g.connected = false
			return
		}
		g.msg = string(data)
		fmt.Printf("Received: %s\n", g.msg)
	}
}
