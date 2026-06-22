package main

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "openxdr/proto"
)

func main() {

	// CHANGE THIS TO YOUR MACBOOK IP
	serverAddr := "192.168.1.173:50051"

	conn, err := grpc.Dial(
		serverAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := pb.NewTelemetryServiceClient(conn)

	agentID := "A001"
	hostname := "WIN-LAPTOP"

	fmt.Println("OpenXDR Agent started...")

	for {

		event := &pb.Event{
			AgentId:   agentID,
			Hostname:  hostname,
			EventType: "process",
			Payload:   "notepad.exe",
			Timestamp: time.Now().Format(time.RFC3339),
		}

		resp, err := client.SendEvent(context.Background(), event)
		if err != nil {
			fmt.Println("send error:", err)
		} else {
			fmt.Println("event sent:", resp.Success)
		}

		time.Sleep(5 * time.Second)
	}
}
