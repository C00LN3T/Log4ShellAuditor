package effectors

import (
	"auto-audit/internal/core"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// ToolDiscovery maps parameter vectors and inputs
type ToolDiscovery struct{}

func (t *ToolDiscovery) ID() string { return "discovery" }

func (t *ToolDiscovery) Execute(target string, kb *core.KnowledgeBase) string {
	fmt.Printf("[ЭФФЕКТОР:discovery] Исследование структуры веб-приложения %s...\n", target)

	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get(target)
	if err != nil {
		return fmt.Sprintf("FAILURE: Ошибка подключения к цели: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Sprintf("FAILURE: Ошибка чтения ответа: %v", err)
	}

	bodyStr := string(bodyBytes)
	foundVectors := []string{}

	if strings.Contains(bodyStr, "input") {
		foundVectors = append(foundVectors, "/api/v1/query?input=")
	} else {
		foundVectors = append(foundVectors, "/?input=")
	}
	foundVectors = append(foundVectors, "X-Api-Version")

	kb.Mu.Lock()
	kb.DiscoveryVectors = append(kb.DiscoveryVectors, foundVectors...)
	kb.Mu.Unlock()

	return fmt.Sprintf("OBSERVATION: Обнаружены потенциальные векторы ввода: GET-параметр '%s' и HTTP-заголовок '%s'.",
		foundVectors[0], foundVectors[1])
}
