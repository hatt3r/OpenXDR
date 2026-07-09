package main

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "openxdr/proto"
)

func getProcesses() []string {

	// Windows command to list processes
	cmd := exec.Command("tasklist")

	out, err := cmd.Output()
	if err != nil {
		return []string{"error_getting_processes"}
	}

	scanner := bufio.NewScanner(strings.NewReader(string(out)))

	processes := []string{}

	for scanner.Scan() {
		line := scanner.Text()

		// Skip header lines
		if strings.Contains(line, "Image Name") ||
			strings.Contains(line, "===") ||
			len(line) < 10 {
			continue
		}

		// Extract process name (first column)
		fields := strings.Fields(line)
		if len(fields) > 0 {
			processes = append(processes, fields[0])
		}
	}

	return processes
}

func main() {

	serverAddr := "192.168.1.173:50051"

	conn, err := grpc.NewClient(
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

	_, err = client.RegisterAgent(context.Background(), &pb.AgentInfo{
		AgentId:  agentID,
		Hostname: hostname,
		Os:       "windows",
	})
	if err != nil {
		fmt.Println("register error:", err)
	}

	go func() {
		for {
			_, err := client.Heartbeat(context.Background(), &pb.HeartbeatRequest{
				AgentId:   agentID,
				Timestamp: time.Now().Format(time.RFC3339),
			})

			if err != nil {
				fmt.Println("heartbeat error:", err)
			}
			time.Sleep(10 * time.Second)
		}
	}()

	fmt.Println("OpenXDR Agent started (REAL MODE)")

	for {

		processes := getProcesses()
		if len(processes) > 50 {
			processes = processes[:50]
		}
		payload := strings.Join(processes, ",")
		// payload := "chrome.exe,mimikatz.exe"
		event := &pb.Event{
			AgentId:   agentID,
			Hostname:  hostname,
			EventType: "process_list",
			Payload:   payload,
			Timestamp: time.Now().Format(time.RFC3339),
		}

		resp, err := client.SendEvent(context.Background(), event)
		if err != nil {
			fmt.Println("send error:", err)
		} else {
			fmt.Println("process telemetry sent:", resp.Success)
		}

		time.Sleep(5 * time.Second)
	}
}
