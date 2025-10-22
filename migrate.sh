#!/bin/bash

# Скрипт для выполнения миграций базы данных на Render
# Использование: ./migrate.sh

echo "🗄️ Выполнение миграций базы данных..."

# Устанавливаем переменные окружения для Go
export PATH=~/go-local/go/bin:$PATH
export GOPATH=~/go
export PATH=$GOPATH/bin:$PATH

# Переходим в директорию проекта
cd "$(dirname "$0")"

# Параметры подключения к базе данных
DB_HOST="dpg-d3rvapq4d50c73fnc320-a.oregon-postgres.render.com"
DB_PORT="5432"
DB_USER="rare_merch_exchange_db_user"
DB_PASSWORD="PVkqh10d2SSP3f6u1DSS8dp3IGUeDcrT"
DB_NAME="rare_merch_exchange_db"
DB_SSL_MODE="require"

# Строка подключения
DB_URL="host=${DB_HOST} port=${DB_PORT} user=${DB_USER} password=${DB_PASSWORD} dbname=${DB_NAME} sslmode=${DB_SSL_MODE}"

echo "📡 Подключение к базе данных: ${DB_HOST}"

# Проверяем, что goose установлен
if ! command -v goose &> /dev/null; then
    echo "📦 Устанавливаем goose для миграций..."
    go install github.com/pressly/goose/v3/cmd/goose@latest
fi

# Проверяем подключение к базе данных
echo "🔍 Проверяем подключение к базе данных..."
if ! goose -dir migrations postgres "${DB_URL}" status; then
    echo "❌ Не удалось подключиться к базе данных"
    echo "🔧 Проверьте настройки подключения в RENDER_FIX.md"
    exit 1
fi

echo "✅ Подключение к базе данных успешно!"

# Выполняем миграции
echo "🚀 Выполняем миграции..."
if goose -dir migrations postgres "${DB_URL}" up; then
    echo "✅ Миграции выполнены успешно!"
    
    # Показываем статус миграций
    echo "📊 Статус миграций:"
    goose -dir migrations postgres "${DB_URL}" status
    
    echo ""
    echo "🎉 База данных готова к использованию!"
    echo "🌐 Теперь можно перезапустить сервис на Render"
else
    echo "❌ Ошибка при выполнении миграций"
    exit 1
fi
