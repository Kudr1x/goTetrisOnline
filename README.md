# GoTetrisOnline

![CI](https://github.com/Kudr1x/GoTetrisOnline/workflows/CI/badge.svg)
![Go Version](https://img.shields.io/badge/Go-1.25-00ADD8?logo=go)
![License](https://img.shields.io/badge/license-MIT-blue.svg)

Multiplayer Tetris game with microservices architecture written in Go.

## Features

- Solo and 1v1 game modes
- Real-time multiplayer with WebSocket
- WASM client for browser play
- Terminal UI client
- gRPC microservices architecture
- Invite system with shareable links
- Random matchmaking

## Architecture

```
┌─────────────┐
│   Client    │  (Terminal/Browser WASM)
└──────┬──────┘
       │ WebSocket
       ▼
┌─────────────┐
│   Gateway   │  :8080
└──────┬──────┘
       │
       ├───────────────┬──────────────┐
       │               │              │
       ▼               ▼              ▼
┌──────────────┐ ┌─────────────┐ ┌──────┐
│ Matchmaking  │ │ Game Engine │ │ ... │
│   :50052     │ │   :50051    │ └──────┘
└──────────────┘ └─────────────┘
```

## Services

- **Matchmaking** - Match creation, invite codes, random opponent search
- **Game Engine** - Tetris game logic, state management
- **Gateway** - WebSocket/HTTP gateway, protocol translation

## Quick Start

### Prerequisites

- Go 1.25+
- Protocol Buffers compiler
- Make

### Installation

```bash
# Clone repository
git clone https://github.com/Kudr1x/GoTetrisOnline.git
cd GoTetrisOnline

# Install dependencies
go mod download

# Generate protobuf files
make gen
```

### Run locally

```bash
# Start all services
make run-all

# In another terminal, start terminal client
make run-terminal

# Or start browser client
make run-browser
# Then open http://localhost:8081
```

### Individual services

```bash
# Start individual services
make run-matchmaking  # :50052
make run-engine       # :50051
make run-gateway      # :8080
```

## Development

### Commands

```bash
make gen          # Generate protobuf files
make test         # Run all tests
make lint         # Run linter
make clean        # Clean generated files
make wasm         # Build WASM client
make kill-ports   # Kill processes on service ports
```

### Testing

```bash
# Run all tests with coverage
go test -v -cover ./...

# Run specific package tests
go test -v ./services/matchmaking/...
```

### Project Structure

```
.
├── api/proto/                # Protocol Buffers definitions
│   ├── game/v1/
│   └── matchmaking/v1/
├── cmd/                      # Client applications
│   ├── terminal/            # Terminal UI client
│   └── wasm-server/         # WASM file server
├── pkg/                      # Shared packages
│   ├── core/                # Tetris game logic
│   └── renderer/            # UI rendering
├── services/                 # Microservices
│   ├── game-engine/         # Game state management
│   ├── gateway/             # WebSocket gateway
│   └── matchmaking/         # Matchmaking service
└── web/                      # WASM client
    ├── browser/
    └── static/
```

## API

### Matchmaking gRPC API

```protobuf
service MatchmakingService {
  rpc CreateMatch(CreateMatchRequest) returns (CreateMatchResponse);
  rpc JoinByInvite(JoinByInviteRequest) returns (JoinByInviteResponse);
  rpc FindRandom(FindRandomRequest) returns (FindRandomResponse);
  rpc GetMatchInfo(GetMatchInfoRequest) returns (GetMatchInfoResponse);
}
```

### Game Modes

- **Solo** - Single player practice mode
- **1v1** - Two players compete, send garbage lines on Tetris clears

## Contributing

Pull requests are welcome. For major changes, please open an issue first.

## License

MIT
