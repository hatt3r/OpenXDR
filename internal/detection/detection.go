package detection

import "strings"

type Alert struct {
	Title       string
	Severity    string
	Description string
}

func CheckProcesses(payload string) []Alert {

	alerts := []Alert{}
	if payload == "" {
		return alerts
	}

	lowerPayload := strings.ToLower(payload)

	// Suspicious process and activity indicators. Keep the original
	// high-signal entries first, then append broader detection coverage.
	suspicious := []string{
		// Credentials and memory access
		"mimikatz.exe",
		"procdump.exe",
		"procdump64.exe",
		"lsass",
		"ntds.dit",
		"secretsdump",
		"lazagne",
		"rubeus",
		"pypykatz",
		"gsecdump",
		"wce.exe",
		"pwdump",
		"fgdump",
		"chntpw",

		// Remote execution and lateral movement
		"psexec.exe",
		"psexec",
		"psexec64.exe",
		"paexec",
		"smbexec",
		"wmiexec",
		"evil-winrm",
		"schtasks",
		"at.exe",

		// Reverse shells and network pivoting
		"nc.exe",
		"ncat.exe",
		"netcat.exe",
		"netcat",
		"socat",
		"responder",
		"bettercap",

		// Reconnaissance and enumeration
		"bloodhound",
		"sharphound",
		"adrecon",
		"seatbelt",
		"winpeas",
		"linpeas",

		// Living-off-the-land utilities
		"certutil",
		"bitsadmin",
		"mshta",
		"regsvr32",
		"rundll32",
		"cmstp",
		"installutil",
		"msiexec",

		// Password cracking, brute force, and scanning
		"hashcat",
		"john",
		"hydra",
		"medusa",
		"sqlmap",
		"nmap",
		"masscan",
		"fscan",

		// Suspicious command-line patterns
		"-enc",
		"-encodedcommand",
		"iex",
		"invoke-expression",
		"downloadstring",
		"frombase64string",
		"hidden",
		"reverse",
	}

	// Collect every indicator present in the payload first.
	matched := []string{}
	for _, proc := range suspicious {
		if strings.Contains(lowerPayload, proc) {
			matched = append(matched, proc)
		}
	}

	// Drop indicators that are already covered by a more specific matched
	// indicator (for example psexec when psexec.exe is also present).
	filtered := []string{}
	for i, candidate := range matched {
		dominated := false
		for j, other := range matched {
			if i == j || len(other) <= len(candidate) {
				continue
			}
			if strings.Contains(other, candidate) {
				dominated = true
				break
			}
		}
		if !dominated {
			filtered = append(filtered, candidate)
		}
	}

	// Avoid exact duplicate alerts.
	seen := make(map[string]struct{})
	for _, proc := range filtered {
		key := strings.ToLower(proc)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}

		alerts = append(alerts, Alert{
			Title:       "Suspicious Process Detected",
			Severity:    "HIGH",
			Description: proc + " observed on endpoint",
		})
	}

	return alerts
}
