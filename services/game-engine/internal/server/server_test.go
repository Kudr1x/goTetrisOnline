package server

import (
	"GoTetrisOnline/pkg/core"
	"GoTetrisOnline/services/game-engine/domain"
	"testing"

	pb "GoTetrisOnline/api/proto/game/v1"
)

func TestMapEventToProto_StateUpdate_IncludesNextPieces(t *testing.T) {
	nextPieces := []core.PieceType{core.PieceI, core.PieceO, core.PieceT}

	stateDTO := domain.GameStateDTO{
		Score: 1000,
		Level: 5,
		Grid:  make([]byte, core.BoardWidth*core.BoardHeight),
		CurrentPiece: core.Piece{
			Type:     core.PieceJ,
			Position: core.Point{X: 4, Y: 10},
			Rotation: 1,
		},
		NextPieces: nextPieces,
	}

	event := domain.GameEvent{
		Type:    "state_update",
		Payload: stateDTO,
	}

	protoMsg := mapEventToProto(event)

	if protoMsg == nil {
		t.Fatal("mapEventToProto returned nil")
	}

	stateUpdate, ok := protoMsg.Payload.(*pb.ServerMessage_State)
	if !ok {
		t.Fatalf("Expected ServerMessage_State, got %T", protoMsg.Payload)
	}

	state := stateUpdate.State

	if state.Score != 1000 {
		t.Errorf("Expected Score=1000, got %d", state.Score)
	}

	if state.Level != 5 {
		t.Errorf("Expected Level=5, got %d", state.Level)
	}

	if len(state.NextPieces) != 3 {
		t.Fatalf("Expected 3 NextPieces, got %d", len(state.NextPieces))
	}

	expectedTypes := []pb.PieceType{pb.PieceType_PIECE_I, pb.PieceType_PIECE_O, pb.PieceType_PIECE_T}
	for i, expected := range expectedTypes {
		if state.NextPieces[i] != expected {
			t.Errorf("NextPieces[%d]: expected %v, got %v", i, expected, state.NextPieces[i])
		}
	}

	if state.CurrentPiece.Type != pb.PieceType_PIECE_J {
		t.Errorf("Expected CurrentPiece type J, got %v", state.CurrentPiece.Type)
	}

	if state.CurrentPiece.X != 4 {
		t.Errorf("Expected CurrentPiece X=4, got %d", state.CurrentPiece.X)
	}

	if state.CurrentPiece.Y != 10 {
		t.Errorf("Expected CurrentPiece Y=10, got %d", state.CurrentPiece.Y)
	}

	if state.CurrentPiece.Rotation != 1 {
		t.Errorf("Expected CurrentPiece Rotation=1, got %d", state.CurrentPiece.Rotation)
	}
}

func TestMapEventToProto_StateUpdate_EmptyNextPieces(t *testing.T) {
	stateDTO := domain.GameStateDTO{
		Score:        0,
		Level:        1,
		Grid:         make([]byte, core.BoardWidth*core.BoardHeight),
		CurrentPiece: core.Piece{Type: core.PieceT},
		NextPieces:   []core.PieceType{},
	}

	event := domain.GameEvent{
		Type:    "state_update",
		Payload: stateDTO,
	}

	protoMsg := mapEventToProto(event)

	if protoMsg == nil {
		t.Fatal("mapEventToProto returned nil")
	}

	stateUpdate := protoMsg.Payload.(*pb.ServerMessage_State)
	state := stateUpdate.State

	if len(state.NextPieces) != 0 {
		t.Errorf("Expected 0 NextPieces, got %d", len(state.NextPieces))
	}
}

func TestMapEventToProto_GameOver(t *testing.T) {
	event := domain.GameEvent{
		Type:    "game_over",
		Payload: int32(5000),
	}

	protoMsg := mapEventToProto(event)

	if protoMsg == nil {
		t.Fatal("mapEventToProto returned nil")
	}

	gameEvent, ok := protoMsg.Payload.(*pb.ServerMessage_Event)
	if !ok {
		t.Fatalf("Expected ServerMessage_Event, got %T", protoMsg.Payload)
	}

	if gameEvent.Event.Type != pb.EventType_EVENT_GAME_OVER {
		t.Errorf("Expected EVENT_GAME_OVER, got %v", gameEvent.Event.Type)
	}

	if gameEvent.Event.Message != "Game Over" {
		t.Errorf("Expected 'Game Over', got '%s'", gameEvent.Event.Message)
	}
}

func TestMapEventToProto_UnknownEvent(t *testing.T) {
	event := domain.GameEvent{
		Type:    "unknown_event",
		Payload: nil,
	}

	protoMsg := mapEventToProto(event)

	if protoMsg != nil {
		t.Errorf("Expected nil for unknown event, got %+v", protoMsg)
	}
}

func TestMapEventToProto_InvalidPayloadType(t *testing.T) {
	event := domain.GameEvent{
		Type:    "state_update",
		Payload: "invalid",
	}

	protoMsg := mapEventToProto(event)

	if protoMsg != nil {
		t.Errorf("Expected nil for invalid payload type, got %+v", protoMsg)
	}
}
