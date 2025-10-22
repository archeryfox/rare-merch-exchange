#!/bin/bash

# Скрипт для установки Docker Compose в локальную директорию
# Использование: ./install-docker-compose.sh

echo "📦 Установка Docker Compose в локальную директорию..."

# Создаем директорию для локальных бинарников
mkdir -p ~/.local/bin

# Скачиваем Docker Compose
echo "⬇️ Скачиваем Docker Compose..."
curl -L "https://github.com/docker/compose/releases/download/v2.20.0/docker-compose-$(uname -s)-$(uname -m)" -o ~/.local/bin/docker-compose

# Делаем исполняемым
chmod +x ~/.local/bin/docker-compose

# Добавляем в PATH
echo "🔧 Добавляем в PATH..."
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc

# Проверяем установку
echo "✅ Docker Compose установлен в ~/.local/bin/docker-compose"
echo "🔄 Перезайдите в терминал или выполните: source ~/.bashrc"

# Проверяем версию
~/.local/bin/docker-compose --version
