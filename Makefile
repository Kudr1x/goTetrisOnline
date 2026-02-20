.PHONY: gen clean test lint run-all

PROTO_PATH = api/proto

gen:
	protoc -I $(PROTO_PATH) \
		--go_out=$(PROTO_PATH) --go_opt=paths=source_relative \
		--go-grpc_out=$(PROTO_PATH) --go-grpc_opt=paths=source_relative \
		$(PROTO_PATH)/game/v1/game.proto

clean:
	rm -f $(PROTO_PATH)/game/v1/*.pb.go
	rm -f web/static/*.wasm

run-engine:
	go run services/game-engine/cmd/main.go

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
	go run services/game-engine/cmd/main.go &
	go run cmd/terminal/main.go

pb:
	GOOS=js GOARCH=wasm go build -o web/static/app.wasm ./web/browser
	go run services/game-engine/cmd/main.go &
	go run services/gateway/cmd/main.go &
	go run cmd/wasm-server/main.go &

kill-ports:
	@echo "Killing processes on ports 50051, 8080, 8081..."
	@for port in 50051 8080 8081; do \
		pid=$$(lsof -ti :$$port); \
		if [ -n "$$pid" ]; then \
			echo "Killing process $$pid on port $$port"; \
			kill -9 $$pid; \
		else \
			echo "No process found on port $$port"; \
		fi \
	done