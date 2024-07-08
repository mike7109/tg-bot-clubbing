#!/bin/bash
set -e

# Создание директории для сборки
mkdir -p build

# Сборка gRPC приложения
echo "Building gRPC Application..."
go build -o build/app ./cmd/main.go
