package server

import (
	"context"
	"testing"

	pb "GoTetrisOnline/api/proto/matchmaking/v1"
	"GoTetrisOnline/services/matchmaking/domain"
	"GoTetrisOnline/services/matchmaking/internal/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestCreateMatch_Success(t *testing.T) {
	store := storage.NewInMemoryStorage()
	server := NewGrpcServer(store, "http://localhost:8080")

	req := &pb.CreateMatchRequest{
		Mode:     pb.GameMode_GAME_MODE_SOLO,
		PlayerId: "player1",
	}

	resp, err := server.CreateMatch(context.Background(), req)
	if err != nil {
		t.Fatalf("CreateMatch failed: %v", err)
	}

	if resp.MatchId == "" {
		t.Error("MatchId should not be empty")
	}

	if resp.InviteCode == "" {
		t.Error("InviteCode should not be empty")
	}

	if resp.InviteLink == "" {
		t.Error("InviteLink should not be empty")
	}
}

func TestCreateMatch_EmptyPlayerId(t *testing.T) {
	store := storage.NewInMemoryStorage()
	server := NewGrpcServer(store, "http://localhost:8080")

	req := &pb.CreateMatchRequest{
		Mode:     pb.GameMode_GAME_MODE_SOLO,
		PlayerId: "",
	}

	_, err := server.CreateMatch(context.Background(), req)
	if err == nil {
		t.Error("Expected error for empty player_id")
	}

	st, ok := status.FromError(err)
	if !ok || st.Code() != codes.InvalidArgument {
		t.Errorf("Expected InvalidArgument, got %v", err)
	}
}

func TestJoinByInvite_Success(t *testing.T) {
	store := storage.NewInMemoryStorage()
	server := NewGrpcServer(store, "http://localhost:8080")

	createResp, _ := server.CreateMatch(context.Background(), &pb.CreateMatchRequest{
		Mode:     pb.GameMode_GAME_MODE_ONE_VS_ONE,
		PlayerId: "player1",
	})

	req := &pb.JoinByInviteRequest{
		InviteCode: createResp.InviteCode,
		PlayerId:   "player2",
	}

	resp, err := server.JoinByInvite(context.Background(), req)
	if err != nil {
		t.Fatalf("JoinByInvite failed: %v", err)
	}

	if resp.MatchId != createResp.MatchId {
		t.Error("MatchId mismatch")
	}

	if resp.Status != pb.MatchStatus_MATCH_STATUS_READY {
		t.Errorf("Expected READY status, got %v", resp.Status)
	}
}

func TestJoinByInvite_InvalidCode(t *testing.T) {
	store := storage.NewInMemoryStorage()
	server := NewGrpcServer(store, "http://localhost:8080")

	req := &pb.JoinByInviteRequest{
		InviteCode: "invalid",
		PlayerId:   "player1",
	}

	_, err := server.JoinByInvite(context.Background(), req)
	if err == nil {
		t.Error("Expected error for invalid invite code")
	}

	st, ok := status.FromError(err)
	if !ok || st.Code() != codes.NotFound {
		t.Errorf("Expected NotFound, got %v", err)
	}
}

func TestFindRandom_Success(t *testing.T) {
	store := storage.NewInMemoryStorage()
	server := NewGrpcServer(store, "http://localhost:8080")

	req1 := &pb.FindRandomRequest{
		Mode:     pb.GameMode_GAME_MODE_ONE_VS_ONE,
		PlayerId: "player1",
	}

	resp1, err := server.FindRandom(context.Background(), req1)
	if err != nil {
		t.Fatalf("FindRandom failed: %v", err)
	}

	if resp1.Status != pb.MatchStatus_MATCH_STATUS_WAITING {
		t.Errorf("First player should be waiting, got %v", resp1.Status)
	}

	req2 := &pb.FindRandomRequest{
		Mode:     pb.GameMode_GAME_MODE_ONE_VS_ONE,
		PlayerId: "player2",
	}

	resp2, err := server.FindRandom(context.Background(), req2)
	if err != nil {
		t.Fatalf("FindRandom failed: %v", err)
	}

	if resp2.MatchId != resp1.MatchId {
		t.Error("Players should be matched together")
	}

	if resp2.Status != pb.MatchStatus_MATCH_STATUS_READY {
		t.Errorf("Match should be ready, got %v", resp2.Status)
	}
}

func TestGetMatchInfo_Success(t *testing.T) {
	store := storage.NewInMemoryStorage()
	server := NewGrpcServer(store, "http://localhost:8080")

	createResp, _ := server.CreateMatch(context.Background(), &pb.CreateMatchRequest{
		Mode:     pb.GameMode_GAME_MODE_SOLO,
		PlayerId: "player1",
	})

	req := &pb.GetMatchInfoRequest{
		MatchId: createResp.MatchId,
	}

	resp, err := server.GetMatchInfo(context.Background(), req)
	if err != nil {
		t.Fatalf("GetMatchInfo failed: %v", err)
	}

	if resp.MatchId != createResp.MatchId {
		t.Error("MatchId mismatch")
	}

	if resp.Mode != pb.GameMode_GAME_MODE_SOLO {
		t.Errorf("Expected SOLO mode, got %v", resp.Mode)
	}

	if len(resp.Players) != 1 {
		t.Errorf("Expected 1 player, got %d", len(resp.Players))
	}
}

func TestGetMatchInfo_NotFound(t *testing.T) {
	store := storage.NewInMemoryStorage()
	server := NewGrpcServer(store, "http://localhost:8080")

	req := &pb.GetMatchInfoRequest{
		MatchId: "nonexistent",
	}

	_, err := server.GetMatchInfo(context.Background(), req)
	if err == nil {
		t.Error("Expected error for nonexistent match")
	}

	st, ok := status.FromError(err)
	if !ok || st.Code() != codes.NotFound {
		t.Errorf("Expected NotFound, got %v", err)
	}
}

func TestMapProtoToGameMode(t *testing.T) {
	tests := []struct {
		proto    pb.GameMode
		expected domain.GameMode
	}{
		{pb.GameMode_GAME_MODE_SOLO, domain.GameModeSolo},
		{pb.GameMode_GAME_MODE_ONE_VS_ONE, domain.GameModeOneVsOne},
		{pb.GameMode_GAME_MODE_UNSPECIFIED, domain.GameModeSolo},
	}

	for _, tt := range tests {
		result := mapProtoToGameMode(tt.proto)
		if result != tt.expected {
			t.Errorf("mapProtoToGameMode(%v) = %v, want %v", tt.proto, result, tt.expected)
		}
	}
}

func TestMapGameModeToProto(t *testing.T) {
	tests := []struct {
		mode     domain.GameMode
		expected pb.GameMode
	}{
		{domain.GameModeSolo, pb.GameMode_GAME_MODE_SOLO},
		{domain.GameModeOneVsOne, pb.GameMode_GAME_MODE_ONE_VS_ONE},
	}

	for _, tt := range tests {
		result := mapGameModeToProto(tt.mode)
		if result != tt.expected {
			t.Errorf("mapGameModeToProto(%v) = %v, want %v", tt.mode, result, tt.expected)
		}
	}
}

func TestMapMatchStatusToProto(t *testing.T) {
	tests := []struct {
		status   domain.MatchStatus
		expected pb.MatchStatus
	}{
		{domain.MatchStatusWaiting, pb.MatchStatus_MATCH_STATUS_WAITING},
		{domain.MatchStatusReady, pb.MatchStatus_MATCH_STATUS_READY},
		{domain.MatchStatusInProgress, pb.MatchStatus_MATCH_STATUS_IN_PROGRESS},
		{domain.MatchStatusFinished, pb.MatchStatus_MATCH_STATUS_FINISHED},
	}

	for _, tt := range tests {
		result := mapMatchStatusToProto(tt.status)
		if result != tt.expected {
			t.Errorf("mapMatchStatusToProto(%v) = %v, want %v", tt.status, result, tt.expected)
		}
	}
}
