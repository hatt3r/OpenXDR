#!/bin/bash

PROJECT_DIR="$HOME/Documents/OpenXDR"
PROMETHEUS_CONFIG="$PROJECT_DIR/prometheus/prometheus.yml"
AGENT_DIR="$PROJECT_DIR/agent"
AGENT_OUTPUT="$AGENT_DIR/agent.exe"

echo ""
echo "=================================="
echo "       Starting OpenXDR"
echo "=================================="
echo ""

cd "$PROJECT_DIR" || exit 1

# ----------------------------------
# 1. Build Windows Agent
# ----------------------------------

echo "[1/4] Building Windows agent..."

rm -f "$AGENT_OUTPUT"

go build -o ./agent.exe ./agent/main.go

if [ $? -ne 0 ]; then
    echo "ERROR: Windows agent build failed."
    exit 1
fi

echo "Windows agent built successfully:"
echo "$AGENT_OUTPUT"
echo ""

# ----------------------------------
# 2. Start Grafana
# ----------------------------------

echo "[2/4] Starting Grafana..."

brew services start grafana

echo ""

# ----------------------------------
# 3. Start Prometheus
# ----------------------------------

echo "[3/4] Starting Prometheus..."

# Kill previous Prometheus process if running
pkill -f "prometheus.*prometheus.yml" 2>/dev/null

nohup prometheus \
    --config.file="$PROMETHEUS_CONFIG" \
    > "$PROJECT_DIR/prometheus/prometheus.log" 2>&1 &

echo "Prometheus started."
echo ""

# ----------------------------------
# 4. Start OpenXDR Server
# ----------------------------------

echo "[4/4] Starting OpenXDR Server..."

# Kill previous OpenXDR server if running
pkill -f "$PROJECT_DIR/server" 2>/dev/null

nohup go run ./server \
    > "$PROJECT_DIR/server/server.log" 2>&1 &

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

echo "Metrics:"
echo "  http://localhost:8080/metrics"

echo "Prometheus:"
echo "  http://localhost:9090"

echo "Grafana:"
echo "  http://localhost:3000"

echo ""
echo "Logs:"
echo "  Server:     $PROJECT_DIR/server/server.log"
echo "  Prometheus: $PROJECT_DIR/prometheus/prometheus.log"
echo ""

echo "=================================="
