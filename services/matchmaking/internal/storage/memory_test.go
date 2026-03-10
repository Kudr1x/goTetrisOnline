package storage

import (
	"GoTetrisOnline/services/matchmaking/domain"
	"testing"
	"time"
)

func TestCreateMatch(t *testing.T) {
	storage := NewInMemoryStorage()

	match, err := storage.CreateMatch(domain.GameModeSolo, "player1")
	if err != nil {
		t.Fatalf("CreateMatch failed: %v", err)
	}

	if match.ID == "" {
		t.Error("Match ID should not be empty")
	}

	if match.Mode != domain.GameModeSolo {
		t.Errorf("Expected mode %v, got %v", domain.GameModeSolo, match.Mode)
	}

	if len(match.Players) != 1 {
		t.Errorf("Expected 1 player, got %d", len(match.Players))
	}

	if match.Status != domain.MatchStatusReady {
		t.Errorf("Solo match should be ready immediately, got %v", match.Status)
	}
}

func TestCreateOneVsOneMatch(t *testing.T) {
	storage := NewInMemoryStorage()

	match, err := storage.CreateMatch(domain.GameModeOneVsOne, "player1")
	if err != nil {
		t.Fatalf("CreateMatch failed: %v", err)
	}

	if match.Status != domain.MatchStatusWaiting {
		t.Errorf("1v1 match should be waiting, got %v", match.Status)
	}

	if match.InviteCode == "" {
		t.Error("Invite code should not be empty")
	}
}

func TestJoinByInvite(t *testing.T) {
	storage := NewInMemoryStorage()

	match1, _ := storage.CreateMatch(domain.GameModeOneVsOne, "player1")
	inviteCode := match1.InviteCode

	match2, err := storage.JoinByInvite(inviteCode, "player2")
	if err != nil {
		t.Fatalf("JoinByInvite failed: %v", err)
	}

	if match2.ID != match1.ID {
		t.Error("Should join the same match")
	}

	if len(match2.Players) != 2 {
		t.Errorf("Expected 2 players, got %d", len(match2.Players))
	}

	if match2.Status != domain.MatchStatusReady {
		t.Errorf("Match should be ready with 2 players, got %v", match2.Status)
	}
}

func TestJoinByInvalidInvite(t *testing.T) {
	storage := NewInMemoryStorage()

	_, err := storage.JoinByInvite("invalid", "player1")
	if err != ErrInviteNotFound {
		t.Errorf("Expected ErrInviteNotFound, got %v", err)
	}
}

func TestJoinFullMatch(t *testing.T) {
	storage := NewInMemoryStorage()

	match, _ := storage.CreateMatch(domain.GameModeOneVsOne, "player1")
	_, err := storage.JoinByInvite(match.InviteCode, "player2")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	_, err = storage.JoinByInvite(match.InviteCode, "player3")
	if err != ErrMatchFull {
		t.Errorf("Expected ErrMatchFull, got %v", err)
	}
}

func TestFindRandomOpponent(t *testing.T) {
	storage := NewInMemoryStorage()

	match1, found1 := storage.FindRandomOpponent(domain.GameModeOneVsOne, "player1")
	if found1 {
		t.Error("First player should not find opponent immediately")
	}

	if match1.Status != domain.MatchStatusWaiting {
		t.Errorf("Match should be waiting, got %v", match1.Status)
	}

	match2, found2 := storage.FindRandomOpponent(domain.GameModeOneVsOne, "player2")
	if !found2 {
		t.Error("Second player should find the first player")
	}

	if match2.ID != match1.ID {
		t.Error("Players should be matched together")
	}

	if len(match2.Players) != 2 {
		t.Errorf("Expected 2 players, got %d", len(match2.Players))
	}

	if match2.Status != domain.MatchStatusReady {
		t.Errorf("Match should be ready, got %v", match2.Status)
	}
}

func TestGetMatch(t *testing.T) {
	storage := NewInMemoryStorage()

	created, _ := storage.CreateMatch(domain.GameModeSolo, "player1")

	retrieved, err := storage.GetMatch(created.ID)
	if err != nil {
		t.Fatalf("GetMatch failed: %v", err)
	}

	if retrieved.ID != created.ID {
		t.Error("Retrieved match ID doesn't match")
	}
}

func TestGetNonExistentMatch(t *testing.T) {
	storage := NewInMemoryStorage()

	_, err := storage.GetMatch("nonexistent")
	if err != ErrMatchNotFound {
		t.Errorf("Expected ErrMatchNotFound, got %v", err)
	}
}

func TestUpdateMatchStatus(t *testing.T) {
	storage := NewInMemoryStorage()

	match, _ := storage.CreateMatch(domain.GameModeSolo, "player1")

	err := storage.UpdateMatchStatus(match.ID, domain.MatchStatusInProgress)
	if err != nil {
		t.Fatalf("UpdateMatchStatus failed: %v", err)
	}

	updated, _ := storage.GetMatch(match.ID)
	if updated.Status != domain.MatchStatusInProgress {
		t.Errorf("Expected status %v, got %v", domain.MatchStatusInProgress, updated.Status)
	}

	if updated.StartedAt == nil {
		t.Error("StartedAt should be set when status changes to InProgress")
	}
}

func TestInviteCodeExpiry(t *testing.T) {
	storage := NewInMemoryStorage()

	match, _ := storage.CreateMatch(domain.GameModeOneVsOne, "player1")
	inviteCode := match.InviteCode

	storage.inviteExpiry[inviteCode] = time.Now().Add(-1 * time.Second)

	_, err := storage.JoinByInvite(inviteCode, "player2")
	if err != ErrInviteExpired {
		t.Errorf("Expected ErrInviteExpired, got %v", err)
	}
}
