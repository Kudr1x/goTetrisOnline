package storage

import (
	"GoTetrisOnline/services/matchmaking/domain"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"sync"
	"time"
)

var (
	ErrMatchNotFound   = errors.New("match not found")
	ErrInviteNotFound  = errors.New("invite code not found")
	ErrInviteExpired   = errors.New("invite code expired")
	ErrMatchFull       = errors.New("match is full")
)

type InMemoryStorage struct {
	mu            sync.RWMutex
	matches       map[string]*domain.Match
	inviteCodes   map[string]string
	inviteExpiry  map[string]time.Time
	waitingQueues map[domain.GameMode][]string
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		matches:       make(map[string]*domain.Match),
		inviteCodes:   make(map[string]string),
		inviteExpiry:  make(map[string]time.Time),
		waitingQueues: make(map[domain.GameMode][]string),
	}
}

func (s *InMemoryStorage) CreateMatch(mode domain.GameMode, playerID string) (*domain.Match, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	matchID := generateID()
	inviteCode := generateInviteCode()

	match := &domain.Match{
		ID:         matchID,
		Mode:       mode,
		Players:    []domain.Player{{ID: playerID, Ready: false}},
		Status:     domain.MatchStatusWaiting,
		InviteCode: inviteCode,
		CreatedAt:  time.Now(),
	}

	if mode == domain.GameModeSolo {
		match.Status = domain.MatchStatusReady
	}

	s.matches[matchID] = match
	s.inviteCodes[inviteCode] = matchID
	s.inviteExpiry[inviteCode] = time.Now().Add(5 * time.Minute)

	return match, nil
}

func (s *InMemoryStorage) GetMatch(matchID string) (*domain.Match, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	match, ok := s.matches[matchID]
	if !ok {
		return nil, ErrMatchNotFound
	}

	return match, nil
}

func (s *InMemoryStorage) JoinByInvite(inviteCode, playerID string) (*domain.Match, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	matchID, ok := s.inviteCodes[inviteCode]
	if !ok {
		return nil, ErrInviteNotFound
	}

	if time.Now().After(s.inviteExpiry[inviteCode]) {
		delete(s.inviteCodes, inviteCode)
		delete(s.inviteExpiry, inviteCode)
		return nil, ErrInviteExpired
	}

	match, ok := s.matches[matchID]
	if !ok {
		return nil, ErrMatchNotFound
	}

	if match.IsFull() {
		return nil, ErrMatchFull
	}

	if !match.AddPlayer(playerID) {
		return nil, ErrMatchFull
	}

	return match, nil
}

func (s *InMemoryStorage) FindRandomOpponent(mode domain.GameMode, playerID string) (*domain.Match, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	queue := s.waitingQueues[mode]

	if len(queue) > 0 {
		waitingPlayerID := queue[0]
		s.waitingQueues[mode] = queue[1:]

		for matchID, match := range s.matches {
			if match.Mode == mode && match.Status == domain.MatchStatusWaiting && match.HasPlayer(waitingPlayerID) {
				if match.AddPlayer(playerID) {
					return s.matches[matchID], true
				}
			}
		}
	}

	matchID := generateID()
	inviteCode := generateInviteCode()

	match := &domain.Match{
		ID:         matchID,
		Mode:       mode,
		Players:    []domain.Player{{ID: playerID, Ready: false}},
		Status:     domain.MatchStatusWaiting,
		InviteCode: inviteCode,
		CreatedAt:  time.Now(),
	}

	s.matches[matchID] = match
	s.inviteCodes[inviteCode] = matchID
	s.inviteExpiry[inviteCode] = time.Now().Add(5 * time.Minute)

	s.waitingQueues[mode] = append(s.waitingQueues[mode], playerID)

	return match, false
}

func (s *InMemoryStorage) UpdateMatchStatus(matchID string, status domain.MatchStatus) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	match, ok := s.matches[matchID]
	if !ok {
		return ErrMatchNotFound
	}

	match.Status = status
	if status == domain.MatchStatusInProgress && match.StartedAt == nil {
		now := time.Now()
		match.StartedAt = &now
	}

	return nil
}

func generateID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)[:16]
}

func generateInviteCode() string {
	b := make([]byte, 6)
	_, _ = rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)[:8]
}
