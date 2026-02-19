package server

import (
	"GoTetrisOnline/services/game-engine/domain"
	"errors"

	pb "GoTetrisOnline/api/proto/game/v1"

	"io"
	"log"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcServer struct {
	pb.UnimplementedGameServiceServer
}

func NewGrpcServer() *GrpcServer {
	return &GrpcServer{}
}

func (s *GrpcServer) Play(stream pb.GameService_PlayServer) error {
	req, err := stream.Recv()
	if err != nil {
		return err
	}

	joinReq, ok := req.Payload.(*pb.ClientMessage_Join)
	if !ok {
		return status.Error(codes.InvalidArgument, "first message must be JoinRequest")
	}

	matchID := joinReq.Join.MatchId
	log.Printf("Player joining match %s", matchID)

	game := domain.NewGame(matchID)
	game.Start()

	g, ctx := errgroup.WithContext(stream.Context())

	g.Go(func() error {
		defer game.Stop()

		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case event, ok := <-game.Events():
				if !ok {
					return nil
				}

				protoMsg := mapEventToProto(event)
				if protoMsg == nil {
					continue
				}

				if err := stream.Send(protoMsg); err != nil {
					return err
				}
			}
		}
	})

	g.Go(func() error {
		for {
			in, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				return nil
			}
			if err != nil {
				return err
			}

			handleInput(game, in.GetInput())
		}
	})

	return g.Wait()
}

func handleInput(game *domain.Game, input *pb.InputRequest) {
	if input == nil {
		return
	}
	switch input.Input {
	case pb.InputType_INPUT_LEFT:
		game.MoveLeft()
	case pb.InputType_INPUT_RIGHT:
		game.MoveRight()
	case pb.InputType_INPUT_ROTATE_CW:
		game.Rotate(1)
	case pb.InputType_INPUT_ROTATE_CCW:
		game.Rotate(-1)
	case pb.InputType_INPUT_HARD_DROP:
		game.HardDrop()
	}
}
func mapEventToProto(event domain.GameEvent) *pb.ServerMessage {
	switch event.Type {
	case "state_update":
		state, ok := event.Payload.(domain.GameStateDTO)
		if !ok {
			return nil
		}

		nextPieces := make([]pb.PieceType, len(state.NextPieces))
		for i, pieceType := range state.NextPieces {
			nextPieces[i] = pb.PieceType(pieceType) //nolint:gosec // piece types are small enums
		}

		return &pb.ServerMessage{
			Payload: &pb.ServerMessage_State{
				State: &pb.StateUpdate{
					Score: state.Score,
					Level: state.Level,
					Grid:  state.Grid,

					CurrentPiece: &pb.Piece{
						Type:     pb.PieceType(state.CurrentPiece.Type), //nolint:gosec // coordinates are small
						X:        int32(state.CurrentPiece.Position.X),  //nolint:gosec // coordinates are small
						Y:        int32(state.CurrentPiece.Position.Y),  //nolint:gosec // coordinates are small
						Rotation: int32(state.CurrentPiece.Rotation),    //nolint:gosec // coordinates are small
					},
					NextPieces: nextPieces,
				},
			},
		}

	case "game_over":
		return &pb.ServerMessage{
			Payload: &pb.ServerMessage_Event{
				Event: &pb.GameEvent{
					Type:    pb.EventType_EVENT_GAME_OVER,
					Message: "Game Over",
				},
			},
		}
	}
	return nil
}
