package domain

import "testing"

func TestMatch_AddPlayer(t *testing.T) {
	match := &Match{
		Mode:    GameModeOneVsOne,
		Status:  MatchStatusWaiting,
		Players: []Player{{ID: "player1"}},
	}

	ok := match.AddPlayer("player2")
	if !ok {
		t.Error("Should be able to add second player")
	}

	if len(match.Players) != 2 {
		t.Errorf("Expected 2 players, got %d", len(match.Players))
	}

	if match.Status != MatchStatusReady {
		t.Errorf("Match should be ready after adding second player, got %v", match.Status)
	}
}

func TestMatch_AddPlayerToFullMatch(t *testing.T) {
	match := &Match{
		Mode:   GameModeOneVsOne,
		Status: MatchStatusReady,
		Players: []Player{
			{ID: "player1"},
			{ID: "player2"},
		},
	}

	ok := match.AddPlayer("player3")
	if ok {
		t.Error("Should not be able to add player to full match")
	}

	if len(match.Players) != 2 {
		t.Errorf("Expected 2 players, got %d", len(match.Players))
	}
}

func TestMatch_IsFull_Solo(t *testing.T) {
	match := &Match{
		Mode:    GameModeSolo,
		Players: []Player{{ID: "player1"}},
	}

	if !match.IsFull() {
		t.Error("Solo match with 1 player should be full")
	}
}

func TestMatch_IsFull_OneVsOne(t *testing.T) {
	match := &Match{
		Mode:    GameModeOneVsOne,
		Players: []Player{{ID: "player1"}},
	}

	if match.IsFull() {
		t.Error("1v1 match with 1 player should not be full")
	}

	match.AddPlayer("player2")

	if !match.IsFull() {
		t.Error("1v1 match with 2 players should be full")
	}
}

func TestMatch_HasPlayer(t *testing.T) {
	match := &Match{
		Players: []Player{
			{ID: "player1"},
			{ID: "player2"},
		},
	}

	if !match.HasPlayer("player1") {
		t.Error("Match should have player1")
	}

	if !match.HasPlayer("player2") {
		t.Error("Match should have player2")
	}

	if match.HasPlayer("player3") {
		t.Error("Match should not have player3")
	}
}
