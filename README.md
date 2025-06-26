# Параллельный веб-краулер с REST API

Веб-краулер на Go, с простым RESTful API для управления задачами обхода и запроса результатов.

Краулер скачивает страницы параллельно, парсит заголовки и внутренние ссылки, сохраняет данные в PostgreSQL и предоставляет эндпоинты для запуска/остановки краула, списка страниц, очистки базы и статистики.

---

## Возможности

- **Параллельная загрузка** с настраиваемым числом воркеров  
- **BFS обход** до указанной пользователем глубины  
- **Парсинг HTML** (`<title>`, внутренние/внешние ссылки) через goquery  
- **Хранение в PostgreSQL** (JSONB, полнотекстовый индекс по заголовкам)  
- **REST API** на базе Gin  
- **Конфигурация через .env** с помощью godotenv  

---

## Установка

```bash
git clone https://github.com/yourusername/Parallel-Web-Crawler-with-REST-API.git
cd Parallel-Web-Crawler-with-REST-API
go mod download
```

---

## Конфигурация

Создайте файл `.env` в корне проекта (добавьте в `.gitignore`):

```dotenv
# DSN для PostgreSQL
DATABASE_URL=postgres://crawler:secret@localhost:5432/crawlerdb?sslmode=disable

# Таймаут HTTP-запросов
FETCH_TIMEOUT=10s

# Число параллельных воркеров
WORKERS=5

# Размер буфера очереди URL
QUEUE_SIZE=100
```

---

## Запуск

```bash
go run cmd/crawler/main.go
```

---

## REST API

Все эндпоинты доступны по префиксу `/api`.

### POST /api/crawl/start

Запустить новый краул.

**Запрос:**
```json
{
  "seeds":    ["http://golang.com"],
  "max_depth": 1
}
```

**Ответ:**
```json
{"status":"started"}
```

### POST /api/crawl/stop

Остановить текущий краул и вернуть финальную статистику.

**Ответ:**
```json
{
  "fetched": 38,
  "errors":   0,
  "in_queue": 0
}
```


### GET /api/pages?q=<keyword>

Получить список всех сохранённых URL, можно фильтровать по ключевому слову.

**Ответ:**
```json
{
  "pages": [
    "http://golang.com",
  ]
}
```

### DELETE /api/pages

Очистить все страницы из базы данных.

**Ответ:**
```json
{"status":"cleared"}
```

### GET /api/stats

Получить текущую статистику краулинга без остановки.

**Ответ:**
```json
{
  "fetched": 12,
  "errors":   1,
  "in_queue": 3
}
```

---

## Примеры использования

```bash
# Очистить старые данные
curl -X DELETE http://localhost:8080/api/pages

# Старт краула (глубина=1)
curl -X POST http://localhost:8080/api/crawl/start   -H "Content-Type: application/json"   -d '{"seeds":["http://golang.com"],"max_depth":1}'

# Остановить и получить статистику
curl -X POST http://localhost:8080/api/crawl/stop

# Список сохранённых страниц
curl http://localhost:8080/api/pages

# Проверить статистику снова
curl http://localhost:8080/api/stats
```

---

## Тестирование

```bash
go test ./pkg/fetcher
go test ./pkg/parser
go test ./pkg/crawler
go test ./pkg/storage
```

