package handler

import (
	pb "GoTetrisOnline/api/proto/game/v1"
	"context"
	"log"
	"net/http"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

type GatewayHandler struct {
	grpcClient pb.GameServiceClient
}

func NewGatewayHandler(conn *grpc.ClientConn) *GatewayHandler {
	return &GatewayHandler{
		grpcClient: pb.NewGameServiceClient(conn),
	}
}

type BrowserCommand struct {
	Cmd string `json:"cmd"`
}

func (h *GatewayHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns: []string{"*"},
	})
	if err != nil {
		log.Printf("failed to accept websocket: %v", err)
		return
	}
	defer c.Close(websocket.StatusInternalError, "internal error")

	ctx := r.Context()

	stream, err := h.grpcClient.Play(ctx)
	if err != nil {
		log.Printf("failed to connect to game engine: %v", err)
		c.Close(websocket.StatusBadGateway, "game engine unavailable")
		return
	}

	err = stream.Send(&pb.ClientMessage{
		Payload: &pb.ClientMessage_Join{
			// todo
			Join: &pb.JoinRequest{
				MatchId: "room-1",
			},
		},
	})
	if err != nil {
		log.Printf("failed to join match: %v", err)
		return
	}

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		for {
			msg, err := stream.Recv()
			if err != nil {
				return err
			}

			writeCtx, cancel := context.WithTimeout(ctx, time.Second*5)
			err = wsjson.Write(writeCtx, c, msg)
			cancel()
			if err != nil {
				return err
			}
		}
	})

	g.Go(func() error {
		for {
			var cmd BrowserCommand
			err := wsjson.Read(ctx, c, &cmd)
			if err != nil {
				if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
					return nil
				}
				return err
			}

			input := mapJsonToProto(cmd.Cmd)
			if input != pb.InputType_INPUT_UNSPECIFIED {
				err = stream.Send(&pb.ClientMessage{
					Payload: &pb.ClientMessage_Input{
						Input: &pb.InputRequest{Input: input},
					},
				})
				if err != nil {
					return err
				}
			}
		}
	})

	if err := g.Wait(); err != nil {
		log.Printf("session closed: %v", err)
	}

	c.Close(websocket.StatusNormalClosure, "bye")
}

func mapJsonToProto(cmd string) pb.InputType {
	switch cmd {
	case "left":
		return pb.InputType_INPUT_LEFT
	case "right":
		return pb.InputType_INPUT_RIGHT
	case "rotate":
		return pb.InputType_INPUT_ROTATE_CW
	case "drop":
		return pb.InputType_INPUT_HARD_DROP
	default:
		return pb.InputType_INPUT_UNSPECIFIED
	}
}
