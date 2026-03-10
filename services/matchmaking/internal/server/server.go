package server

import (
	"context"
	"fmt"

	pb "GoTetrisOnline/api/proto/matchmaking/v1"
	"GoTetrisOnline/services/matchmaking/domain"
	"GoTetrisOnline/services/matchmaking/internal/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcServer struct {
	pb.UnimplementedMatchmakingServiceServer
	storage    *storage.InMemoryStorage
	inviteBase string
}

func NewGrpcServer(storage *storage.InMemoryStorage, inviteBase string) *GrpcServer {
	return &GrpcServer{
		storage:    storage,
		inviteBase: inviteBase,
	}
}

func (s *GrpcServer) CreateMatch(ctx context.Context, req *pb.CreateMatchRequest) (*pb.CreateMatchResponse, error) {
	if req.PlayerId == "" {
		return nil, status.Error(codes.InvalidArgument, "player_id is required")
	}

	mode := mapProtoToGameMode(req.Mode)
	match, err := s.storage.CreateMatch(mode, req.PlayerId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create match: %v", err)
	}

	return &pb.CreateMatchResponse{
		MatchId:    match.ID,
		InviteCode: match.InviteCode,
		InviteLink: fmt.Sprintf("%s/join/%s", s.inviteBase, match.InviteCode),
	}, nil
}

func (s *GrpcServer) JoinByInvite(ctx context.Context, req *pb.JoinByInviteRequest) (*pb.JoinByInviteResponse, error) {
	if req.InviteCode == "" {
		return nil, status.Error(codes.InvalidArgument, "invite_code is required")
	}
	if req.PlayerId == "" {
		return nil, status.Error(codes.InvalidArgument, "player_id is required")
	}

	match, err := s.storage.JoinByInvite(req.InviteCode, req.PlayerId)
	if err != nil {
		switch err {
		case storage.ErrInviteNotFound:
			return nil, status.Error(codes.NotFound, "invite code not found")
		case storage.ErrInviteExpired:
			return nil, status.Error(codes.DeadlineExceeded, "invite code expired")
		case storage.ErrMatchFull:
			return nil, status.Error(codes.ResourceExhausted, "match is full")
		default:
			return nil, status.Errorf(codes.Internal, "failed to join match: %v", err)
		}
	}

	return &pb.JoinByInviteResponse{
		MatchId: match.ID,
		Status:  mapMatchStatusToProto(match.Status),
	}, nil
}

func (s *GrpcServer) FindRandom(ctx context.Context, req *pb.FindRandomRequest) (*pb.FindRandomResponse, error) {
	if req.PlayerId == "" {
		return nil, status.Error(codes.InvalidArgument, "player_id is required")
	}

	mode := mapProtoToGameMode(req.Mode)
	match, _ := s.storage.FindRandomOpponent(mode, req.PlayerId)

	return &pb.FindRandomResponse{
		MatchId: match.ID,
		Status:  mapMatchStatusToProto(match.Status),
	}, nil
}

func (s *GrpcServer) GetMatchInfo(ctx context.Context, req *pb.GetMatchInfoRequest) (*pb.GetMatchInfoResponse, error) {
	if req.MatchId == "" {
		return nil, status.Error(codes.InvalidArgument, "match_id is required")
	}

	match, err := s.storage.GetMatch(req.MatchId)
	if err != nil {
		if err == storage.ErrMatchNotFound {
			return nil, status.Error(codes.NotFound, "match not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get match: %v", err)
	}

	players := make([]*pb.Player, len(match.Players))
	for i, p := range match.Players {
		players[i] = &pb.Player{
			PlayerId: p.ID,
			Ready:    p.Ready,
		}
	}

	return &pb.GetMatchInfoResponse{
		MatchId: match.ID,
		Mode:    mapGameModeToProto(match.Mode),
		Players: players,
		Status:  mapMatchStatusToProto(match.Status),
	}, nil
}

func mapProtoToGameMode(mode pb.GameMode) domain.GameMode {
	switch mode {
	case pb.GameMode_GAME_MODE_SOLO:
		return domain.GameModeSolo
	case pb.GameMode_GAME_MODE_ONE_VS_ONE:
		return domain.GameModeOneVsOne
	default:
		return domain.GameModeSolo
	}
}

func mapGameModeToProto(mode domain.GameMode) pb.GameMode {
	switch mode {
	case domain.GameModeSolo:
		return pb.GameMode_GAME_MODE_SOLO
	case domain.GameModeOneVsOne:
		return pb.GameMode_GAME_MODE_ONE_VS_ONE
	default:
		return pb.GameMode_GAME_MODE_UNSPECIFIED
	}
}

func mapMatchStatusToProto(status domain.MatchStatus) pb.MatchStatus {
	switch status {
	case domain.MatchStatusWaiting:
		return pb.MatchStatus_MATCH_STATUS_WAITING
	case domain.MatchStatusReady:
		return pb.MatchStatus_MATCH_STATUS_READY
	case domain.MatchStatusInProgress:
		return pb.MatchStatus_MATCH_STATUS_IN_PROGRESS
	case domain.MatchStatusFinished:
		return pb.MatchStatus_MATCH_STATUS_FINISHED
	default:
		return pb.MatchStatus_MATCH_STATUS_UNSPECIFIED
	}
}
