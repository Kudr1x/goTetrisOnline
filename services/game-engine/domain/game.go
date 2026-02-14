package domain

import (
	"GoTetrisOnline/pkg/core"
	"sync"
	"time"
)

type GameStatus int

type GameEvent struct {
	Type    string
	Payload interface{}
}

const (
	StatusWaiting GameStatus = iota
	StatusRunning
	StatusFinished
)

type Game struct {
	mu sync.RWMutex

	Board        *core.Board
	CurrentPiece core.Piece
	NextPieces   []core.PieceType

	Score int
	Level int

	UID    string
	Status GameStatus

	events chan GameEvent
	quit   chan struct{}
}

func NewGame(uid string) *Game {
	return &Game{
		UID:    uid,
		Status: StatusWaiting,
		Board:  core.NewBoard(),
		events: make(chan GameEvent, 100),
		quit:   make(chan struct{}),
	}
}

func (g *Game) Start() {
	g.mu.Lock()
	g.Status = StatusRunning

	// g.spawnPiece()
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

	if g.Status == StatusFinished {
		return
	}

	select {
	case g.events <- GameEvent{Type: "state_update", Payload: g.GetSnapshot()}:
	default:
		// todo
	}
}

func (g *Game) GetSnapshot() map[string]interface{} {
	return map[string]interface{}{
		"score": g.Score,
	}

	// todo
}
