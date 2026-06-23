## 🛡️ OpenXDR (currently in phase 1)

OpenXDR is an open-source, Go-based experimental Extended Detection and Response (XDR) system designed to simulate how modern security platforms collect, transmit, and process endpoint telemetry.

This project is built for learning purposes and demonstrates a full telemetry pipeline using gRPC between a Windows endpoint agent and a MacBook-based central server.

while this is in the nascent stage, it is a stepping stone for later stages.


---

# 🎯 What This Project Does

OpenXDR builds a simple but realistic security monitoring pipeline:


endpoint Agent → gRPC → admin Server


The system collects process-level telemetry from a endpoint machine (preferable windows) and streams it in real time to a central ingestion server.

---

# 🧠 Core Idea

Instead of building dashboards or UI first, OpenXDR focuses on the **core of all security platforms: telemetry ingestion**.

---

# 🏗️ Architecture

```
┌──────────────────────────────┐
│ endpoint Laptop (Agent) │
│ │
│ - Process Collector │
│ - gRPC Client │
│ - Telemetry Sender │
└──────────────┬───────────────┘

│ gRPC :50051
▼

┌──────────────┘────────────────┐
│ admin Server (SOC Core) │
│ │
│ - gRPC Server │
│ - Event Logger (Zap) │
│ - In-memory Storage │
└──────────────────────────────┘
```

---

# 💻 Tech Stack

- Go (Golang)
- gRPC
- Protocol Buffers
- Zap logging
- Windows `tasklist` (Phase 1 telemetry source)

---

# 📁 Project Structure

```
openxdr/
├── agent/ #  endpoint agent
├── server/ #  ingestion server
├── proto/ # gRPC contract (telemetry definition)
├── internal/logger/ # structured logging system
├── go.mod
└── README.md
```

---

# 📡 Telemetry Flow

1. endpoint agent collects running processes
2. Data is packaged into an Event
3. Event is sent via gRPC every 5 seconds
4. admin server receives and logs the event

---

# ▶️ How to Run

Update server IP in agent:

```serverAddr := "YOUR_dev_server_IP:50051"```

1. Start Server

```go run server/main.go```

2. Run Agent on endpoirnt

```go run agent/main.go```
