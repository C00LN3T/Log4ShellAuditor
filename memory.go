package main

import (
	"sync"
)

// ============================================================================
// НАУЧНОЕ ОБОСНОВАНИЕ: Расширение Пространства Состояний (Extended State Space)
// 
// Для симуляции полного цикла авто-ремедиации (Self-Healing / Auto-Remediation)
// в пространство состояний S добавлены новые переменные:
// - PatchApplied (Устранение уязвимости на объекте воздействия)
// - PatchVerified (Верификация отсутствия уязвимости повторным зондированием)
// - ReportGenerated (Формирование формального отчета)
//
// Это замыкает цикл DevSecOps внутри одного рефлексивного агента.
// ============================================================================

// ToolStats хранит апостериорную статистику эффективности применения действия
type ToolStats struct {
	UsageCount       int
	SuccessCount     int
	LootCount        int
	EfficiencyScore  float64
}

// KnowledgeBase представляет расширенную модель среды агента
type KnowledgeBase struct {
	Mu               sync.RWMutex
	Target           string
	DiscoveryVectors []string
	CustomPayloads   []string
	Loot             []string
	Observations     []string
	ToolPerformance  map[string]*ToolStats
	CurrentPhase     string

	// Расширение состояния для авто-патчинга и комплаенса
	PatchApplied     bool   // Указывает, применен ли патч
	PatchVerified    bool   // Указывает, проверен ли патч повторной атакой
	ReportGenerated  bool   // Указывает, сгенерирован ли отчет по стандартам
}

func NewKnowledgeBase(target string) *KnowledgeBase {
	return &KnowledgeBase{
		Target:          target,
		ToolPerformance: make(map[string]*ToolStats),
		CurrentPhase:    "Reconnaissance",
	}
}

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
