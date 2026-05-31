package main

import (
	"fmt"
	"strings"
	"time"
)

// ============================================================================
// НАУЧНОЕ ОБОСНОВАНИЕ: Расширенная политика принятия решений (Extended Decision Policy)
// 
// Жизненный цикл агента теперь охватывает полную петлю управления инцидентами безопасности:
// 1. Оценка защищенности (Reconnaissance + Discovery)
// 2. Тестирование на проникновение (Evasion + Exploitation)
// 3. Автоматическое реагирование (Remediation)
// 4. Повторное тестирование контроля защищенности (Verification)
// 5. Обеспечение соответствия требованиям регуляторов (Reporting)
//
// Отношение принятия решений Think() теперь включает расширенные рефлексивные ветвления,
// позволяющие агенту изменять тактику в зависимости от статуса применения патча.
// ============================================================================

// StandaloneExecutor - Встраиваемый автономный рефлексивный агент
type StandaloneExecutor struct {
	Memory   *KnowledgeBase
	Tools    map[string]Tool
	LastTool string
}

// NewStandaloneExecutor конструирует экземпляр агента с привязкой к цели и набору инструментов
func NewStandaloneExecutor(target string) *StandaloneExecutor {
	exec := &StandaloneExecutor{
		Memory: NewKnowledgeBase(target),
		Tools:  make(map[string]Tool),
	}

	// Инициализация реестра эффекторов (инструментов)
	exec.Tools["port_scanner"] = &ToolPortScanner{}
	exec.Tools["discovery"] = &ToolDiscovery{}
	exec.Tools["payload_generator"] = &ToolPayloadGenerator{}
	exec.Tools["prober"] = &ToolProber{}
	exec.Tools["semantic_fuzzer"] = &ToolSemanticFuzzer{}
	exec.Tools["remediator"] = &ToolRemediator{}
	exec.Tools["reporter"] = &ToolReporter{}

	return exec
}

// Think реализует расширенное решающее правило выбора действия (Policy Mapping)
func (ae *StandaloneExecutor) Think() string {
	ae.Memory.Mu.RLock()
	defer ae.Memory.Mu.RUnlock()

	vectorsCount := len(ae.Memory.DiscoveryVectors)
	payloadsCount := len(ae.Memory.CustomPayloads)
	lootCount := len(ae.Memory.Loot)

	// Шаг 1: Разведка
	if _, scanned := ae.Memory.ToolPerformance["port_scanner"]; !scanned {
		return "port_scanner"
	}

	// Шаг 2: Маппинг
	if vectorsCount == 0 {
		return "discovery"
	}

	// Шаг 3: Синтез воздействия
	if payloadsCount == 0 {
		return "payload_generator"
	}

	// Шаг 4: Эксплуатация
	if lootCount == 0 {
		// Адаптивный обход WAF
		if proberStats, tried := ae.Memory.ToolPerformance["prober"]; tried && proberStats.LootCount == 0 {
			if _, fuzzed := ae.Memory.ToolPerformance["semantic_fuzzer"]; !fuzzed {
				return "semantic_fuzzer"
			}
		}
		return "prober"
	}

	// Шаг 5: Устранение уязвимости (Remediation)
	if !ae.Memory.PatchApplied {
		return "remediator"
	}

	// Шаг 6: Верификация устранения (Повторная атака)
	if !ae.Memory.PatchVerified {
		return "prober" // Prober вызовется в режиме верификации, так как PatchApplied == true
	}

	// Шаг 7: Генерация документации по ГОСТ/ФСТЭК/ФЗ (Reporting)
	if !ae.Memory.ReportGenerated {
		return "reporter"
	}

	// Шаг 8: Цикл завершен
	return "stop"
}

// Run запускает автономный цикл Sense-Think-Act
func (ae *StandaloneExecutor) Run() {
	fmt.Printf("[*] Инициализация агента-исполнителя для цели: %s\n", ae.Memory.Target)
	fmt.Println("================================================================")

	for {
		// 1. THINK (Интеллектуальный вывод / Рассуждение)
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

		// 2. ACT (Выполнение действия эффектором)
		fmt.Printf("[ВЫВОД] Выбран инструмент: '%s' (Текущая фаза: %s)\n", toolID, ae.Memory.CurrentPhase)
		observation := tool.Execute(ae.Memory.Target, ae.Memory)
		
		// 3. SENSE (Обратная связь / Восприятие сигналов среды)
		success := strings.Contains(observation, "SUCCESS") || strings.Contains(observation, "OBSERVATION")
		lootFound := strings.Contains(observation, "Loot:")
		
		// Обновление локальной функции полезности
		ae.Memory.UpdateStats(toolID, success, lootFound)
		
		ae.Memory.Mu.Lock()
		ae.Memory.Observations = append(ae.Memory.Observations, observation)
		ae.Memory.Mu.Unlock()

		fmt.Printf("[НАБЛЮДЕНИЕ] %s\n\n", observation)
		ae.LastTool = toolID
		
		time.Sleep(800 * time.Millisecond)
	}
}
