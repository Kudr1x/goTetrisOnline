package main

import (
	pb "GoTetrisOnline/api/proto/game/v1"
	"GoTetrisOnline/pkg/core"
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewGameServiceClient(conn)

	stream, err := client.Play(context.Background())
	if err != nil {
		log.Printf("error creating stream: %v", err)
		return
	}

	if err := stream.Send(&pb.ClientMessage{
		Payload: &pb.ClientMessage_Join{
			Join: &pb.JoinRequest{
				MatchId: "room-1",
				Token:   "token",
			},
		},
	}); err != nil {
		log.Printf("failed to send join: %v", err)
		return
	}

	go func() {
		for {
			msg, err := stream.Recv()
			if err != nil {
				log.Fatalf("Stream closed: %v", err)
			}

			switch payload := msg.Payload.(type) {
			case *pb.ServerMessage_State:
				render(payload.State)
			case *pb.ServerMessage_Event:
				fmt.Printf("EVENT: %s", payload.Event.Message)
			}
		}
	}()

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Controls: a (left), d (right), w (rotate), s (drop), q (quit)")

	for {
		char, _ := reader.ReadString('\n')
		char = strings.TrimSpace(char)

		var input pb.InputType

		switch char {
		case "a":
			input = pb.InputType_INPUT_LEFT
		case "d":
			input = pb.InputType_INPUT_RIGHT
		case "w":
			input = pb.InputType_INPUT_ROTATE_CW
		case "s":
			input = pb.InputType_INPUT_HARD_DROP
		case "q":
			return
		default:
			continue
		}

		if err := stream.Send(&pb.ClientMessage{
			Payload: &pb.ClientMessage_Input{
				Input: &pb.InputRequest{
					Input: input,
				},
			},
		}); err != nil {
			log.Printf("Failed to send input: %v", err)
			return
		}
	}
}

func render(state *pb.StateUpdate) {
	fmt.Print("\033[H\033[2J")

	currentPiece := core.Piece{
		Type: core.PieceType(state.CurrentPiece.Type), //nolint:gosec // coordinates are small
		Position: core.Point{
			X: int(state.CurrentPiece.X),
			Y: int(state.CurrentPiece.Y),
		},
		Rotation: int(state.CurrentPiece.Rotation),
	}

	minos := core.GetRotatedMinos(currentPiece.Type, currentPiece.Rotation)

	display := make([]string, core.BoardHeight)

	for y := 0; y < core.BoardHeight; y++ {
		row := ""

		for x := 0; x < core.BoardWidth; x++ {
			char := " . "

			idx := y*core.BoardWidth + x
			if idx < len(state.Grid) && state.Grid[idx] != 0 {
				char = "[#]"
			}

			for _, mino := range minos {
				absX := currentPiece.Position.X + mino.X
				absY := currentPiece.Position.Y + mino.Y

				if absX == x && absY == y {
					char = "[@]"
				}
			}
			row += char
		}
		display[y] = fmt.Sprintf("|%s|", row)
	}

	fmt.Println(strings.Join(display, "\n"))
	fmt.Printf("Score: %d | Level: %d\n", state.Score, state.Level)
}
