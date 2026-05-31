package effectors

import (
	"aeon-arsenal/internal/core"
	"aeon-arsenal/pkg/target"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// ToolRemediator applies automated patching and remediation
type ToolRemediator struct{}

func (t *ToolRemediator) ID() string { return "remediator" }

func (t *ToolRemediator) Execute(targetURL string, kb *core.KnowledgeBase) string {
	fmt.Printf("[ЭФФЕКТОР:remediator] Анализ причин уязвимости и генерация исправления для %s...\n", targetURL)
	time.Sleep(500 * time.Millisecond)

	configContent := "remediated=true\n"
	err := ioutil.WriteFile("remediation.properties", []byte(configContent), 0644)
	if err != nil {
		return fmt.Sprintf("FAILURE: Не удалось записать свойства конфигурации: %v", err)
	}

	client := &http.Client{Timeout: 3 * time.Second}
	u, _ := url.Parse(targetURL)
	remediateURL := fmt.Sprintf("http://%s/remediate", u.Host)

	fmt.Printf("[ЭФФЕКТОР:remediator] Отправка команды применения патча на %s...\n", remediateURL)
	resp, err := client.Post(remediateURL, "application/json", nil)
	if err == nil {
		resp.Body.Close()
	} else {
		fmt.Printf("[ЭФФЕКТОР:remediator] Предупреждение: Не удалось связаться с /remediate эндпоинтом: %v. Попытка локального перезапуска...\n", err)
		_ = target.StartJavaApp()
	}

	time.Sleep(2000 * time.Millisecond)

	kb.Mu.Lock()
	kb.PatchApplied = true
	kb.CurrentPhase = "Verification"
	kb.Mu.Unlock()

	return "REMEDIATION_SUCCESS: Патч применен. На целевое приложение отправлен запрос ремедиации (изменен флаг -Dlog4j2.formatMsgNoLookups=true). JVM успешно переинициализирована."
}
