package main

import (
	"context"
	"net"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"openxdr/internal/detection"
	"openxdr/internal/logger"
	"openxdr/internal/store"
	pb "openxdr/proto"
)

type Agent struct {
	AgentID  string
	Hostname string
	LastSeen string
	Status   string
}

type Server struct {
	pb.UnimplementedTelemetryServiceServer
	events []pb.Event
}

var (
	agents  = map[string]*Agent{}
	dbStore *store.Store
)

//register agent

func (s *Server) RegisterAgent(ctx context.Context, a *pb.AgentInfo) (*pb.Ack, error) {

	// memory registry
	agents[a.AgentId] = &Agent{
		AgentID:  a.AgentId,
		Hostname: a.Hostname,
		LastSeen: time.Now().Format(time.RFC3339),
		Status:   "ONLINE",
	}

	logger.Log.Info("agent registered",
		zap.String("agent_id", a.AgentId),
		zap.String("hostname", a.Hostname),
	)

	//stick to db
	_, _ = dbStore.DB.Exec(`
	INSERT OR REPLACE INTO agents(agent_id, hostname, status, last_seen)
	VALUES (?, ?, ?, ?)
	`, a.AgentId, a.Hostname, "ONLINE", time.Now().Format(time.RFC3339))

	return &pb.Ack{Success: true}, nil
}

//Can you feel my heartbeat?

func (s *Server) Heartbeat(ctx context.Context, h *pb.HeartbeatRequest) (*pb.Ack, error) {

	// update memory
	if agent, ok := agents[h.AgentId]; ok {
		agent.LastSeen = time.Now().Format(time.RFC3339)
		agent.Status = "ONLINE"
	}

	logger.Log.Info("heartbeat received",
		zap.String("agent_id", h.AgentId),
	)

	// persist heartbeat
	_, _ = dbStore.DB.Exec(`
	UPDATE agents
	SET last_seen = ?, status = ?
	WHERE agent_id = ?
	`, time.Now().Format(time.RFC3339), "ONLINE", h.AgentId)

	return &pb.Ack{Success: true}, nil
}

//event telemetry

func (s *Server) SendEvent(ctx context.Context, e *pb.Event) (*pb.Ack, error) {

	// log event
	logger.Log.Info("event received",
		zap.String("agent_id", e.AgentId),
		zap.String("hostname", e.Hostname),
		zap.String("event_type", e.EventType),
		zap.String("payload", e.Payload),
		zap.String("timestamp", e.Timestamp),
	)

	_, _ = dbStore.DB.Exec(`
	INSERT INTO events(agent_id, hostname, event_type, payload, timestamp)
	VALUES (?, ?, ?, ?, ?)
	`, e.AgentId, e.Hostname, e.EventType, e.Payload, e.Timestamp)

	// detection engin
	alerts := detection.CheckProcesses(e.Payload)

	for _, alert := range alerts {

		logger.Log.Warn("ALERT GENERATED",
			zap.String("agent_id", e.AgentId),
			zap.String("title", alert.Title),
			zap.String("severity", alert.Severity),
			zap.String("description", alert.Description),
		)

		// store alert in DB
		_, _ = dbStore.DB.Exec(`
		INSERT INTO alerts(agent_id, title, severity, description, timestamp)
		VALUES (?, ?, ?, ?, ?)
		`, e.AgentId, alert.Title, alert.Severity, alert.Description, time.Now().Format(time.RFC3339))
	}

	// memory fallback (optional)
	s.events = append(s.events, *e)

	return &pb.Ack{Success: true}, nil
}

//main

func main() {

	//intialize logger
	logger.Init()
	defer logger.Sync()

	//initialize database
	dbStore = store.New("openxdr.db")

	// gRPC server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterTelemetryServiceServer(grpcServer, &Server{
		events: make([]pb.Event, 0),
	})

	logger.Log.Info("OpenXDR Server started on :50051")

	// ready? start server
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}
