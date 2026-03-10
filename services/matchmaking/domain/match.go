package domain

import (
	"time"
)

type GameMode int

const (
	GameModeSolo GameMode = iota + 1
	GameModeOneVsOne
)

type MatchStatus int

const (
	MatchStatusWaiting MatchStatus = iota + 1
	MatchStatusReady
	MatchStatusInProgress
	MatchStatusFinished
)

type Match struct {
	ID         string
	Mode       GameMode
	Players    []Player
	Status     MatchStatus
	InviteCode string
	CreatedAt  time.Time
	StartedAt  *time.Time
}

type Player struct {
	ID    string
	Ready bool
}

func (m *Match) AddPlayer(playerID string) bool {
	if m.IsFull() {
		return false
	}

	m.Players = append(m.Players, Player{
		ID:    playerID,
		Ready: false,
	})

	if m.IsFull() {
		m.Status = MatchStatusReady
	}

	return true
}

func (m *Match) IsFull() bool {
	requiredPlayers := 1
	if m.Mode == GameModeOneVsOne {
		requiredPlayers = 2
	}
	return len(m.Players) >= requiredPlayers
}

func (m *Match) HasPlayer(playerID string) bool {
	for _, p := range m.Players {
		if p.ID == playerID {
			return true
		}
	}
	return false
}
