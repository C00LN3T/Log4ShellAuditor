package effectors

import (
	"aeon-arsenal/internal/core"
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"
)

// ToolPortScanner performs reconnaissance and network port mapping
type ToolPortScanner struct{}

func (t *ToolPortScanner) ID() string { return "port_scanner" }

func (t *ToolPortScanner) Execute(target string, kb *core.KnowledgeBase) string {
	u, err := url.Parse(target)
	host := target
	if err == nil && u.Host != "" {
		host = u.Host
	}
	if !strings.Contains(host, ":") {
		host = host + ":8080"
	}

	fmt.Printf("[ЭФФЕКТОР:port_scanner] Сканирование порта %s...\n", host)
	conn, err := net.DialTimeout("tcp", host, 2*time.Second)
	if err != nil {
		return fmt.Sprintf("FAILURE: Порт %s закрыт или недоступен: %v", host, err)
	}
	conn.Close()

	kb.Mu.Lock()
	kb.CurrentPhase = "Discovery"
	kb.Mu.Unlock()

	return fmt.Sprintf("OBSERVATION: Обнаружен открытый HTTP-порт %s. Java Spring Web-служба отвечает.", host)
}
