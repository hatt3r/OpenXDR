# 🛡️ OpenXDR

OpenXDR is an open-source, Go-based experimental Extended Detection and Response (XDR) system designed to simulate how modern security platforms collect, transmit, and process endpoint telemetry.

This project is built for learning purposes and demonstrates a full telemetry pipeline using gRPC between a Windows endpoint agent and a MacBook-based central server.

---

# 🎯 What This Project Does

OpenXDR builds a simple but realistic security monitoring pipeline:


Windows Agent → gRPC → MacBook Server


The system collects process-level telemetry from a Windows machine and streams it in real time to a central ingestion server.

---

# 🧠 Core Idea

Instead of building dashboards or UI first, OpenXDR focuses on the **core of all security platforms: telemetry ingestion**.

---

# 🏗️ Architecture

```
┌──────────────────────────────┐
│ Windows Laptop (Agent) │
│ │
│ - Process Collector │
│ - gRPC Client │
│ - Telemetry Sender │
└──────────────┬───────────────┘
│ gRPC :50051
▼
┌──────────────────────────────┐
│ MacBook Server (SOC Core) │
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
├── agent/ # Windows endpoint agent
├── server/ # MacBook ingestion server
├── proto/ # gRPC contract (telemetry definition)
├── internal/logger/ # structured logging system
├── go.mod
└── README.md
```

---

# 📡 Telemetry Flow

1. Windows agent collects running processes
2. Data is packaged into an Event
3. Event is sent via gRPC every 5 seconds
4. MacBook server receives and logs the event

---

# 📦 Event Schema

proto
```
message Event {
  string agent_id = 1;
  string hostname = 2;
  string event_type = 3;
  string payload = 4;
  string timestamp = 5;
}
```

▶️ How to Run

Update server IP in agent:

```serverAddr := "YOUR_dev_server_IP:50051"```

1. Start Server 
```go run server/main.go```

2. Run Agent on endpoirnt

```go run agent/main.go```
