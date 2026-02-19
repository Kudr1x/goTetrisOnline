package domain

import (
	"GoTetrisOnline/pkg/core"
	"sync"
	"time"
)

const (
	StatusWaiting GameStatus = iota
	StatusRunning
	StatusFinished
)

type GameStatus int

type GameEvent struct {
	Type    string
	Payload interface{}
}

type GameStateDTO struct {
	Score        int32
	Level        int32
	Grid         []byte
	CurrentPiece core.Piece
	NextPieces   []core.PieceType
}

type Game struct {
	mu sync.RWMutex

	Board        *core.Board
	CurrentPiece core.Piece
	NextPieces   []core.PieceType

	Score int32
	Level int32
	Lines int32

	UID    string
	Status GameStatus

	events chan GameEvent
	quit   chan struct{}

	bag *core.Bag
}

func NewGame(uid string) *Game {
	return &Game{
		UID:    uid,
		Status: StatusWaiting,
		Board:  core.NewBoard(),
		events: make(chan GameEvent, 100),
		quit:   make(chan struct{}),
		bag:    core.NewBag(),
	}
}

func (g *Game) Start() {
	g.mu.Lock()
	g.Status = StatusRunning

	g.CurrentPiece = g.spawnPiece()
	g.mu.Unlock()

	go g.loop()
}

func (g *Game) Stop() {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.Status != StatusFinished {
		g.Status = StatusFinished
		close(g.quit)
		close(g.events)
	}
}

func (g *Game) Events() <-chan GameEvent {
	return g.events
}

func (g *Game) loop() {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-g.quit:
			return
		case <-ticker.C:
			g.ApplyGravity()
		}
	}
}

func (g *Game) ApplyGravity() {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.Status != StatusRunning {
		return
	}

	next := g.CurrentPiece
	next.Position.Y++

	if !g.Board.HasCollision(next) {
		g.CurrentPiece = next
		g.broadcast()
	} else {
		g.lockAndSpawn()
		g.broadcast()
	}
}

func (g *Game) GetSnapshot() GameStateDTO {
	return GameStateDTO{
		Score:        g.Score,
		Level:        g.Level,
		Grid:         g.Board.ToBytes(),
		CurrentPiece: g.CurrentPiece,
		NextPieces:   g.bag.Peek(3),
	}
}

func (g *Game) broadcast() {
	if g.Status == StatusFinished {
		return
	}

	select {
	case g.events <- GameEvent{Type: "state_update", Payload: g.GetSnapshot()}:
	default:
		// skip
	}
}

func (g *Game) lockAndSpawn() {
	g.Board.LockPiece(g.CurrentPiece)

	lines := g.Board.ClearLines()
	g.updateScore(lines)

	g.CurrentPiece = g.spawnPiece()

	if g.Board.HasCollision(g.CurrentPiece) {
		g.Status = StatusFinished
		g.events <- GameEvent{Type: "game_over", Payload: g.Score}
		close(g.quit)
		close(g.events)
	}
}

func (g *Game) updateScore(lines int32) {
	switch lines {
	case 1:
		g.Score += 40 * (g.Level + 1)
	case 2:
		g.Score += 100 * (g.Level + 1)
	case 3:
		g.Score += 300 * (g.Level + 1)
	case 4:
		g.Score += 1200 * (g.Level + 1)
	}
	g.Lines += lines
}

func (g *Game) MoveLeft() {
	g.mu.Lock()
	defer g.mu.Unlock()

	next := g.CurrentPiece
	next.Position.X--

	if !g.Board.HasCollision(next) {
		g.CurrentPiece = next
		g.broadcast()
	}
}

func (g *Game) MoveRight() {
	g.mu.Lock()
	defer g.mu.Unlock()

	next := g.CurrentPiece
	next.Position.X++

	if !g.Board.HasCollision(next) {
		g.CurrentPiece = next
		g.broadcast()
	}
}

func (g *Game) Rotate(direction int) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.Status != StatusRunning {
		return
	}

	rotated, ok := core.TryRotate(g.Board, g.CurrentPiece, direction)
	if ok {
		g.CurrentPiece = rotated
		g.broadcast()
	}
}

func (g *Game) HardDrop() {
	g.mu.Lock()
	defer g.mu.Unlock()

	for {
		next := g.CurrentPiece
		next.Position.Y++

		if g.Board.HasCollision(next) {
			g.lockAndSpawn()
			g.broadcast()
			return
		}

		g.CurrentPiece = next
	}
}

func (g *Game) spawnPiece() core.Piece {
	return core.Piece{
		Type:     g.bag.Next(),
		Position: core.Point{X: 4, Y: 0},
		Rotation: 0,
	}
}
