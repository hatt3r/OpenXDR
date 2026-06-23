package detection

import "strings"

type Alert struct {
	Title       string
	Severity    string
	Description string
}

func CheckProcesses(payload string) []Alert {

	alerts := []Alert{}

	lowerPayload := strings.ToLower(payload)

	suspicious := []string{
		"mimikatz.exe",
		"procdump.exe",
		"psexec.exe",
		"nc.exe",
		"netcat.exe",
	}

	for _, proc := range suspicious {

		if strings.Contains(lowerPayload, proc) {

			alerts = append(alerts, Alert{
				Title:       "Suspicious Process Detected",
				Severity:    "HIGH",
				Description: proc + " observed on endpoint",
			})
		}
	}

	return alerts
}
