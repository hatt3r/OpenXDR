# 🛡️ OpenXDR

OpenXDR is an open-source, Go-based experimental Extended Detection and Response (XDR) system. It demonstrates a full telemetry pipeline from an endpoint agent to a central server, with process inspection, gRPC ingestion, SQLite storage, detection rules, a REST API, Prometheus metrics, and Grafana dashboards.

The project is currently in **Phase 2**.

## ✨ Features

- Endpoint agent collects running processes.
- Agent sends telemetry to the central server over gRPC on port `50051`.
- Server stores agents, events, and alerts in SQLite (`openxdr.db`).
- Detection engine flags suspicious processes and command-line indicators.
- REST API exposes agents, events, alerts, health, and Prometheus metrics.
- Prometheus scrapes server metrics on port `8080`.
- Grafana visualizes the collected metrics on port `3000`.
- Helper scripts start and stop the full local stack.

## 🏗️ Architecture

```
┌──────────────────────────────┐
│ Endpoint (Agent)             │
│                              │
│ - Process collector          │
│ - gRPC client                │
│ - Telemetry sender           │
└──────────────┬───────────────┘
               │ gRPC :50051
               ▼
┌──────────────────────────────┐
│ OpenXDR Server (SOC Core)    │
│                              │
│ - gRPC server                │
│ - Event logger (Zap)         │
│ - SQLite storage             │
│ - Detection engine           │
│ - REST API                   │
│ - Agent registration         │
│ - Heartbeats                 │
│ - Alerts                     │
│ - Prometheus metrics         │
└──────────────┬───────────────┘
               │ scrape :8080/metrics
               ▼
┌──────────────────────────────┐
│ Prometheus / Grafana         │
│                              │
│ - Metrics store :9090        │
│ - Dashboards :3000           │
└──────────────────────────────┘
```

## 📁 Project Structure

```
openxdr/
├── agent/              # Endpoint agent
├── server/             # Central ingestion server
├── api/                # REST API handlers and router
├── internal/
│   ├── detection/      # Detection rules
│   ├── logger/         # Zap logging
│   └── store/          # SQLite storage
├── proto/              # gRPC contract and generated code
├── prometheus/         # Prometheus configuration
├── startopenxdr.sh     # Start the full stack
├── stopopenxdr.sh      # Stop the full stack
├── go.mod
└── README.md
```

## 💻 Tech Stack

- Go
- gRPC and Protocol Buffers
- Gin REST API
- Zap logging
- SQLite
- Prometheus
- Grafana

## 🚀 Getting Started

### Prerequisites

- Go (matching `go.mod`)
- Prometheus installed and available on `PATH`
- Grafana installed through Homebrew, if using the start script
- A macOS/Linux shell with `bash`

### 1. Configure the agent

Update the server address in `agent/main.go`:

```go
serverAddr := "YOUR_SERVER_IP:50051"
```

### 2. Make scripts executable

```bash
chmod +x startopenxdr.sh stopopenxdr.sh
```

### 3. Start OpenXDR

```bash
./startopenxdr.sh
```

The script:

- Cross-compiles the Windows agent to `agent/agent.exe` (`GOOS=windows GOARCH=amd64`).
- Cross-compiles the Linux agent to `agent/agent_linux` (`GOOS=linux GOARCH=amd64`).
- Builds the server into `server/openxdr-server`.
- Starts Grafana.
- Starts Prometheus with `prometheus/prometheus.yml`.
- Starts the OpenXDR server.

### 4. Stop OpenXDR

```bash
./stopopenxdr.sh
```

## 🌐 Service URLs

| Service | URL |
| --- | --- |
| OpenXDR API | http://localhost:8080 |
| Health | http://localhost:8080/health |
| Agents | http://localhost:8080/agents |
| Events | http://localhost:8080/events |
| Alerts | http://localhost:8080/alerts |
| Metrics | http://localhost:8080/metrics |
| gRPC server | http://localhost:50051 |
| Prometheus | http://localhost:9090 |
| Grafana | http://localhost:3000 |

## 📡 Telemetry Flow

1. The endpoint agent collects running processes.
2. Process data is packaged into an `Event`.
3. The event is sent through gRPC every 5 seconds.
4. The server logs, stores, and evaluates the event.
5. Detection alerts are stored in SQLite.
6. Prometheus metrics are exposed and scraped.
7. Grafana can visualize the metrics.

## 🗄️ Data Storage

SQLite is used by default. The database file is:

```text
openxdr.db
```

Tables:

- `agents`
- `events`
- `alerts`

## 📋 Logs

- Server: `server/server.log`
- Prometheus: `prometheus/prometheus.log`
- PID files: `.pids/`

## 🧪 Manual Run (without scripts)

Start the server:

```bash
go run ./server
```

Run the agent on an endpoint:

```bash
go run ./agent
```

For a Windows agent binary, cross-compile from macOS or Linux:

```bash
GOOS=windows GOARCH=amd64 go build -o agent/agent.exe ./agent
```

For a Linux agent binary, cross-compile from macOS or another platform:

```bash
GOOS=linux GOARCH=amd64 go build -o agent/agent_linux ./agent
```
