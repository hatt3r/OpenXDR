package main

import (
	"context"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"openxdr/internal/logger"
	pb "openxdr/proto"

	"openxdr/internal/detection"
)

type Server struct {
	pb.UnimplementedTelemetryServiceServer
	events []pb.Event
}

func (s *Server) SendEvent(
	ctx context.Context,
	e *pb.Event,
) (*pb.Ack, error) {

	// Log incoming event
	logger.Log.Info("event received",
		zap.String("agent_id", e.AgentId),
		zap.String("hostname", e.Hostname),
		zap.String("event_type", e.EventType),
		zap.String("payload", e.Payload),
		zap.String("timestamp", e.Timestamp),
	)

	alerts := detection.CheckProcesses(e.Payload)

	for _, alert := range alerts {

		logger.Log.Warn(
			"ALERT GENERATED",
			zap.String("title", alert.Title),
			zap.String("severity", alert.Severity),
			zap.String("description", alert.Description),
		)
	}
	// Store in memory (temporary Phase 1 storage)
	s.events = append(s.events, *e)

	return &pb.Ack{
		Success: true,
	}, nil
}

func main() {

	// Init logger
	logger.Init()
	defer logger.Sync()

	// Listen on gRPC port
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()

	// Register service
	s := &Server{
		events: make([]pb.Event, 0),
	}

	pb.RegisterTelemetryServiceServer(grpcServer, s)

	logger.Log.Info("OpenXDR Server started on :50051")

	// Start server
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}
