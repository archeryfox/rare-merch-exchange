# 🔧 Исправление проекта для локальной разработки

## ❌ Проблема:
```
"Failed to connect to database","error":"failed to ping database: dial tcp 127.0.0.1:5432: connect: connection refused"
```

## ✅ Решение:

### 1. Установите Docker и Docker Compose

```bash
# Установка Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER

# Установка Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/download/v2.20.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Перезайдите в систему после установки
```

### 2. Запустите локальную разработку

```bash
./dev.sh
```

Этот скрипт:
- Создаст .env файл с правильными настройками
- Запустит PostgreSQL и Redis в Docker
- Выполнит миграции базы данных
- Запустит приложение

### 3. Альтернативный способ (без Docker)

Если Docker недоступен, установите PostgreSQL и Redis локально:

```bash
# Установка PostgreSQL
sudo apt update
sudo apt install postgresql postgresql-contrib

# Установка Redis
sudo apt install redis-server

# Запуск сервисов
sudo systemctl start postgresql
sudo systemctl start redis-server
sudo systemctl enable postgresql
sudo systemctl enable redis-server

# Создание базы данных
sudo -u postgres createdb rare_merch_exchange
sudo -u postgres psql -c "ALTER USER postgres PASSWORD 'password';"

# Выполнение миграций
go install github.com/pressly/goose/v3/cmd/goose@latest
goose -dir migrations postgres "host=localhost port=5432 user=postgres password=password dbname=rare_merch_exchange sslmode=disable" up

# Запуск приложения
go run cmd/api/main.go
```

## 📁 Созданные файлы:

- `docker-compose.yml` - конфигурация Docker для PostgreSQL и Redis
- `dev.sh` - скрипт для запуска локальной разработки
- `.env` - переменные окружения (создается автоматически)

## 🌐 После запуска:

- **API:** http://localhost:8080/api/v1
- **Swagger UI:** http://localhost:8080/swagger/index.html

## 🛑 Остановка:

```bash
# Остановка Docker контейнеров
docker-compose down

# Остановка локальных сервисов
sudo systemctl stop postgresql redis-server
```

## 🔍 Проверка:

```bash
# Проверка PostgreSQL
psql -h localhost -U postgres -d rare_merch_exchange -c "\dt"

# Проверка Redis
redis-cli ping
```
