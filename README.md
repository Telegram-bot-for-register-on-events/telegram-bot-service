# telegram-bot-service

## Назначение:

### Сервис отвечает за:

- Взаимодействие с Telegram Bot API
- Хранение данных пользователей
- Отображение списка событий
- Регистрацию пользователя на событие

### Функциональные требования

- Обработка команд Telegram-бота (/start и др.)
- Отображение списка событий (через Event-Service)
- Регистрация пользователя на событие
- Хранение информации о пользователях

## Требования к запуску:
 - Docker
 - Git
 - Токен telegram-bot, полученный от @BotFather

### 1. Клонирование репозитория 
`git clone github.com/Telegram-bot-for-register-on-events/telegram-bot-service`
`cd telegram-bot-service`

### 2. Cоздание .env
В корне проекта выполните команду:
`mv .envexample .env`.
В поле `TELEGRAM_BOT_TOKEN` вставьте токен, полученный от @BotFather.

### 3. Запуск микросервиса
В корне проекта выполните команду:
`docker compose up -d --build`
