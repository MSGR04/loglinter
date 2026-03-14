# 🔍 LogLinter

Линтер для проверки лог-записей в Go. Анализирует логи на соответствие правилам форматирования и безопасности.

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![CI](https://github.com/MSGR04/loglinter/actions/workflows/ci.yml/badge.svg)](https://github.com/MSGR04/loglinter/actions/workflows/ci.yml)

---

## ✨ Возможности

- ✅ **Строчная буква** - проверка первой буквы лог-сообщения
- ✅ **Английский язык** - только английские символы в сообщениях
- ✅ **Нет спецсимволов** - запрет на спецсимволы и эмодзи
- ✅ **Нет чувствительных данных** - поиск паролей, токенов, ключей
- ✅ **Конфигурация** - через JSON файл `.loglinter.json`
- ✅ **Кастомные паттерны** - свои правила для чувствительных данных
- ✅ **Авто-исправление** - флаг `-fix` для автоматического исправления
- ✅ **golangci-lint** - интеграция как плагин
- ✅ **Поддержка логгеров** - `log`, `slog`, `zap`

---

## 📁 Структура проекта
```bash
loglinter/
├── cmd/
│ └── loglinter/ # точка входа
├── pkg/
│ └── analyzer/ # основная логика
│ ├── analyzer.go
│ ├── analyzer_test.go
│ ├── config.go
│ ├── rules.go
│ └── testdata/ # тестовые файлы
├── plugin/
│ └── main.go # плагин для golangci-lint
├── examples/ # примеры использования
├── .golangci.yml
└── README.md
```
---
## 📦 Установка

### Вариант 1: Установка как отдельного инструмента

```bash
go install github.com/MSGR04/loglinter.git/cmd/loglinter@latest
```

### Вариант 2: Сборка из исходников

```bash
git clone https://github.com/MSGR04/loglinter.git
cd loglinter
go install ./cmd/loglinter
```
### Вариант 3: Как плагин для golangci-lint

## Сборка плагина (в Linux/WSL)

```bash
go build -buildmode=plugin -o loglinter.so plugin/main.go
```

Добавьте в .golangci.yml:
```yaml
linters-settings:
  custom:
    loglinter:
      path: ./loglinter.so
      description: Линтер для проверки лог-записей
      original-url: github.com/MSGR04/loglinter.git

linters:
  enable:
    - loglinter
```

--- 

## 🚀 Использование

### Базовый запуск

```bash
# Проверить текущую директорию
loglinter ./...

# Проверить конкретный файл
loglinter main.go

# Проверить с авто-исправлением
loglinter -fix ./...
```

### Пример кода

```go
package main

import "log"

func main() {
    // ❌ ОШИБКИ (будут обнаружены)
    log.Println("Starting server")              // заглавная буква
    log.Println("запуск сервера")                // не английский
    log.Println("server started!🚀")             // спецсимволы
    log.Println("user password: secret")        // чувствительные данные
    
    // ✅ ПРАВИЛЬНО
    log.Println("starting server")              // строчная буква
    log.Println("server started")               // без спецсимволов
    log.Println("user authenticated")           // без чувствительных данных
}
```

### Результат выполнения

```bash
$ loglinter main.go
main.go:7:2: лог-сообщение должно начинаться со строчной буквы: "Starting server"
main.go:8:2: лог-сообщение должно содержать только английские символы: "запуск сервера"
main.go:9:2: лог-сообщение не должно содержать спецсимволы или эмодзи: "server started!🚀"
main.go:10:2: лог-сообщение не должно содержать потенциально чувствительные данные: "user password: secret"
```

---
### ⚙️ Конфигурация

Создайте файл .loglinter.json в корне вашего проекта:

```json
{
  "enable_lowercase": true,
  "enable_english_only": true,
  "enable_special_chars": true,
  "enable_sensitive_data": true,
  "use_default_patterns": true,
  "sensitive_patterns": [
    "confidential",
    "secret_.*",
    "internal_[0-9]+"
  ],
  "log_packages": ["log", "slog", "zap"]
}
```

### Опции конфигурации

| Параметр | Описание | По умолчанию |
|----------|----------|--------------|
| `enable_lowercase` | Проверять строчную букву | `true` |
| `enable_english_only` | Проверять английский язык | `true` |
| `enable_special_chars` | Проверять спецсимволы | `true` |
| `enable_sensitive_data` | Проверять чувствительные данные | `true` |
| `use_default_patterns` | Использовать стандартные паттерны | `true` |
| `sensitive_patterns` | Кастомные паттерны | `[]` |
| `log_packages` | Поддерживаемые пакеты логирования | `["log", "slog", "zap"]` |

### Поиск конфига

Линтер ищет .loglinter.json в текущей директории и во всех родительских. Это позволяет иметь один конфиг для всего проекта.

---
## 🔧 Авто-исправление


```bash
# Исправить все автоматически исправимые ошибки
loglinter -fix ./...
```

### Что исправляется

| Тип ошибки | Было | Стало |
|------------|------|-------|
| Заглавная буква | `"Starting server"` | `"starting server"` |
| Спецсимволы | `"Hello!!! World??? 😊"` | `"Hello World"` |

---
## 🎯 Кастомные паттерны


```json
{
  "enable_sensitive_data": true,
  "use_default_patterns": false,
  "sensitive_patterns": [
    "confidential",
    "project_[a-z]+",
    "api_key_[0-9a-f]+"
  ]
}
```

Поддерживаются регулярные выражения (синтаксис Go regexp).

### Примеры

| Паттерн | Найдет | Не найдет |
|---------|--------|-----------|
| `"secret_.*"` | `"secret_project"`, `"secret_data"` | `"mysecret"` |
| `"internal_[0-9]+"` | `"internal_123"` | `"internal_abc"` |
| `"confidential"` | `"confidential report"` | `"public info"` |

---
## 🔄 Интеграция с CI/CD


### GitHub Actions

Создайте файл .github/workflows/ci.yml:

```yaml
name: Lint

on:
  push:
    branches: [ main, master ]
  pull_request:
    branches: [ main, master ]

jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.22'
      - name: Install loglinter
        run: go install github.com/MSGR04/loglinter.git/cmd/loglinter@latest
      - name: Run linter
        run: loglinter ./...
```

### GitLab CI

Создайте файл .gitlab-ci.yml:

```yaml
stages:
  - lint

lint:
  stage: lint
  image: golang:1.22
  script:
    - go install github.com/MSGR04/loglinter.git/cmd/loglinter@latest
    - loglinter ./...
  only:
    - main
    - merge_requests
```

---
## 📚 Примеры

Все примеры находятся в папке [`examples/`](examples/):

| Папка | Описание |
|-------|----------|
| [`01_basic/`](examples/01_basic) | Базовые проверки всех правил |
| [`02_config/`](examples/02_config) | Тестирование конфигурации |
| [`03_custom_patterns/`](examples/03_custom_patterns) | Кастомные паттерны чувствительных данных |
| [`04_regex_patterns/`](examples/04_regex_patterns) | Regex паттерны |
| [`05_fix/`](examples/05_fix) | Авто-исправление |
| [`06_nested/`](examples/06_nested) | Поиск конфига в родительской папке |

Запуск примеров:
```bash
cd examples/01_basic && loglinter check_rules.go
cd ../02_config && loglinter config_test.go
# и так далее..
```
---
## 🧪 Тестирование

```bash
# Запустить тесты линтера
go test ./pkg/analyzer/...

# Проверить на примерах
loglinter examples/...

# Проверить с авто-исправлением
loglinter -fix examples/fix_test.go
```

## 📋 Требования

- Go 1.22 или выше

## ⭐ Поддержка проекта

Автор: MSGR04
Если проект оказался полезным, поставьте звезду на GitHub!
