package effectors

import (
	"os"
)

// getAgentHost returns the agent's callback host
func getAgentHost() string {
	val := os.Getenv("AGENT_HOST")
	if val != "" {
		return val
	}
	return "127.0.0.1"
}
