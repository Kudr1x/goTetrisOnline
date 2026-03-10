.PHONY: gen clean test lint run-all run-engine run-matchmaking run-gateway run-terminal run-browser-server run-all-terminal run-all-browser wasm pt pb kill-ports

PROTO_PATH = api/proto

gen:
	protoc -I $(PROTO_PATH) \
		--go_out=$(PROTO_PATH) --go_opt=paths=source_relative \
		--go-grpc_out=$(PROTO_PATH) --go-grpc_opt=paths=source_relative \
		$(PROTO_PATH)/game/v1/game.proto
	protoc -I $(PROTO_PATH) \
		--go_out=$(PROTO_PATH) --go_opt=paths=source_relative \
		--go-grpc_out=$(PROTO_PATH) --go-grpc_opt=paths=source_relative \
		$(PROTO_PATH)/matchmaking/v1/matchmaking.proto

clean:
	rm -f $(PROTO_PATH)/game/v1/*.pb.go
	rm -f $(PROTO_PATH)/matchmaking/v1/*.pb.go
	rm -f web/static/*.wasm

run-engine:
	go run services/game-engine/cmd/main.go

run-matchmaking:
	go run services/matchmaking/cmd/main.go

run-gateway:
	go run services/gateway/cmd/main.go

run-terminal:
	go run cmd/terminal/main.go

run-browser-server:
	go run cmd/wasm-server/main.go

wasm:
	GOOS=js GOARCH=wasm go build -o web/static/app.wasm ./web/browser
	@echo "WASM built: web/static/app.wasm"

test:
	go test ./... -v

lint:
	golangci-lint run ./...
	@echo "ok"

pt:
	@echo "Starting Terminal Play (Solo mode, old flow)"
	go run services/game-engine/cmd/main.go &
	sleep 1
	go run cmd/terminal/main.go

pb:
	@echo "Building WASM and starting Browser Play (old flow)"
	GOOS=js GOARCH=wasm go build -o web/static/app.wasm ./web/browser
	go run services/game-engine/cmd/main.go &
	go run services/gateway/cmd/main.go &
	go run cmd/wasm-server/main.go &
	@echo "Services started. Open http://localhost:8081"

run-all:
	@echo "Starting ALL services (Matchmaking + Engine + Gateway)"
	@echo "   - Matchmaking: :50052"
	@echo "   - Game Engine: :50051"
	@echo "   - Gateway: :8080"
	@echo ""
	go run services/matchmaking/cmd/main.go &
	sleep 0.5
	go run services/game-engine/cmd/main.go &
	sleep 0.5
	go run services/gateway/cmd/main.go &
	@echo ""
	@echo "All services running!"
	@echo "   Now run: make run-terminal or make run-browser"
	@wait

run-all-terminal:
	@echo "Starting full stack + terminal client"
	@$(MAKE) -j3 run-matchmaking run-engine run-gateway & sleep 2 && go run cmd/terminal/main.go

run-all-browser:
	@echo "Starting full stack + browser client"
	GOOS=js GOARCH=wasm go build -o web/static/app.wasm ./web/browser
	@$(MAKE) -j3 run-matchmaking run-engine run-gateway & sleep 2 && go run cmd/wasm-server/main.go
	@echo "Open http://localhost:8081"

kill-ports:
	@echo "Killing processes on ports 50051, 50052, 8080, 8081..."
	@for port in 50051 50052 8080 8081; do \
		pid=$$(lsof -ti :$$port); \
		if [ -n "$$pid" ]; then \
			echo "Killing process $$pid on port $$port"; \
			kill -9 $$pid; \
		else \
			echo "No process found on port $$port"; \
		fi \
	done