#!/bin/bash

# Скрипт для запуска Rare Merch Exchange API
# Использование: ./run.sh

echo "🚀 Запуск Rare Merch Exchange API..."

# Устанавливаем переменные окружения для Go
export PATH=~/go-local/go/bin:$PATH
export GOPATH=~/go
export PATH=$GOPATH/bin:$PATH

# Переходим в директорию проекта
cd "$(dirname "$0")"

# Проверяем, что .env файл существует
if [ ! -f ".env" ]; then
    echo "❌ Файл .env не найден!"
    echo "📝 Создаём .env из примера..."
    cp env.example .env
    echo "✅ Файл .env создан. Отредактируйте его при необходимости."
fi

# Проверяем, что Go установлен
if ! command -v go &> /dev/null; then
    echo "❌ Go не найден. Устанавливаем..."
    curl -L https://go.dev/dl/go1.23.0.linux-amd64.tar.gz -o /tmp/go.tar.gz
    mkdir -p ~/go-local && tar -C ~/go-local -xzf /tmp/go.tar.gz
    export PATH=~/go-local/go/bin:$PATH
fi

# Загружаем зависимости
echo "📦 Загружаем зависимости..."
go mod tidy

# Запускаем приложение
echo "🎯 Запускаем сервер..."
echo "🌐 API будет доступен по адресу: http://localhost:8080"
echo "📚 Swagger UI: http://localhost:8080/swagger/index.html"
echo ""
echo "Для остановки нажмите Ctrl+C"
echo ""

go run cmd/api/main.go
