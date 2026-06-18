# robo-ruka

## Макет руки

![Макет робо-руки](assets/maket.png)

## Скриншоты

Включённое состояние:

![Состояние On](assets/on.png)

Выключенное cостояние:

![Состояние Off](assets/off.png)

Выполненных команд:

![Лог команд](assets/commands.png)

## Как запустить

Требуется Go 1.22+.

```bash
git clone git@github.com:AlexeyM0L/robo-ruka.git
cd robo-ruka
cp .env.example .env
go run ./cmd/server
```

По умолчанию сервер слушает `http://localhost:8080`.
При первом запуске рядом автоматически создаётся файл базы данных `robo-ruka.db`.

### Конфигурация

Настройки читаются из переменных окружения (можно положить в `.env`, пример — в `.env.example`):

| Переменная       | По умолчанию       | Назначение                       |
|------------------|--------------------|----------------------------------|
| `HOST`           | `localhost`        | Хост                             |
| `PORT`           | `8080`             | Порт                             |
| `TEMPLATE_PATH`  | `web/index.html`   | Путь до HTML-шаблона             |
| `DB_PATH`        | `robo-ruka.db`     | Путь до файла базы данных SQLite |

### Сборка бинаря

```bash
go build -o server ./cmd/server
./server
```

## База данных

В качестве хранилища использовал **SQLite** — это БД, которая живёт в одном
файле. Драйвер для работы с базой данных на Go — `modernc.org/sqlite`.

Одна таблица с одной строкой:

```sql
CREATE TABLE IF NOT EXISTS status (
    id    INTEGER PRIMARY KEY,
    value TEXT NOT NULL  -- "on" или "off"
);
```

Запись делается через upsert (`INSERT ... ON CONFLICT(id) DO UPDATE`):
если строки ещё нет — она создаётся, если есть — обновляется.


## Архитектура

```
HTTP-запрос
   |
handler   ── разбирает запрос, рендерит HTML
   |
service   ── бизнес-логика и валидация
   |
repository ── работа с базой данных (SQLite)
   |
domain
```

Внешние слои не зависят от внутренних, что позволяет беспрепятсвенно переходить от одной базы данных к другой, не ломая логику всей программы

### Что за что отвечает

| Файл                                | Слой       | Назначение                                                                 |
|-------------------------------------|------------|----------------------------------------------------------------------------|
| `cmd/server/main.go`                | точка входа| Читает конфиг, открывает БД, собирает слои вместе и запускает HTTP-сервер. |
| `internal/config/config.go`         | конфиг     | Читает настройки из окружения/`.env` (хост, порт, путь к БД и шаблону).    |
| `internal/domain/status.go`         | domain     | Тип `Status` (`on`/`off`) и его разбор (`ParseStatus`). Без зависимостей.  |
| `internal/repository/repository.go` | repository | Интерфейс `Status` (`Get`/`Set`) и сборка репозиториев.                    |
| `internal/repository/sqlite.go`     | repository | Открытие БД и создание таблицы (`NewDB`).                                   |
| `internal/repository/status.go`     | repository | Чтение/запись статуса в SQLite.                                            |
| `internal/service/service.go`       | service    | Интерфейс сервиса статуса и его сборка.                                     |
| `internal/service/status.go`        | service    | Логика: валидирует ввод и обращается к репозиторию.                         |
| `internal/service/errors.go`        | service    | Доменные ошибки сервиса (`ErrInvalidStatus`).                              |
| `internal/handler/http.go`          | handler    | HTTP-обработчик `/`: меняет статус по `?status=on/off` и рендерит шаблон.  |
| `web/index.html`                    | web        | HTML-шаблон страницы с тумблером.                                          |

### Как всё связано

В файле `main.go` вызываются и пробрасываются зависимости при старте программы.

```go
db, _ := repository.NewDB(cfg.DBPath)   // открыли SQLite
repo := repository.NewRepository(db)    // репозиторий поверх БД
svc  := service.NewService(repo)        // сервис поверх репозитория
h    := handler.New(svc, tmpl)          // handler поверх сервиса
```

Дальше при запросе `GET /?status=on`:
`handler` берёт параметр -> зовёт `service.Update` -> тот валидирует через `domain.ParseStatus`
и сохраняет через `repository.Set` → ответ рендерится из `web/index.html`.
Если параметра нет — текущий статус читается через `repository.Get`.
