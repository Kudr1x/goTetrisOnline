package main

import (
	pb "GoTetrisOnline/api/proto/game/v1"
	"GoTetrisOnline/pkg/core"
	"context"
	"fmt"
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	colorI     = lipgloss.NewStyle().Foreground(lipgloss.Color("51")).Bold(true)  // Cyan
	colorO     = lipgloss.NewStyle().Foreground(lipgloss.Color("226")).Bold(true) // Yellow
	colorT     = lipgloss.NewStyle().Foreground(lipgloss.Color("129")).Bold(true) // Purple
	colorS     = lipgloss.NewStyle().Foreground(lipgloss.Color("46")).Bold(true)  // Green
	colorZ     = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true) // Red
	colorJ     = lipgloss.NewStyle().Foreground(lipgloss.Color("21")).Bold(true)  // Blue
	colorL     = lipgloss.NewStyle().Foreground(lipgloss.Color("208")).Bold(true) // Orange
	colorGray  = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))            // Dark gray
	colorPiece = lipgloss.NewStyle().Foreground(lipgloss.Color("255")).Bold(true) // White (current piece)

	boardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(0, 1)

	sidebarStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(0, 1)
)

type model struct {
	stream     pb.GameService_PlayClient
	state      *pb.StateUpdate
	gameOver   bool
	finalScore int32
	err        error
	width      int
	height     int
}

type gameStateMsg struct {
	state *pb.StateUpdate
}

type gameOverMsg struct {
	score int32
}

type errMsg struct {
	err error
}

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

	m := model{
		stream: stream,
	}

	p := tea.NewProgram(&m, tea.WithAltScreen())

	go receiveLoop(stream, p)

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func receiveLoop(stream pb.GameService_PlayClient, p *tea.Program) {
	for {
		msg, err := stream.Recv()
		if err != nil {
			p.Send(errMsg{err: err})
			return
		}

		switch payload := msg.Payload.(type) {
		case *pb.ServerMessage_State:
			p.Send(gameStateMsg{state: payload.State})
		case *pb.ServerMessage_Event:
			if payload.Event.Type == pb.EventType_EVENT_GAME_OVER {
				p.Send(gameOverMsg{score: 0})
			}
		}
	}
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.gameOver {
			if msg.String() == "q" || msg.String() == "ctrl+c" {
				return m, tea.Quit
			}
			return m, nil
		}

		var input pb.InputType
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "a", "left":
			input = pb.InputType_INPUT_LEFT
		case "d", "right":
			input = pb.InputType_INPUT_RIGHT
		case "w", "up":
			input = pb.InputType_INPUT_ROTATE_CW
		case "e":
			input = pb.InputType_INPUT_ROTATE_CCW
		case "s", "down":
			input = pb.InputType_INPUT_SOFT_DROP
		case " ":
			input = pb.InputType_INPUT_HARD_DROP
		default:
			return m, nil
		}

		if err := m.stream.Send(&pb.ClientMessage{
			Payload: &pb.ClientMessage_Input{
				Input: &pb.InputRequest{
					Input: input,
				},
			},
		}); err != nil {
			m.err = err
			return m, tea.Quit
		}

	case gameStateMsg:
		m.state = msg.state

	case gameOverMsg:
		m.gameOver = true
		if m.state != nil {
			m.finalScore = m.state.Score
		}

	case errMsg:
		m.err = msg.err
		return m, tea.Quit

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

func (m *model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n", m.err)
	}

	if m.gameOver {
		return fmt.Sprintf("\nGAME OVER!\n\nFinal Score: %d\n\nPress 'q' to quit\n", m.finalScore)
	}

	if m.state == nil {
		return "Connecting to game server...\n"
	}

	return renderGame(m.state)
}

func renderGame(state *pb.StateUpdate) string {
	boardContent := renderBoard(state)
	sidebarContent := renderSidebar(state)

	board := boardStyle.Render(boardContent)
	sidebar := sidebarStyle.Render(sidebarContent)

	return lipgloss.JoinHorizontal(lipgloss.Top, board, sidebar)
}

func renderBoard(state *pb.StateUpdate) string {
	var b strings.Builder

	currentPiece := core.Piece{
		Type: core.PieceType(state.CurrentPiece.Type), //nolint:gosec // piece types are small enums
		Position: core.Point{
			X: int(state.CurrentPiece.X),
			Y: int(state.CurrentPiece.Y),
		},
		Rotation: int(state.CurrentPiece.Rotation),
	}

	minos := core.GetRotatedMinos(currentPiece.Type, currentPiece.Rotation)

	for y := core.Space; y < core.BoardHeight; y++ {
		for x := 0; x < core.BoardWidth; x++ {
			var char string
			cellType := core.PieceNone

			idx := y*core.BoardWidth + x
			if idx < len(state.Grid) && state.Grid[idx] != 0 {
				cellType = core.PieceType(state.Grid[idx]) //nolint:gosec // piece types are small enums
			}

			isPiece := false
			for _, mino := range minos {
				absX := currentPiece.Position.X + mino.X
				absY := currentPiece.Position.Y + mino.Y

				if absX == x && absY == y {
					isPiece = true
					break
				}
			}

			switch {
			case isPiece:
				char = colorPiece.Render("██")
			case cellType != core.PieceNone:
				char = getColorForPiece(cellType).Render("██")
			default:
				char = colorGray.Render("░░")
			}

			b.WriteString(char)
		}

		if y < core.BoardHeight-1 {
			b.WriteString("\n")
		}
	}

	return b.String()
}

func renderSidebar(state *pb.StateUpdate) string {
	var b strings.Builder

	b.WriteString("NEXT:\n\n")

	if len(state.NextPieces) > 0 {
		nextType := core.PieceType(state.NextPieces[0]) //nolint:gosec // piece types are small enums
		b.WriteString(renderNextPiece(nextType))
	}

	b.WriteString("\n\n")
	b.WriteString(fmt.Sprintf("Score: %d\n", state.Score))
	b.WriteString(fmt.Sprintf("Level: %d\n", state.Level))
	b.WriteString("\n")
	b.WriteString("CONTROLS\n")
	b.WriteString("A: Left\n")
	b.WriteString("D: Right\n")
	b.WriteString("W: Rotate CW\n")
	b.WriteString("E: Rotate CCW\n")
	b.WriteString("S: Soft Drop\n")
	b.WriteString("Space: Drop!\n")
	b.WriteString("Q: Quit")

	return b.String()
}

func renderNextPiece(t core.PieceType) string {
	var b strings.Builder
	minos := core.GetRotatedMinos(t, 0)

	minX, maxX := 0, 0
	minY, maxY := 0, 0
	for i, m := range minos {
		if i == 0 {
			minX, maxX = m.X, m.X
			minY, maxY = m.Y, m.Y
		} else {
			if m.X < minX {
				minX = m.X
			}
			if m.X > maxX {
				maxX = m.X
			}
			if m.Y < minY {
				minY = m.Y
			}
			if m.Y > maxY {
				maxY = m.Y
			}
		}
	}

	width := maxX - minX + 1
	height := maxY - minY + 1
	gridSize := 4
	offsetX := (gridSize - width) / 2
	offsetY := (gridSize - height) / 2

	grid := make([][]bool, gridSize)
	for i := range grid {
		grid[i] = make([]bool, gridSize)
	}

	for _, m := range minos {
		x := m.X - minX + offsetX
		y := m.Y - minY + offsetY
		if x >= 0 && x < gridSize && y >= 0 && y < gridSize {
			grid[y][x] = true
		}
	}

	style := getColorForPiece(t)
	for y := 0; y < gridSize; y++ {
		for x := 0; x < gridSize; x++ {
			if grid[y][x] {
				b.WriteString(style.Render("██"))
			} else {
				b.WriteString(colorGray.Render("░░"))
			}
		}
		if y < gridSize-1 {
			b.WriteString("\n")
		}
	}

	return b.String()
}

func getColorForPiece(t core.PieceType) lipgloss.Style {
	switch t {
	case core.PieceI:
		return colorI
	case core.PieceO:
		return colorO
	case core.PieceT:
		return colorT
	case core.PieceS:
		return colorS
	case core.PieceZ:
		return colorZ
	case core.PieceJ:
		return colorJ
	case core.PieceL:
		return colorL
	default:
		return colorGray
	}
}
