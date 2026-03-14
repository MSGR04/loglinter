#!/bin/bash

echo "=== 01_basic: Базовые проверки ==="
cd 01_basic && loglinter check_rules.go && cd ..

echo -e "\n=== 02_config: Тестирование конфигурации ==="
cd 02_config && loglinter config_test.go && cd ..

echo -e "\n=== 03_custom_patterns: Кастомные паттерны ==="
cd 03_custom_patterns && loglinter custom_patterns.go && cd ..

echo -e "\n=== 04_regex_patterns: Regex паттерны ==="
cd 04_regex_patterns && loglinter regex_patterns.go && cd ..

echo -e "\n=== 05_fix: Авто-исправление ==="
cd 05_fix && loglinter -fix fix_test.go && cd ..

echo -e "\n=== 06_nested: Поиск конфига ==="
cd 06_nested/cmd/app && loglinter main.go && cd ../../..

