package main

import (
	pb "GoTetrisOnline/api/proto/game/v1"
	"GoTetrisOnline/pkg/core"
	"GoTetrisOnline/pkg/renderer"
	"context"
	"fmt"
	"image/color"
	"log"
	"time"

	"github.com/coder/websocket"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"google.golang.org/protobuf/proto"
)

const (
	wsURL    = "ws://localhost:8081/ws"
	cellSize = 20
	boardX   = 20
	boardY   = 20
	sidebarX = boardX + core.BoardWidth*cellSize + 40
	screenW  = 640
	screenH  = 480
)

type Game struct {
	conn          *websocket.Conn
	state         *pb.StateUpdate
	connected     bool
	err           error
	lastInput     time.Time
	inputCooldown time.Duration
	ctx           context.Context
	cancel        context.CancelFunc
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		return fmt.Errorf("quit")
	}

	if g.conn == nil || g.state == nil {
		return nil
	}

	if time.Since(g.lastInput) < g.inputCooldown {
		return nil
	}

	var input pb.InputType
	sendInput := false

	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		input = pb.InputType_INPUT_LEFT
		sendInput = true
	} else if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		input = pb.InputType_INPUT_RIGHT
		sendInput = true
	} else if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		input = pb.InputType_INPUT_ROTATE_CW
		sendInput = true
	} else if ebiten.IsKeyPressed(ebiten.KeyE) {
		input = pb.InputType_INPUT_ROTATE_CCW
		sendInput = true
	} else if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		input = pb.InputType_INPUT_SOFT_DROP
		sendInput = true
	} else if ebiten.IsKeyPressed(ebiten.KeySpace) {
		input = pb.InputType_INPUT_HARD_DROP
		sendInput = true
	}

	if sendInput {
		g.lastInput = time.Now()
		msg := &pb.ClientMessage{
			Payload: &pb.ClientMessage_Input{
				Input: &pb.InputRequest{Input: input},
			},
		}
		data, err := proto.Marshal(msg)
		if err != nil {
			log.Printf("Marshal error: %v", err)
			return nil
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		if err := g.conn.Write(ctx, websocket.MessageBinary, data); err != nil {
			log.Printf("Write error: %v", err)
			return nil
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{20, 20, 20, 255})

	if g.err != nil {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("Error: %v", g.err))
		return
	}

	if !g.connected {
		ebitenutil.DebugPrint(screen, "Connecting...")
		return
	}

	if g.state == nil {
		ebitenutil.DebugPrint(screen, "Waiting for game state...")
		return
	}

	view := renderer.StateToView(g.state)
	g.drawBoard(screen, view)
	g.drawSidebar(screen, view)
}

func (g *Game) drawBoard(screen *ebiten.Image, view *renderer.GameView) {
	for y := 0; y < view.Height; y++ {
		for x := 0; x < view.Width; x++ {
			cell := view.Board[y][x]
			var clr color.RGBA

			switch cell.Type {
			case renderer.CellPiece, renderer.CellFixed:
				clr = renderer.GetPieceColor(cell.PieceType)
			default:
				clr = color.RGBA{40, 40, 40, 255}
			}

			px := float32(boardX + x*cellSize)
			py := float32(boardY + y*cellSize)
			vector.DrawFilledRect(screen, px, py, cellSize-1, cellSize-1, clr, false)
		}
	}
}

func (g *Game) drawSidebar(screen *ebiten.Image, view *renderer.GameView) {
	y := boardY
	ebitenutil.DebugPrintAt(screen, "NEXT:", sidebarX, y)
	y += 30

	if view.NextPiece != core.PieceNone {
		nextGrid := renderer.RenderNextPieceGrid(view.NextPiece)
		clr := renderer.GetPieceColor(view.NextPiece)

		for py := 0; py < nextGrid.Size; py++ {
			for px := 0; px < nextGrid.Size; px++ {
				if nextGrid.Grid[py][px] {
					fx := float32(sidebarX + px*cellSize)
					fy := float32(y + py*cellSize)
					vector.DrawFilledRect(screen, fx, fy, cellSize-1, cellSize-1, clr, false)
				}
			}
		}
		y += nextGrid.Size*cellSize + 30
	}

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Score: %d", view.Score), sidebarX, y)
	y += 20
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Level: %d", view.Level), sidebarX, y)
	y += 40
	ebitenutil.DebugPrintAt(screen, "CONTROLS:", sidebarX, y)
	y += 20
	ebitenutil.DebugPrintAt(screen, "A/←: Left", sidebarX, y)
	y += 15
	ebitenutil.DebugPrintAt(screen, "D/→: Right", sidebarX, y)
	y += 15
	ebitenutil.DebugPrintAt(screen, "W/↑: Rotate", sidebarX, y)
	y += 15
	ebitenutil.DebugPrintAt(screen, "Space: Drop", sidebarX, y)
	y += 15
	ebitenutil.DebugPrintAt(screen, "Q: Quit", sidebarX, y)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenW, screenH
}

func main() {
	g := &Game{inputCooldown: 150 * time.Millisecond}
	go g.connectToGateway()

	ebiten.SetWindowSize(screenW, screenH)
	ebiten.SetWindowTitle("Tetris WASM Client")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

func (g *Game) connectToGateway() {
	dialCtx, dialCancel := context.WithTimeout(context.Background(), time.Minute)
	defer dialCancel()

	c, _, err := websocket.Dial(dialCtx, wsURL, nil)
	if err != nil {
		g.err = err
		log.Printf("Connection error: %v", err)
		return
	}

	g.conn = c
	g.connected = true
	g.ctx, g.cancel = context.WithCancel(context.Background())

	joinMsg := &pb.ClientMessage{
		Payload: &pb.ClientMessage_Join{
			Join: &pb.JoinRequest{MatchId: "room-1", Token: "token"},
		},
	}
	data, _ := proto.Marshal(joinMsg)
	_ = c.Write(g.ctx, websocket.MessageBinary, data)

	go g.readLoop()
}

func (g *Game) readLoop() {
	defer func() {
		if g.cancel != nil {
			g.cancel()
		}
		if g.conn != nil {
			g.conn.Close(websocket.StatusNormalClosure, "")
		}
	}()

	for {
		select {
		case <-g.ctx.Done():
			return
		default:
		}

		_, data, err := g.conn.Read(g.ctx)
		if err != nil {
			if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
				log.Println("Connection closed normally")
				return
			}
			log.Printf("Read error: %v", err)
			g.err = err
			g.connected = false
			return
		}

		var msg pb.ServerMessage
		if err := proto.Unmarshal(data, &msg); err != nil {
			log.Printf("Unmarshal error: %v", err)
			continue
		}

		switch payload := msg.Payload.(type) {
		case *pb.ServerMessage_State:
			g.state = payload.State
		case *pb.ServerMessage_Event:
			if payload.Event.Type == pb.EventType_EVENT_GAME_OVER {
				log.Println("Game Over!")
			}
		}
	}
}
