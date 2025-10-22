# 🔧 Настройка переменных окружения

## 📁 Файлы конфигурации

- `env.example` - пример файла с переменными окружения
- `.env` - ваш локальный файл с переменными (создан из env.example)

## 🚀 Быстрый старт

1. **Скопируйте файл примера:**
   ```bash
   cp env.example .env
   ```

2. **Отредактируйте `.env` под ваши нужды:**
   ```bash
   nano .env
   # или
   vim .env
   ```

3. **Запустите приложение:**
   ```bash
   export PATH=~/go-local/go/bin:$PATH
   go run cmd/api/main.go
   ```

## ⚙️ Основные настройки

### База данных PostgreSQL
```bash
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=rare_merch_exchange
DB_SSL_MODE=disable
```

### Redis
```bash
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
```

### JWT секрет
```bash
# ВАЖНО: Измените в продакшене!
JWT_SECRET=your-super-secret-jwt-key-change-in-production
JWT_EXPIRE_HOURS=24
```

## 🔒 Безопасность

### Для разработки:
- Используйте значения по умолчанию
- JWT секрет можно оставить как есть

### Для продакшена:
```bash
GIN_MODE=release
LOG_LEVEL=warn
DB_SSL_MODE=require
JWT_SECRET=your-very-secure-production-secret-key-here
REDIS_PASSWORD=your-redis-password
```

## 📧 Email настройки (опционально)

Раскомментируйте и настройте для отправки уведомлений:
```bash
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password
SMTP_FROM=noreply@rare-merch-exchange.com
```

## 💳 Платежные системы (опционально)

### Stripe:
```bash
STRIPE_SECRET_KEY=sk_test_...
STRIPE_PUBLISHABLE_KEY=pk_test_...
```

### PayPal:
```bash
PAYPAL_CLIENT_ID=your-paypal-client-id
PAYPAL_CLIENT_SECRET=your-paypal-client-secret
```

## 📊 Мониторинг (опционально)

### Sentry:
```bash
SENTRY_DSN=your-sentry-dsn
```

### New Relic:
```bash
NEW_RELIC_LICENSE_KEY=your-newrelic-license-key
```

## 🐳 Docker

Если используете Docker, создайте `.env.docker`:
```bash
cp env.example .env.docker
# Отредактируйте под Docker окружение
```

## 🔄 Обновление конфигурации

При изменении переменных окружения:
1. Отредактируйте `.env`
2. Перезапустите приложение
3. Конфигурация загрузится автоматически

## 📝 Примечания

- Все переменные имеют значения по умолчанию в `config/config.go`
- Переменные окружения переопределяют значения из `config.yaml`
- Файл `.env` автоматически загружается через `godotenv.Load()`
- Не коммитьте `.env` в git (добавлен в `.gitignore`)
