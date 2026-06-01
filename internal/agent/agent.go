package agent

import (
	"auto-audit/internal/core"
	"auto-audit/internal/effectors"
	"fmt"
	"strings"
	"time"
)

// StandaloneExecutor implements the autonomous agent's cognitive loop (Sense-Think-Act)
type StandaloneExecutor struct {
	Memory   *core.KnowledgeBase
	Tools    map[string]core.Tool
	LastTool string
}

// NewStandaloneExecutor instantiates the agent, binds it to the state database, and registers effectors
func NewStandaloneExecutor(target string) *StandaloneExecutor {
	exec := &StandaloneExecutor{
		Memory: core.NewKnowledgeBase(target),
		Tools:  make(map[string]core.Tool),
	}

	exec.Tools["port_scanner"] = &effectors.ToolPortScanner{}
	exec.Tools["discovery"] = &effectors.ToolDiscovery{}
	exec.Tools["payload_generator"] = &effectors.ToolPayloadGenerator{}
	exec.Tools["prober"] = &effectors.ToolProber{}
	exec.Tools["semantic_fuzzer"] = &effectors.ToolSemanticFuzzer{}
	exec.Tools["remediator"] = &effectors.ToolRemediator{}
	exec.Tools["reporter"] = &effectors.ToolReporter{}

	return exec
}

// Think implements POMDP utility decision model mapped to discrete security phases
func (ae *StandaloneExecutor) Think() string {
	ae.Memory.Mu.RLock()
	defer ae.Memory.Mu.RUnlock()

	vectorsCount := len(ae.Memory.DiscoveryVectors)
	payloadsCount := len(ae.Memory.CustomPayloads)
	lootCount := len(ae.Memory.Loot)

	// Step 1: Reconnaissance
	if _, scanned := ae.Memory.ToolPerformance["port_scanner"]; !scanned {
		return "port_scanner"
	}

	// Step 2: Mapping
	if vectorsCount == 0 {
		return "discovery"
	}

	// Step 3: Signature Synthesis
	if payloadsCount == 0 {
		return "payload_generator"
	}

	// Step 4: Exploitation
	if lootCount == 0 {
		if proberStats, tried := ae.Memory.ToolPerformance["prober"]; tried && proberStats.LootCount == 0 {
			if _, fuzzed := ae.Memory.ToolPerformance["semantic_fuzzer"]; !fuzzed {
				return "semantic_fuzzer"
			}
		}
		return "prober"
	}

	// Step 5: Active Remediation
	if !ae.Memory.PatchApplied {
		return "remediator"
	}

	// Step 6: Patch Verification
	if !ae.Memory.PatchVerified {
		return "prober" // Calls prober in verification mode (PatchApplied == true)
	}

	// Step 7: Compliance Documentation
	if !ae.Memory.ReportGenerated {
		return "reporter"
	}

	// Step 8: Execution complete
	return "stop"
}

// Run executes the Sense-Think-Act loop until stopping condition is met
func (ae *StandaloneExecutor) Run() {
	fmt.Printf("[*] Инициализация агента-исполнителя для цели: %s\n", ae.Memory.Target)
	fmt.Println("================================================================")

	for {
		// 1. THINK
		toolID := ae.Think()
		if toolID == "stop" {
			fmt.Println("================================================================")
			fmt.Println("[INFO] Жизненный цикл аудита, патчинга и комплаенса завершен.")
			break
		}

		tool := ae.Tools[toolID]
		if tool == nil {
			fmt.Printf("[ERROR] Инструмент %s не найден в реестре эффекторов.\n", toolID)
			break
		}

		// 2. ACT
		fmt.Printf("[ВЫВОД] Выбран инструмент: '%s' (Текущая фаза: %s)\n", toolID, ae.Memory.CurrentPhase)
		observation := tool.Execute(ae.Memory.Target, ae.Memory)

		// 3. SENSE
		success := strings.Contains(observation, "SUCCESS") || strings.Contains(observation, "OBSERVATION")
		lootFound := strings.Contains(observation, "Loot:")

		ae.Memory.UpdateStats(toolID, success, lootFound)

		ae.Memory.Mu.Lock()
		ae.Memory.Observations = append(ae.Memory.Observations, observation)
		ae.Memory.Mu.Unlock()

		fmt.Printf("[НАБЛЮДЕНИЕ] %s\n\n", observation)
		ae.LastTool = toolID

		time.Sleep(800 * time.Millisecond)
	}
}
