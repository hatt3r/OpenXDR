#!/bin/bash

set -e

PROJECT_DIR="$HOME/Documents/OpenXDR"
PROMETHEUS_CONFIG="$PROJECT_DIR/prometheus/prometheus.yml"
AGENT_DIR="$PROJECT_DIR/agent"
AGENT_OUTPUT="$AGENT_DIR/agent.exe"
SERVER_BIN="$PROJECT_DIR/server/openxdr-server"
SERVER_LOG="$PROJECT_DIR/server/server.log"
PROMETHEUS_LOG="$PROJECT_DIR/prometheus/prometheus.log"
PID_DIR="$PROJECT_DIR/.pids"

echo ""
echo "=================================="
echo "       Starting OpenXDR"
echo "=================================="
echo ""

cd "$PROJECT_DIR" || exit 1

mkdir -p "$PID_DIR"

# ----------------------------------
# 1. Build Windows Agent
# ----------------------------------

echo "[1/5] Building Windows agent..."

rm -f "$AGENT_OUTPUT"

GOOS=windows GOARCH=amd64 go build -o "$AGENT_OUTPUT" ./agent

echo "Windows agent built successfully:"
echo "$AGENT_OUTPUT"
echo ""

# ----------------------------------
# 2. Build OpenXDR Server
# ----------------------------------

echo "[2/5] Building OpenXDR server..."

go build -o "$SERVER_BIN" ./server

echo "OpenXDR server built successfully."
echo ""

# ----------------------------------
# 3. Start Grafana
# ----------------------------------

echo "[3/5] Starting Grafana..."

brew services start grafana

echo ""

# ----------------------------------
# 4. Start Prometheus
# ----------------------------------

echo "[4/5] Starting Prometheus..."

# Stop a previous Prometheus instance if it is still running.
if [ -f "$PID_DIR/prometheus.pid" ]; then
    kill "$(cat "$PID_DIR/prometheus.pid")" 2>/dev/null || true
    rm -f "$PID_DIR/prometheus.pid"
fi
pkill -f "prometheus.*prometheus.yml" 2>/dev/null || true

nohup prometheus \
    --config.file="$PROMETHEUS_CONFIG" \
    > "$PROMETHEUS_LOG" 2>&1 &

echo $! > "$PID_DIR/prometheus.pid"

echo "Prometheus started."
echo ""

# ----------------------------------
# 5. Start OpenXDR Server
# ----------------------------------

echo "[5/5] Starting OpenXDR Server..."

# Stop a previous OpenXDR server if it is still running.
if [ -f "$PID_DIR/server.pid" ]; then
    kill "$(cat "$PID_DIR/server.pid")" 2>/dev/null || true
    rm -f "$PID_DIR/server.pid"
fi
pkill -f "go run ./server" 2>/dev/null || true
pkill -f "$SERVER_BIN" 2>/dev/null || true

nohup "$SERVER_BIN" \
    > "$SERVER_LOG" 2>&1 &

echo $! > "$PID_DIR/server.pid"

echo "OpenXDR server started."
echo ""

# ----------------------------------
# Finished
# ----------------------------------

echo "=================================="
echo "       OpenXDR is running"
echo "=================================="
echo ""

echo "Windows Agent:"
echo "  $AGENT_OUTPUT"
echo ""

echo "OpenXDR API:"
echo "  http://localhost:8080"

echo ""
echo "Metrics:"
echo "  http://localhost:8080/metrics"

echo ""
echo "Prometheus:"
echo "  http://localhost:9090"

echo ""
echo "Grafana:"
echo "  http://localhost:3000"

echo ""
echo "gRPC:"
echo "  localhost:50051"

echo ""
echo "Logs:"
echo "  Server:     $SERVER_LOG"
echo "  Prometheus: $PROMETHEUS_LOG"

echo ""
echo "=================================="