package core

import (
	"sync"
)

// ToolStats stores posterior efficiency statistics for a tool action
type ToolStats struct {
	UsageCount      int
	SuccessCount    int
	LootCount       int
	EfficiencyScore float64
}

// KnowledgeBase represents the agent's environment state model and long-term memory
type KnowledgeBase struct {
	Mu               sync.RWMutex
	Target           string
	DiscoveryVectors []string
	CustomPayloads   []string
	Loot             []string
	Observations     []string
	ToolPerformance  map[string]*ToolStats
	CurrentPhase     string

	// State variables for auto-remediation and compliance
	PatchApplied    bool // Indicates if remediation was applied
	PatchVerified   bool // Indicates if remediation was verified via a second test
	ReportGenerated bool // Indicates if regulatory report was generated
}

// NewKnowledgeBase constructs a new state database
func NewKnowledgeBase(target string) *KnowledgeBase {
	return &KnowledgeBase{
		Target:          target,
		ToolPerformance: make(map[string]*ToolStats),
		CurrentPhase:    "Reconnaissance",
	}
}

// UpdateStats updates utility values for a tool run
func (kb *KnowledgeBase) UpdateStats(toolID string, success bool, loot bool) {
	kb.Mu.Lock()
	defer kb.Mu.Unlock()

	stats, exists := kb.ToolPerformance[toolID]
	if !exists {
		stats = &ToolStats{}
		kb.ToolPerformance[toolID] = stats
	}

	stats.UsageCount++
	if success {
		stats.SuccessCount++
	}
	if loot {
		stats.LootCount++
	}

	stats.EfficiencyScore = float64(stats.SuccessCount) / float64(stats.UsageCount)
}
