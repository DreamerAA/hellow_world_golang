# Задача 12

Cоздание простого REST API для работы с некоторыми сущностями, например, задачами, пользователями или книгами: Добавить, изменять, удалять и получать данные о них.

## Описание задачи

### Создать REST API

. Разработайте базовый веб-сервер с помощью net/http.
. Реализуйте маршруты для следующих операций:
.. Получение всех элементов (GET /items).
.. Получение одного элемента по ID (GET /items/{id}).
.. Создание нового элемента (POST /items).
.. Обновление существующего элемента по ID (PUT /items/{id}).
.. Удаление элемента по ID (DELETE /items/{id}).

### Работа с данными

Для хранения данных о сущности с помощью базу данных PostgreSQL.

### Структура данных

. Определите структуру для данных, например:

```[go]
    type Item struct {
        ID      int    `json:"id"`
        Name    string `json:"name"`
        Details string `json:"details"`
    }
```

. Реализуйте автоинкремент ID для уникальности каждого элемента.

### JSON-формат

Поддержите ввод и вывод данных в формате JSON, используя encoding/json.

### Обработка ошибок

Добавьте обработку ошибок для каждого типа запроса (например, если элемент с заданным ID не найден, возвращайте статус 404).

## Основные шаги

. Настройка веб-сервера.
. Реализация CRUD-операций (Create, Read, Update, Delete).
. Тестирование API с помощью инструмента вроде Postman или curl.
