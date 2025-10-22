# Настройка Swagger документации

## Установка необходимых инструментов

### 1. Установка Go (если не установлен)

```bash
# Через snap (рекомендуется)
sudo snap install go --classic

# Или через apt
sudo apt update
sudo apt install golang-go

# Проверка установки
go version
```

### 2. Установка swag CLI

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

## Генерация Swagger документации

### 1. Генерация документации

Из корневой директории проекта выполните:

```bash
swag init -g cmd/api/main.go
```

Эта команда:
- Сканирует код в поиске Swagger аннотаций
- Генерирует файлы `docs.go`, `swagger.json` и `swagger.yaml` в папке `docs/`
- Обновляет существующий файл `docs/docs.go`

### 2. Проверка генерации

После выполнения команды в папке `docs/` должны появиться файлы:
- `docs.go` - Go код с документацией
- `swagger.json` - JSON схема API
- `swagger.yaml` - YAML схема API

## Доступ к Swagger UI

После запуска сервера Swagger UI будет доступен по адресу:
```
http://localhost:8080/swagger/index.html
```

## Структура Swagger аннотаций

### Основные аннотации в main.go

```go
// @title Rare Merch Exchange API
// @version 1.0
// @description API для биржи раритетного мерча с аукционами, лотереями и конкурсами
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
```

### Аннотации для методов handlers

```go
// MethodName описание метода
// @Summary Краткое описание
// @Tags группа методов (auth, users, items, auctions, etc.)
// @Accept json
// @Produce json
// @Param name type location required "описание"
// @Success status {type} model "описание"
// @Failure status {type} model "описание"
// @Security BearerAuth (если требуется аутентификация)
// @Router /path [method]
func (h *Handler) MethodName(c *gin.Context) {
    // реализация
}
```

### Примеры аннотаций

#### GET запрос с параметрами
```go
// @Summary Поиск товаров
// @Tags items
// @Accept json
// @Produce json
// @Param query query string false "Поисковый запрос"
// @Param limit query int false "Лимит результатов" default(20)
// @Success 200 {array} item.Item
// @Router /items [get]
```

#### POST запрос с телом
```go
// @Summary Создать товар
// @Tags items
// @Accept json
// @Produce json
// @Param request body item.CreateItemRequest true "Данные товара"
// @Success 201 {object} item.Item
// @Failure 400 {object} ErrorResponse
// @Security BearerAuth
// @Router /items [post]
```

#### GET запрос с path параметром
```go
// @Summary Получить товар по ID
// @Tags items
// @Accept json
// @Produce json
// @Param id path string true "ID товара"
// @Success 200 {object} item.Item
// @Failure 404 {object} ErrorResponse
// @Router /items/{id} [get]
```

## Обновление документации

При изменении аннотаций в коде необходимо повторно выполнить:

```bash
swag init -g cmd/api/main.go
```

## Полезные команды

### Генерация только JSON
```bash
swag init -g cmd/api/main.go --parseDependency --parseInternal
```

### Генерация с кастомным выводом
```bash
swag init -g cmd/api/main.go -o ./docs --parseDependency --parseInternal
```

## Структура ответов

В проекте используются стандартные структуры ответов:

### ErrorResponse
```go
type ErrorResponse struct {
    Error   string `json:"error"`
    Message string `json:"message"`
}
```

### SuccessResponse
```go
type SuccessResponse struct {
    Message string `json:"message"`
}
```

## Теги для группировки методов

- `auth` - аутентификация (регистрация, вход)
- `users` - пользователи
- `items` - товары
- `auctions` - аукционы
- `lotteries` - лотереи
- `contests` - конкурсы
- `transactions` - транзакции
- `verification` - верификация

## Безопасность

Для методов, требующих аутентификации, используется:
```go
// @Security BearerAuth
```

Это соответствует определению в main.go:
```go
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
```
