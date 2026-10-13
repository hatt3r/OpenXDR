#!/bin/bash

PROJECT_DIR="$HOME/Documents/OpenXDR"
SERVER_BIN="$PROJECT_DIR/server/openxdr-server"
PID_DIR="$PROJECT_DIR/.pids"

echo ""
echo "=================================="
echo "       Stopping OpenXDR"
echo "=================================="
echo ""

# ----------------------------------
# 1. Stop OpenXDR Server
# ----------------------------------

echo "[1/3] Stopping OpenXDR Server..."

if [ -f "$PID_DIR/server.pid" ]; then
    pid="$(cat "$PID_DIR/server.pid")"
    kill "$pid" 2>/dev/null || true
    for _ in 1 2 3 4 5; do
        kill -0 "$pid" 2>/dev/null || break
        sleep 0.2
    done
    kill -9 "$pid" 2>/dev/null || true
    rm -f "$PID_DIR/server.pid"
fi

pkill -f "go run ./server" 2>/dev/null || true
pkill -f "$SERVER_BIN" 2>/dev/null || true

echo "OpenXDR Server stopped."
echo ""

# ----------------------------------
# 2. Stop Prometheus
# ----------------------------------

echo "[2/3] Stopping Prometheus..."

if [ -f "$PID_DIR/prometheus.pid" ]; then
    pid="$(cat "$PID_DIR/prometheus.pid")"
    kill "$pid" 2>/dev/null || true
    for _ in 1 2 3 4 5; do
        kill -0 "$pid" 2>/dev/null || break
        sleep 0.2
    done
    kill -9 "$pid" 2>/dev/null || true
    rm -f "$PID_DIR/prometheus.pid"
fi

pkill -f "prometheus.*prometheus.yml" 2>/dev/null || true

echo "Prometheus stopped."
echo ""

# ----------------------------------
# 3. Stop Grafana
# ----------------------------------

echo "[3/3] Stopping Grafana..."

brew services stop grafana 2>/dev/null || true

echo "Grafana stopped."
echo ""

echo "=================================="
echo "       OpenXDR stopped"
echo "=================================="
echo ""