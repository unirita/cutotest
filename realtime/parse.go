package realtime

import (
	"encoding/json"
	"testing"
)

type result struct {
	Status  int            `json:"status"`
	Message string         `json:"message"`
	PID     int            `json:"pid"`
	Network *networkResult `json:"network"`
}

type networkResult struct {
	Instance int    `json:"instance"`
	Name     string `json:"name"`
}

func parseRealtimeResult(t *testing.T, output string) *result {
	res := new(result)
	if err := json.Unmarshal([]byte(output), res); err != nil {
		t.Fatalf("Output JSON has error: %s", err)
	}
	return res
}
