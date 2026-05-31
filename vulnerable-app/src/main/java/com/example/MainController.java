package com.example;

import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.springframework.http.ResponseEntity;
import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestHeader;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.PostMapping;

@RestController
public class MainController {
    private static final Logger logger = LogManager.getLogger(MainController.class);

    private boolean isRemediated() {
        String[] paths = {"/app/shared/remediation.properties", "remediation.properties"};
        for (String p : paths) {
            try {
                java.nio.file.Path path = java.nio.file.Paths.get(p);
                if (java.nio.file.Files.exists(path)) {
                    String content = java.nio.file.Files.readString(path);
                    if (content.contains("remediated=true")) {
                        return true;
                    }
                }
            } catch (Exception e) {
                // ignore
            }
        }
        return false;
    }

    @GetMapping("/")
    public ResponseEntity<String> index(
        @RequestHeader(value = "X-Api-Version", defaultValue = "1.0") String apiVersion,
        @RequestParam(value = "input", required = false) String input
    ) {
        // Имитация сигнатурного WAF: блокировка прямого вхождения ${jndi:
        if (apiVersion != null && apiVersion.contains("${jndi:")) {
            System.out.println("[WAF] Запрос заблокирован по сигнатуре: " + apiVersion);
            return ResponseEntity.status(HttpStatus.FORBIDDEN)
                .body("<html><body><h1>403 Forbidden</h1><p>Blocked by WAF (Signature Match)</p></body></html>");
        }

        // Trigger Log4j vulnerable logging if NOT remediated
        if (!isRemediated()) {
            if (input != null && !input.isEmpty()) {
                logger.info("[AUDIT] Input parameter logged: {}", input);
            }
            logger.info("[AUDIT] API Version header logged: {}", apiVersion);
        } else {
            System.out.println("[PATCH] Механизм JNDI Lookup отключен (Remediation применен). Логирование заголовка заблокировано.");
        }
        
        return ResponseEntity.ok("<html><body><h1>Enterprise Portal (Simple)</h1><p>Status: Active</p></body></html>");
    }

    @PostMapping("/remediate")
    public ResponseEntity<String> remediate() {
        try {
            java.nio.file.Path path = java.nio.file.Paths.get("remediation.properties");
            java.nio.file.Files.writeString(path, "remediated=true\n");
            System.out.println("[REMEDIATION HOOK] Успешно создана конфигурация remediation.properties (remediated=true).");
            return ResponseEntity.ok("Remediated successfully");
        } catch (Exception e) {
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body("Error: " + e.getMessage());
        }
    }
}
