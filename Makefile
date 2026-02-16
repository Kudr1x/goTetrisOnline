.PHONY: gen clean

PROTO_PATH = api/proto

gen:
	protoc -I $(PROTO_PATH) \
		--go_out=$(PROTO_PATH) --go_opt=paths=source_relative \
		--go-grpc_out=$(PROTO_PATH) --go-grpc_opt=paths=source_relative \
		$(PROTO_PATH)/game/v1/game.proto
	@echo "ok"

clean:
	rm -f $(PROTO_PATH)/game/v1/*.pb.go

run-engine:
	go run services/game-engine/cmd/main.go

run-gateway:
	go run services/gateway/cmd/main.go

run-client:
	go run client/cmd/main.go


.PHONY: lint
lint:
	golangci-lint run ./...
	@echo "ok"