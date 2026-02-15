package server

import (
	"GoTetrisOnline/services/game-engine/domain"

	pb "GoTetrisOnline/api/proto/game/v1"

	"io"
	"log"

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
	ctx := stream.Context()

	req, err := stream.Recv()
	if err != nil {
		return err
	}

	joinReq, ok := req.Payload.(*pb.ClientMessage_Join)
	if !ok {
		return status.Error(codes.InvalidArgument, "first message must be JoinRequest")
	}

	// todo
	playerID := "player-1"
	log.Printf("Player %s joined match %s", playerID, joinReq.Join.MatchId)

	game := domain.NewGame(joinReq.Join.MatchId)
	game.Start()
	defer game.Stop()

	writeErrChan := make(chan error, 1)
	go func() {
		defer close(writeErrChan)

		for {
			select {
			case <-ctx.Done():
				return
			case event, ok := <-game.Events():
				if !ok {
					return
				}

				protoMsg := mapEventToProto(event)

				if err := stream.Send(protoMsg); err != nil {
					writeErrChan <- err
					return
				}
			}
		}
	}()

	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			select {
			case wErr := <-writeErrChan:
				return wErr
			default:
				return err
			}
		}

		if input := in.GetInput(); input != nil {
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
	}
}

func mapEventToProto(event domain.GameEvent) *pb.ServerMessage {
	switch event.Type {
	case "state_update":
		state, ok := event.Payload.(domain.GameStateDTO)
		if !ok {
			return nil
		}

		return &pb.ServerMessage{
			Payload: &pb.ServerMessage_State{
				State: &pb.StateUpdate{
					Score: state.Score,
					Level: state.Level,
					Grid:  state.Grid,

					CurrentPiece: &pb.Piece{
						Type:     pb.PieceType(state.CurrentPiece.Type),
						X:        int32(state.CurrentPiece.Position.X),
						Y:        int32(state.CurrentPiece.Position.Y),
						Rotation: int32(state.CurrentPiece.Rotation),
					},
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
