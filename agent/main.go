package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "openxdr/proto"
)

func getProcesses() []string {
	var cmd *exec.Cmd

	// Choose a process-listing command that is available on the current OS.
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("tasklist")
	default:
		// "comm=" prints only the process name with no header on Unix systems.
		cmd = exec.Command("ps", "-eo", "comm=")
	}

	out, err := cmd.Output()
	if err != nil {
		return []string{"error_getting_processes"}
	}

	scanner := bufio.NewScanner(strings.NewReader(string(out)))

	processes := []string{}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		switch runtime.GOOS {
		case "windows":
			// Skip tasklist header lines.
			if strings.Contains(line, "Image Name") ||
				strings.Contains(line, "===") {
				continue
			}

			// Extract the process name from the first column.
			fields := strings.Fields(line)
			if len(fields) > 0 {
				processes = append(processes, fields[0])
			}
		default:
			processes = append(processes, line)
		}
	}

	return processes
}

// operatingSystem returns a human-readable OS string in a cross-platform way.
func operatingSystem() string {
	switch runtime.GOOS {
	case "windows":
		if out, err := exec.Command("cmd", "/C", "ver").Output(); err == nil {
			if s := strings.TrimSpace(string(out)); s != "" {
				return s
			}
		}
	default:
		if out, err := exec.Command("uname", "-sr").Output(); err == nil {
			if s := strings.TrimSpace(string(out)); s != "" {
				return s
			}
		}
	}

	return runtime.GOOS + "/" + runtime.GOARCH
}

// hostName returns the machine hostname without depending on shell utilities.
func hostName() string {
	if name, err := os.Hostname(); err == nil && name != "" {
		return name
	}

	if out, err := exec.Command("hostname").Output(); err == nil {
		if s := strings.TrimSpace(string(out)); s != "" {
			return s
		}
	}

	return "unknown"
}

func main() {

	serverAddr := "192.168.1.25:50051"

	conn, err := grpc.NewClient(
		serverAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := pb.NewTelemetryServiceClient(conn)

	OSstring := operatingSystem()
	hostnameString := hostName()

	agentID := "A001"

	_, err = client.RegisterAgent(context.Background(), &pb.AgentInfo{
		AgentId:  agentID,
		Hostname: hostnameString,
		Os:       OSstring,
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
			Hostname:  hostnameString,
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
