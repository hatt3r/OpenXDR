#!/bin/bash

PROJECT_DIR="$HOME/Documents/OpenXDR"

echo ""
echo "=================================="
echo "       Stopping OpenXDR"
echo "=================================="
echo ""

# ----------------------------------
# 1. Stop OpenXDR Server
# ----------------------------------

echo "[1/3] Stopping OpenXDR Server..."

pkill -f "$PROJECT_DIR/server" 2>/dev/null

echo "OpenXDR Server stopped."
echo ""

# ----------------------------------
# 2. Stop Prometheus
# ----------------------------------

echo "[2/3] Stopping Prometheus..."

pkill -f "prometheus.*prometheus.yml" 2>/dev/null

echo "Prometheus stopped."
echo ""

# ----------------------------------
# 3. Stop Grafana
# ----------------------------------

echo "[3/3] Stopping Grafana..."

brew services stop grafana

echo "Grafana stopped."
echo ""

echo "=================================="
echo "       OpenXDR stopped"
echo "=================================="
echo ""