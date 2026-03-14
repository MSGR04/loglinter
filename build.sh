#!/bin/bash

# Сборка плагина для golangci-lint
go build -buildmode=plugin -o loglinter.so plugin/main.go

echo "Плагин собран: loglinter.so"