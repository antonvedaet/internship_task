# Тестовое задание. Сервис назначения ревьюеров для Pull Request’ов

Микросервис для автоматического назначения ревьюеров на Pull Request'ы внутри команд.

## Возможности

- Управление командами и пользователями
- Автоматическое назначение ревьюеров (до 2-х активных пользователей из команды автора)
- Переназначение ревьюеров из той же команды
- Идемпотентный мерж PR - повторный вызов не приводит к ошибке
- Массовая деактивация пользователей команды
- Получение списка PR'ов назначенных пользователю

## Технологии

- Go 1.24 - язык программирования
- PostgreSQL - база данных
- Docker & Docker Compose - контейнеризация

## Установка и запуск

### Требования
- Docker
- Docker Compose

### Быстрый старт

```bash
# Клонирование репозитория
git clone <repository-url>
cd pr-reviewer-service

# Запуск сервиса
make build
make run

# Проверка здоровья
make health
```

### Ручная сборка

```bash
docker-compose up --build -d
```

## API Endpoints

### Команды
- `POST /team/add` - Создать команду с участниками
- `GET /team/get?team_name=name` - Получить команду с участниками
- `POST /team/deactivate` - Массовая деактивация пользователей команды

### Пользователи
- `POST /users/setIsActive` - Установить флаг активности пользователя
- `GET /users/getReview?user_id=id` - Получить PR'ы пользователя для ревью

### Pull Requests
- `POST /pullRequest/create` - Создать PR и назначить ревьюеров
- `POST /pullRequest/merge` - Пометить PR как MERGED
- `POST /pullRequest/reassign` - Переназначить ревьюера

### Системные
- `GET /health` - Проверка здоровья сервиса

## Структура базы данных

```sql
teams (team_name)
users (user_id, username, team_name, is_active)
pull_requests (pull_request_id, author_id, status, assigned_reviewers[], ...)
```

## Структура

```
cmd/server/
└── main.go              # Точка входа
internal/
├── http/                # Роутинг
|    └──handlers/        # HTTP обработчики
|               
├── service/             # Бизнес-логика
├── database/            # Слой работы с БД
└── models/              # Модели данных
```

## Линтер

- govet - стандартный анализатор Go
- staticcheck - статический анализ кода
- ineffassign - обнаружение неэффективных присваиваний
- unused - поиск неиспользуемого кода

## Форматтер

- gofmt - проверка форматирования
- goimports - сортировка импортов

## Команды Makefile

```bash
make build      # Сборка контейнеров
make run        # Запуск сервиса
make stop       # Остановка сервиса
make clean      # Остановка с удалением volumes
make test-db    # Тест подключения к БД
make logs       # Просмотр логов приложения
make health     # Проверка здоровья сервиса
make lint       # Запуск линтера
```

## Логика

### Назначение ревьюеров
- Автоматически назначаются до 2 активных ревьюеров из команды автора
- Автор исключается из списка кандидатов
- Если кандидатов меньше двух - назначается доступное количество (0/1)

### Переназначение
- Заменяет одного ревьюера на случайного активного участника из команды заменяемого
- После MERGED менять ревьюеров нельзя
- Новый ревьюер должен быть из той же команды и не быть уже назначенным на PR

### Деактивация пользователей
- Массовая деактивация всех пользователей команды
- Не затрагивает уже назначенные PR (только флаг активности)

## Тестирование

Примеры запросов:

```bash
# Создание команды
curl -X POST http://localhost:8080/team/add \
  -H "Content-Type: application/json" \
  -d '{"team_name":"backend","members":[{"user_id":"u1","username":"A","is_active":true}]}'

# Создание PR
curl -X POST http://localhost:8080/pullRequest/create \
  -H "Content-Type: application/json" \
  -d '{"pull_request_id":"pr-1","pull_request_name":"Add bugs","author_id":"u1"}'

# Массовая деактивация
curl -X POST http://localhost:8080/team/deactivate \
  -H "Content-Type: application/json" \
  -d '{"team_name":"backend"}'
```

## Особенности реализации

- Идемпотентность операций - повторные вызовы merge не вызывают ошибок
- Валидация по спецификации - все ошибки соответствуют OpenAPI
- Встроенный роутинг Go 1.24 без внешних зависимостей