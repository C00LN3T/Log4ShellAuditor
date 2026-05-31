package main

import (
	"fmt"
	"os"
	"time"
)

// ============================================================================
// Точка входа в демонстрационный стенд реактивного исполнительного агента.
//
// В докер-окружении: ожидает готовности контейнера vulnerable-app,
// запускает HTTP-сервер для отдачи Exploit.class, LDAP-сервер перенаправления,
// и запускает автономный цикл аудита.
//
// В локальном окружении: запускает ЛОКАЛЬНОЕ Java Spring Boot приложение,
// уязвимое к Log4Shell, поднимает встроенные слушатели и запускает Go-агента.
// ============================================================================

func main() {
	fmt.Println("=== ДЕМОНСТРАЦИОННЫЙ СТЕНД РЕАКТИВНОГО АГЕНТА (JAVA SPRING TARGET) ===")

	// Очищаем старые артефакты перед новым запуском
	_ = os.Remove("remediation.properties")
	_ = os.Remove("cve_2021_44228_report.md")
	_ = os.Remove("java_target.log")
	_ = os.RemoveAll("reports")

	// 1. Запуск LDAP TCP-слушателя на порту :1389 для фиксации и перенаправления JNDI
	StartLDAPCallbackListener()
	defer StopLDAPCallbackListener()

	// 2. Запуск HTTP-сервера на порту :8000 для раздачи Exploit.class и приема Loot
	StartHTTPServer()
	defer StopHTTPServer()

	// 3. Управление запуском Java-цели
	agentHost := os.Getenv("AGENT_HOST")
	if agentHost == "" {
		// Локальный запуск на хосте
		_ = os.WriteFile("flag.txt", []byte("FLAG{LOCAL_HOST_LOG4SHELL_SECRET_2026}\n"), 0644)
		defer os.Remove("flag.txt")

		fmt.Println("[*] Запуск скомпилированного уязвимого Java Spring приложения локально...")
		err := StartJavaApp()
		if err != nil {
			fmt.Printf("[!] КРИТИЧЕСКАЯ ОШИБКА: Не удалось запустить Java-приложение: %v\n", err)
			return
		}
		defer StopJavaApp()

		fmt.Println("[*] Ожидание инициализации веб-контекста Spring (3.5 сек)...")
		time.Sleep(3500 * time.Millisecond)
	} else {
		// Запуск внутри контейнера Docker
		fmt.Printf("[*] Обнаружена контейнерная среда (AGENT_HOST=%s).\n", agentHost)
		fmt.Println("[*] Ожидание инициализации сетевого контекста Java-приложения (5 сек)...")
		time.Sleep(5000 * time.Millisecond)
	}

	targetURL := os.Getenv("TARGET_URL")
	if targetURL == "" {
		targetURL = "http://localhost:8080"
	}

	// 4. Инициализация и запуск агента-аудитора
	agent := NewStandaloneExecutor(targetURL)
	agent.Run()
}
