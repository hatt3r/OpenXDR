package store

type Agent struct {
	AgentID  string `json:"agent_id"`
	Hostname string `json:"hostname"`
	Status   string `json:"status"`
	LastSeen string `json:"last_seen"`
}

type Alert struct {
	ID          int    `json:"id"`
	AgentID     string `json:"agent_id"`
	Title       string `json:"title"`
	Severity    string `json:"severity"`
	Description string `json:"description"`
	Timestamp   string `json:"timestamp"`
}

type Event struct {
	ID        int    `json:"id"`
	AgentID   string `json:"agent_id"`
	Hostname  string `json:"hostname"`
	EventType string `json:"event_type"`
	Payload   string `json:"payload"`
	Timestamp string `json:"timestamp"`
}

func (s *Store) GetAgents() ([]Agent, error) {

	rows, err := s.DB.Query(`
		SELECT agent_id, hostname, status, last_seen
		FROM agents
	`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var agents []Agent

	for rows.Next() {

		var a Agent

		err := rows.Scan(
			&a.AgentID,
			&a.Hostname,
			&a.Status,
			&a.LastSeen,
		)

		if err != nil {
			return nil, err
		}

		agents = append(agents, a)
	}

	return agents, nil
}

func (s *Store) GetAlerts() ([]Alert, error) {

	rows, err := s.DB.Query(`
		SELECT 
			id,
			agent_id,
			title,
			severity,
			description,
			timestamp
		FROM alerts
		ORDER BY id DESC
	`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var alerts []Alert

	for rows.Next() {

		var a Alert

		err := rows.Scan(
			&a.ID,
			&a.AgentID,
			&a.Title,
			&a.Severity,
			&a.Description,
			&a.Timestamp,
		)

		if err != nil {
			return nil, err
		}

		alerts = append(alerts, a)
	}

	return alerts, nil
}

func (s *Store) GetEvents() ([]Event, error) {

	rows, err := s.DB.Query(`
		SELECT
			id,
			agent_id,
			hostname,
			event_type,
			payload,
			timestamp
		FROM events
		ORDER BY id DESC
	`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var events []Event

	for rows.Next() {

		var e Event

		err := rows.Scan(
			&e.ID,
			&e.AgentID,
			&e.Hostname,
			&e.EventType,
			&e.Payload,
			&e.Timestamp,
		)

		if err != nil {
			return nil, err
		}

		events = append(events, e)
	}

	return events, nil
}
