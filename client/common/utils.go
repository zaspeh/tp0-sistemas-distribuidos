package common

import "strings"

func countLines(s string) int {
	if strings.TrimSpace(s) == "" {
		return 0
	}

	lines := strings.Split(s, "\n")

	count := 0
	for _, l := range lines {
		if strings.TrimSpace(l) != "" {
			count++
		}
	}

	return count
}

func log_send_error(clientID string, err error) {
	log.Errorf(
		"action: send_message | result: fail | client_id: %v | error: %v",
		clientID,
		err,
	)
}
