package handler

import (
	pb "GoTetrisOnline/api/proto/game/v1"
	"context"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/coder/websocket"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

type GatewayHandler struct {
	grpcClient pb.GameServiceClient
}

func NewGatewayHandler(conn *grpc.ClientConn) *GatewayHandler {
	return &GatewayHandler{
		grpcClient: pb.NewGameServiceClient(conn),
	}
}

func (h *GatewayHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
		OriginPatterns:     []string{"*"},
		CompressionMode:    websocket.CompressionDisabled,
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

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		for {
			msg, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					return nil
				}
				return err
			}

			data, err := proto.Marshal(msg)
			if err != nil {
				log.Printf("failed to marshal message: %v", err)
				continue
			}

			writeCtx, cancel := context.WithTimeout(ctx, time.Second*5)
			err = c.Write(writeCtx, websocket.MessageBinary, data)
			cancel()
			if err != nil {
				return err
			}
		}
	})

	g.Go(func() error {
		for {
			msgType, data, err := c.Read(ctx)
			if err != nil {
				if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
					return nil
				}
				return err
			}

			if msgType != websocket.MessageBinary {
				log.Printf("unexpected message type: %v", msgType)
				continue
			}

			var clientMsg pb.ClientMessage
			if err := proto.Unmarshal(data, &clientMsg); err != nil {
				log.Printf("failed to unmarshal client message: %v", err)
				continue
			}

			if err := stream.Send(&clientMsg); err != nil {
				return err
			}
		}
	})

	if err := g.Wait(); err != nil {
		log.Printf("session closed: %v", err)
	}

	c.Close(websocket.StatusNormalClosure, "bye")
}
