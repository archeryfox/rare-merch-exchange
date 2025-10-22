#!/bin/bash

# Скрипт для генерации Swagger документации
# Использование: ./generate-swagger.sh

echo "🚀 Генерация Swagger документации..."

# Устанавливаем переменные окружения для Go
export PATH=~/go-local/go/bin:$PATH
export GOPATH=~/go
export PATH=$GOPATH/bin:$PATH

# Переходим в директорию проекта
cd "$(dirname "$0")"

# Проверяем, что swag установлен
if ! command -v swag &> /dev/null; then
    echo "❌ swag не найден. Устанавливаем..."
    go install github.com/swaggo/swag/cmd/swag@latest
fi

# Генерируем документацию
echo "📝 Генерируем Swagger документацию..."
swag init -g cmd/api/main.go

if [ $? -eq 0 ]; then
    echo "✅ Swagger документация успешно сгенерирована!"
    echo "📁 Файлы созданы в папке docs/:"
    echo "   - docs.go"
    echo "   - swagger.json"
    echo "   - swagger.yaml"
    echo ""
    echo "🌐 Swagger UI будет доступен по адресу:"
    echo "   http://localhost:8080/swagger/index.html"
    echo ""
    echo "💡 Для запуска сервера используйте:"
    echo "   go run cmd/api/main.go"
else
    echo "❌ Ошибка при генерации документации"
    exit 1
fi
