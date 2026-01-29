Вот вариант **простого и аккуратного README**, без перегруза, но с понятной идеей и примерами 👌
Можно класть прямо в `README.md`.

---

# errorx

`errorx` — небольшой Go-пакет для **структурированных ошибок** с поддержкой:

* namespace’ов
* кодов ошибок
* произвольных полей
* error wrapping (`errors.Is / As / Unwrap`)
* сериализации **цепочки ошибок в JSON**

Подходит для API, микросервисов и логирования.

---

## Возможности

* 📦 Ошибки с `namespace` и `code`
* 🔗 Цепочки ошибок через `Wrap`
* 🧾 JSON-представление всей цепочки ошибок
* 🔍 Удобные проверки (`IsCode`, `IsNamespace`)
* 🧠 Совместимость со стандартным `errors`

---

## Установка

```bash
go get github.com/AlexsRyzhkov/errorx
```

---

## Быстрый старт

### Создание namespace

```go
authErr := errorx.NewNamespace("auth")
```

### Новая ошибка

```go
err := authErr.New(
    "USER_NOT_FOUND",
    "user not found",
    "user_id", 42,
)
```

---

## Оборачивание ошибок

```go
err := db.QueryRow(...)
if err != nil {
    return authErr.Wrap(
        err,
        "DB_ERROR",
        "failed to load user",
        "user_id", 42,
    )
}
```

Или без предварительного namespace:

```go
return errorx.Wrap(
    err,
    "auth",
    "DB_ERROR",
    "failed to load user",
)
```

---

## Проверки ошибок

### Проверка типа errorx

```go
if errorx.IsErrx(err) {
    // это errorx
}
```

### Проверка кода

```go
if errorx.IsCode(err, "USER_NOT_FOUND") {
    // пользователь не найден
}
```

### Проверка namespace

```go
if errorx.IsNamespace(err, "auth") {
    // ошибка из auth
}
```

### Код + namespace

```go
if errorx.IsCodeNS(err, "auth", "DB_ERROR") {
    // auth.DB_ERROR
}
```

---

## JSON-сериализация

`errorx` сериализует **всю цепочку ошибок**:

```go
data, _ := json.Marshal(err)
fmt.Println(string(data))
```

Пример вывода:

```json
[
  {
    "namespace": "auth",
    "code": "DB_ERROR",
    "message": "failed to load user",
    "fields": {
      "user_id": 42
    },
    "caller": "service.go:37"
  },
  {
    "type": "*pq.Error",
    "message": "connection refused"
  }
]
```

---

## Generic `As`

```go
if e, ok := errorx.As[*errorxError](err); ok {
    fmt.Println(e.Code)
}
```

Работает поверх стандартного `errors.As`.

---

## Лицензия

MIT

---
