#!/bin/bash
echo "Starting Game Engine on :50051..."
go run services/game-engine/cmd/main.go &
ENGINE_PID=$!

sleep 2

echo "Starting Gateway on :8081..."
go run services/gateway/cmd/main.go &
GATEWAY_PID=$!

sleep 2

echo "Starting Web Server on :8080..."
go run client/cmd/web/main.go &
WEB_PID=$!

echo ""
echo "All services started:"
echo "  Game Engine: $ENGINE_PID"
echo "  Gateway: $GATEWAY_PID"
echo "  Web Server: $WEB_PID"
echo ""
echo "Open: http://localhost:8080"
echo ""
echo "Press Ctrl+C to stop all services"

trap "kill $ENGINE_PID $GATEWAY_PID $WEB_PID 2>/dev/null; exit" INT TERM

wait
