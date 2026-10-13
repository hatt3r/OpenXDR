package main

import (
	"context"
	"net"
	api "openxdr/api"
	"openxdr/internal/detection"
	"openxdr/internal/logger"
	"openxdr/internal/store"
	pb "openxdr/proto"
	"openxdr/server/internal/metrics"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
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

// Register agent
func (s *Server) RegisterAgent(
	ctx context.Context,
	a *pb.AgentInfo,
) (*pb.Ack, error) {

	now := time.Now().Format(time.RFC3339)

	// Memory registry
	agents[a.AgentId] = &Agent{
		AgentID:  a.AgentId,
		Hostname: a.Hostname,
		LastSeen: now,
		Status:   "ONLINE",
	}

	// Prometheus metrics
	metrics.ActiveAgents.WithLabelValues(
		a.AgentId,
		a.Hostname,
		a.Os,
	).Set(1)

	metrics.GRPCRequestsTotal.WithLabelValues(
		"RegisterAgent",
		"success",
	).Inc()

	logger.Log.Info(
		"agent registered",
		zap.String("agent_id", a.AgentId),
		zap.String("hostname", a.Hostname),
	)

	// Store in database
	_, err := dbStore.DB.Exec(`
		INSERT OR REPLACE INTO agents(
			agent_id,
			hostname,
			status,
			last_seen
		)
		VALUES (?, ?, ?, ?)
	`,
		a.AgentId,
		a.Hostname,
		"ONLINE",
		now,
	)

	if err != nil {
		metrics.GRPCRequestsTotal.WithLabelValues(
			"RegisterAgent",
			"error",
		).Inc()

		logger.Log.Error(
			"failed to store agent",
			zap.String("agent_id", a.AgentId),
			zap.Error(err),
		)

		return &pb.Ack{Success: false}, err
	}

	return &pb.Ack{Success: true}, nil
}

// Heartbeat
func (s *Server) Heartbeat(
	ctx context.Context,
	h *pb.HeartbeatRequest,
) (*pb.Ack, error) {

	now := time.Now().Format(time.RFC3339)

	// Update memory
	if agent, ok := agents[h.AgentId]; ok {
		agent.LastSeen = now
		agent.Status = "ONLINE"
	}

	metrics.GRPCRequestsTotal.WithLabelValues(
		"Heartbeat",
		"success",
	).Inc()

	logger.Log.Info(
		"heartbeat received",
		zap.String("agent_id", h.AgentId),
	)

	// Persist heartbeat
	_, err := dbStore.DB.Exec(`
		UPDATE agents
		SET last_seen = ?, status = ?
		WHERE agent_id = ?
	`,
		now,
		"ONLINE",
		h.AgentId,
	)

	if err != nil {
		metrics.GRPCRequestsTotal.WithLabelValues(
			"Heartbeat",
			"error",
		).Inc()

		logger.Log.Error(
			"failed to update heartbeat",
			zap.String("agent_id", h.AgentId),
			zap.Error(err),
		)

		return &pb.Ack{Success: false}, err
	}

	return &pb.Ack{Success: true}, nil
}

// Event telemetry
func (s *Server) SendEvent(
	ctx context.Context,
	e *pb.Event,
) (*pb.Ack, error) {

	// Event received metric
	metrics.EventsReceived.WithLabelValues(
		e.AgentId,
		e.EventType,
	).Inc()

	logger.Log.Info(
		"event received",
		zap.String("agent_id", e.AgentId),
		zap.String("hostname", e.Hostname),
		zap.String("event_type", e.EventType),
		zap.String("payload", e.Payload),
		zap.String("timestamp", e.Timestamp),
	)

	// Store event in database
	_, err := dbStore.DB.Exec(`
		INSERT INTO events(
			agent_id,
			hostname,
			event_type,
			payload,
			timestamp
		)
		VALUES (?, ?, ?, ?, ?)
	`,
		e.AgentId,
		e.Hostname,
		e.EventType,
		e.Payload,
		e.Timestamp,
	)

	if err != nil {

		metrics.EventsProcessed.WithLabelValues(
			e.AgentId,
			"failed",
		).Inc()

		metrics.GRPCRequestsTotal.WithLabelValues(
			"SendEvent",
			"error",
		).Inc()

		logger.Log.Error(
			"failed to store event",
			zap.String("agent_id", e.AgentId),
			zap.Error(err),
		)

		return &pb.Ack{Success: false}, err
	}

	// Event successfully processed
	metrics.EventsProcessed.WithLabelValues(
		e.AgentId,
		"success",
	).Inc()

	metrics.GRPCRequestsTotal.WithLabelValues(
		"SendEvent",
		"success",
	).Inc()

	// Detection engine: only process telemetry should be evaluated for
	// suspicious-process rules, and empty payloads produce no alerts.
	if e.EventType == "process_list" && e.Payload != "" {
		alerts := detection.CheckProcesses(e.Payload)

		for _, alert := range alerts {

			// Prometheus alert metric
			metrics.AlertsGenerated.WithLabelValues(
				alert.Title,
				alert.Severity,
			).Inc()

			logger.Log.Warn(
				"ALERT GENERATED",
				zap.String("agent_id", e.AgentId),
				zap.String("title", alert.Title),
				zap.String("severity", alert.Severity),
				zap.String("description", alert.Description),
			)

			// Store alert in database
			_, err := dbStore.DB.Exec(`
			INSERT INTO alerts(
				agent_id,
				title,
				severity,
				description,
				timestamp
			)
			VALUES (?, ?, ?, ?, ?)
		`,
				e.AgentId,
				alert.Title,
				alert.Severity,
				alert.Description,
				time.Now().Format(time.RFC3339),
			)

			if err != nil {
				logger.Log.Error(
					"failed to store alert",
					zap.String("agent_id", e.AgentId),
					zap.Error(err),
				)
			}
		}
	}

	// Memory fallback
	s.events = append(s.events, *e)

	return &pb.Ack{Success: true}, nil
}

func main() {

	// Initialize logger
	logger.Init()
	defer logger.Sync()

	// Initialize database
	dbStore = store.New("openxdr.db")

	// Give API access to database
	api.Init(dbStore)

	// gRPC server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterTelemetryServiceServer(
		grpcServer,
		&Server{
			events: make([]pb.Event, 0),
		},
	)

	logger.Log.Info(
		"OpenXDR Server started on :50051",
	)

	// REST API
	router := api.SetupRouter()

	go func() {
		if err := router.Run(":8080"); err != nil {
			panic(err)
		}
	}()

	// Start gRPC server
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}
