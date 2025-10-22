#!/bin/bash

# Скрипт для запуска локальной разработки
# Использование: ./dev.sh

echo "🚀 Запуск локальной разработки Rare Merch Exchange..."

# Устанавливаем переменные окружения для Go
export PATH=~/go-local/go/bin:$PATH
export GOPATH=~/go
export PATH=$GOPATH/bin:$PATH

# Переходим в директорию проекта
cd "$(dirname "$0")"

# Проверяем, что Docker установлен
if ! command -v docker &> /dev/null; then
    echo "❌ Docker не установлен. Устанавливаем..."
    echo "📝 Инструкции по установке Docker:"
    echo "   curl -fsSL https://get.docker.com -o get-docker.sh"
    echo "   sudo sh get-docker.sh"
    echo "   sudo usermod -aG docker \$USER"
    echo "   # Перезайдите в систему после установки"
    exit 1
fi

# Проверяем, что Docker Compose установлен
if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose не установлен. Устанавливаем..."
    echo "📝 Инструкции по установке Docker Compose:"
    echo "   sudo curl -L \"https://github.com/docker/compose/releases/download/v2.20.0/docker-compose-\$(uname -s)-\$(uname -m)\" -o /usr/local/bin/docker-compose"
    echo "   sudo chmod +x /usr/local/bin/docker-compose"
    exit 1
fi

# Создаем .env файл для локальной разработки
echo "📝 Создаём .env файл для локальной разработки..."
cat > .env << 'EOF'
# ============================================
# RARE MERCH EXCHANGE - Локальная разработка
# ============================================

# СЕРВЕР
SERVER_PORT=8080
SERVER_HOST=localhost
GIN_MODE=debug

# БАЗА ДАННЫХ (PostgreSQL - Docker)
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=rare_merch_exchange
DB_SSL_MODE=disable

# REDIS (Docker)
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT АУТЕНТИФИКАЦИЯ
JWT_SECRET=your-super-secret-jwt-key-change-in-production
JWT_EXPIRE_HOURS=24

# ЗАГРУЗКА ФАЙЛОВ
UPLOAD_MAX_SIZE=10485760
UPLOAD_PATH=./uploads

# АУКЦИОНЫ
AUCTION_SOFT_CLOSE_MINUTES=5
AUCTION_MIN_INCREMENT_PERCENT=5

# ЛОТЕРЕИ
LOTTERY_MAX_TICKETS_PER_USER=10
LOTTERY_MIN_TICKETS=1

# КОМИССИИ
COMMISSION_RATE_VERIFIED=0.05
COMMISSION_RATE_REGULAR=0.08

# ЭСКРОУ
ESCROW_AUTO_RELEASE_DAYS=7
ESCROW_DISPUTE_PERIOD_DAYS=3

# АНТИФРОД
ANTIFRAUD_MAX_BIDS_PER_MINUTE=10
ANTIFRAUD_SUSPICIOUS_ACTIVITY_THRESHOLD=5

# ЛОГИРОВАНИЕ
LOG_LEVEL=info
LOG_FORMAT=json
EOF

echo "✅ .env файл создан"

# Запускаем Docker контейнеры
echo "🐳 Запускаем PostgreSQL и Redis..."
docker-compose up -d

# Ждем, пока база данных будет готова
echo "⏳ Ждём готовности базы данных..."
sleep 10

# Проверяем подключение к PostgreSQL
echo "🔍 Проверяем подключение к PostgreSQL..."
if docker exec rare-merch-postgres pg_isready -U postgres; then
    echo "✅ PostgreSQL готов"
else
    echo "❌ PostgreSQL не готов. Ждём ещё..."
    sleep 10
fi

# Проверяем подключение к Redis
echo "🔍 Проверяем подключение к Redis..."
if docker exec rare-merch-redis redis-cli ping | grep -q PONG; then
    echo "✅ Redis готов"
else
    echo "❌ Redis не готов. Ждём ещё..."
    sleep 5
fi

# Выполняем миграции
echo "🗄️ Выполняем миграции базы данных..."
if command -v goose &> /dev/null; then
    goose -dir migrations postgres "host=localhost port=5432 user=postgres password=password dbname=rare_merch_exchange sslmode=disable" up
    echo "✅ Миграции выполнены"
else
    echo "⚠️ Goose не установлен. Устанавливаем..."
    go install github.com/pressly/goose/v3/cmd/goose@latest
    goose -dir migrations postgres "host=localhost port=5432 user=postgres password=password dbname=rare_merch_exchange sslmode=disable" up
    echo "✅ Миграции выполнены"
fi

# Загружаем зависимости
echo "📦 Загружаем зависимости..."
go mod tidy

# Запускаем приложение
echo "🎯 Запускаем приложение..."
echo "🌐 API будет доступен по адресу: http://localhost:8080"
echo "📚 Swagger UI: http://localhost:8080/swagger/index.html"
echo ""
echo "Для остановки нажмите Ctrl+C"
echo "Для остановки Docker контейнеров: docker-compose down"
echo ""

go run cmd/api/main.go
