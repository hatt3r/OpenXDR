package store

type Agent struct {
	AgentID  string `json:"agent_id"`
	Hostname string `json:"hostname"`
	Status   string `json:"status"`
	LastSeen string `json:"last_seen"`
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
