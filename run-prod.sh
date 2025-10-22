#!/bin/bash

# Скрипт для запуска с продакшен настройками Render
# Использование: ./run-prod.sh

echo "🚀 Запуск с продакшен настройками Render..."

# Устанавливаем переменные окружения для Go
export PATH=~/go-local/go/bin:$PATH
export GOPATH=~/go
export PATH=$GOPATH/bin:$PATH

# Переходим в директорию проекта
cd "$(dirname "$0")"

# Загружаем продакшен переменные окружения
echo "📝 Загружаем продакшен настройки из prod.env..."
if [ -f "prod.env" ]; then
    export $(cat prod.env | grep -v '^#' | xargs)
    echo "✅ Продакшен настройки загружены"
else
    echo "❌ Файл prod.env не найден!"
    exit 1
fi

# Проверяем подключение к базе данных
echo "🔍 Проверяем подключение к PostgreSQL..."
if command -v psql &> /dev/null; then
    PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -U $DB_USER -d $DB_NAME -c "SELECT version();" > /dev/null 2>&1
    if [ $? -eq 0 ]; then
        echo "✅ Подключение к PostgreSQL успешно"
    else
        echo "⚠️ Не удалось подключиться к PostgreSQL, но продолжаем..."
    fi
else
    echo "⚠️ psql не установлен, пропускаем проверку БД"
fi

# Проверяем подключение к Redis
echo "🔍 Проверяем подключение к Redis..."
if command -v redis-cli &> /dev/null; then
    # Для Render Redis нужен TLS и пароль
    REDISCLI_AUTH=$REDIS_PASSWORD redis-cli -h $REDIS_HOST -p $REDIS_PORT --tls ping > /dev/null 2>&1
    if [ $? -eq 0 ]; then
        echo "✅ Подключение к Redis успешно"
    else
        echo "⚠️ Не удалось подключиться к Redis, но продолжаем..."
        echo "💡 Убедитесь, что Redis требует TLS подключение"
    fi
else
    echo "⚠️ redis-cli не установлен, пропускаем проверку Redis"
fi

# Загружаем зависимости
echo "📦 Загружаем зависимости..."
go mod tidy

# Запускаем приложение
echo "🎯 Запускаем приложение с продакшен настройками..."
echo "🌐 API будет доступен по адресу: http://localhost:$SERVER_PORT"
echo "📚 Swagger UI: http://localhost:$SERVER_PORT/swagger/index.html"
echo ""
echo "🔧 Настройки:"
echo "   - Сервер: $SERVER_HOST:$SERVER_PORT"
echo "   - База данных: $DB_HOST:$DB_PORT"
echo "   - Redis: $REDIS_HOST:$REDIS_PORT"
echo "   - Режим: $GIN_MODE"
echo ""
echo "Для остановки нажмите Ctrl+C"
echo ""

go run cmd/api/main.go
