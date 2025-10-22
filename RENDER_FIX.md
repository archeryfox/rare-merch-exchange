# 🔧 Исправление проблемы с базой данных на Render

## ❌ Проблема
```
Failed to connect to database: dial tcp [::1]:5432: connect: connection refused
```

## ✅ Решение

### 1. Исправлен хост базы данных в `render.yaml`

**Было:**
```yaml
- key: DB_HOST
  value: dpg-d3rvapq4d50c73fnc320-a
```

**Стало:**
```yaml
- key: DB_HOST
  value: dpg-d3rvapq4d50c73fnc320-a.oregon-postgres.render.com
```

### 2. Исправлен хост Redis в `render.yaml`

**Было:**
```yaml
- key: REDIS_HOST
  value: red-d3rvarjuibrs739i7b3g
```

**Стало:**
```yaml
- key: REDIS_HOST
  value: red-d3rvarjuibrs739i7b3g.render.com
```

### 3. Добавлен пароль для Redis

```yaml
- key: REDIS_PASSWORD
  value: ""
```

## 🚀 Следующие шаги

### 1. Перезапустите сервис на Render
- Зайдите в [Render Dashboard](https://dashboard.render.com/static/srv-d3rvc3euk2gs73bulkf0)
- Нажмите "Manual Deploy" или дождитесь автоматического деплоя

### 2. Проверьте подключение к базе данных

Если у вас есть доступ к PostgreSQL клиенту:
```bash
PGPASSWORD=PVkqh10d2SSP3f6u1DSS8dp3IGUeDcrT psql -h dpg-d3rvapq4d50c73fnc320-a.oregon-postgres.render.com -U rare_merch_exchange_db_user rare_merch_exchange_db
```

### 3. Выполните миграции базы данных

Если база данных пустая, нужно выполнить миграции. Создайте скрипт для миграций:

```bash
# Установите goose для миграций
go install github.com/pressly/goose/v3/cmd/goose@latest

# Выполните миграции
goose -dir migrations postgres "host=dpg-d3rvapq4d50c73fnc320-a.oregon-postgres.render.com port=5432 user=rare_merch_exchange_db_user password=PVkqh10d2SSP3f6u1DSS8dp3IGUeDcrT dbname=rare_merch_exchange_db sslmode=require" up
```

## 📋 Проверка конфигурации

### Текущие настройки в `render.yaml`:

```yaml
envVars:
  - key: DB_HOST
    value: dpg-d3rvapq4d50c73fnc320-a.oregon-postgres.render.com
  - key: DB_PORT
    value: 5432
  - key: DB_USER
    value: rare_merch_exchange_db_user
  - key: DB_PASSWORD
    value: PVkqh10d2SSP3f6u1DSS8dp3IGUeDcrT
  - key: DB_NAME
    value: rare_merch_exchange_db
  - key: DB_SSL_MODE
    value: require
  - key: REDIS_HOST
    value: red-d3rvarjuibrs739i7b3g.render.com
  - key: REDIS_PORT
    value: 6379
  - key: REDIS_PASSWORD
    value: ""
  - key: REDIS_DB
    value: 0
```

## 🔍 Диагностика

### Если проблема остается:

1. **Проверьте логи сервиса** в Render Dashboard
2. **Убедитесь, что база данных активна** в Render Dashboard
3. **Проверьте, что Redis активен** в Render Dashboard
4. **Убедитесь, что все сервисы в одном регионе** (Oregon)

### Возможные дополнительные проблемы:

1. **База данных не создана** - выполните миграции
2. **Неправильные права доступа** - проверьте пользователя БД
3. **Проблемы с сетью** - убедитесь, что сервисы в одной сети

## 📚 Полезные команды

### Проверка подключения к PostgreSQL:
```bash
PGPASSWORD=PVkqh10d2SSP3f6u1DSS8dp3IGUeDcrT psql -h dpg-d3rvapq4d50c73fnc320-a.oregon-postgres.render.com -U rare_merch_exchange_db_user rare_merch_exchange_db -c "\dt"
```

### Проверка подключения к Redis:
```bash
redis-cli -h red-d3rvarjuibrs739i7b3g.render.com -p 6379 ping
```

### Выполнение миграций:
```bash
goose -dir migrations postgres "host=dpg-d3rvapq4d50c73fnc320-a.oregon-postgres.render.com port=5432 user=rare_merch_exchange_db_user password=PVkqh10d2SSP3f6u1DSS8dp3IGUeDcrT dbname=rare_merch_exchange_db sslmode=require" up
```
