# resp — довідник

Повний довідник пакета `resp`: ментальна модель, помічники відповідей, білдер
`Response`, опції, гарантія безпечного за замовчуванням кодування, безпекові
помічники, помилки й практичні рецепти.

Англійська версія: **[DOC.md](DOC.md)**.

## Зміст

- [Ментальна модель](#ментальна-модель)
- [Помічники відповідей](#помічники-відповідей)
- [Білдер Response](#білдер-response)
- [Опції](#опції)
- [Відповіді-помилки](#відповіді-помилки)
- [Безпечне за замовчуванням кодування](#безпечне-за-замовчуванням-кодування)
- [JSONP](#jsonp)
- [Кастомний JSON-енкодер](#кастомний-json-енкодер)
- [Проєкція полів](#проєкція-полів)
- [Файли й завантаження](#файли-й-завантаження)
- [Безпекові заголовки](#безпекові-заголовки)
- [Конкурентність](#конкурентність)
- [Помилки](#помилки)
- [Рецепти й поради](#рецепти-й-поради)

## Ментальна модель

`resp` створює HTTP-відповіді поверх стандартного `net/http`. Замість того щоб
вручну виставляти тип вмісту, писати статус і кодувати тіло, ви викликаєте один
помічник, що робить усі три речі узгоджено — і повідомляє про помилку кодування
до того, як щось зафіксовано.

Дві форми покривають усі випадки:

1. **Пакетні помічники** — `resp.JSON(w, data, opts...)` і подібні: один виклик,
   `http.ResponseWriter`, опційні функціональні опції.
2. **Білдер `Response`** — `resp.NewResponse(w).SetStatus(...).SetHeader(...)`:
   ланцюгова конфігурація, тоді термінальний метод (`JSON`, `HTML`, …).

Обидва приймають ті самі значення `Option`, тож усе, що можна зробити інлайн,
можна й на білдері.

```go
import "github.com/goloop/resp/v2"
```

Зручний тип `resp.R` — це `map[string]any`, зручно для ad-hoc JSON:

```go
resp.JSON(w, resp.R{"message": "Hello", "status": "success"})
```

## Помічники відповідей

Кожен помічник бере `w`, своє навантаження й опційні `...Option`:

| Помічник | Надсилає |
|----------|----------|
| `JSON(w, data)`            | тіло `application/json` |
| `JSONIndent(w, data, prefix, indent)` | форматований JSON |
| `JSONP(w, data, callback)` | валідований JSONP |
| `XML(w, data)`             | тіло `application/xml` |
| `HTML(w, html)`            | тіло `text/html` |
| `String(w, text)`          | тіло `text/plain` |
| `Blob(w, contentType, bytes)` | сирі байти, `Content-Length` виставлено |
| `Stream(w, reader)`        | потокове тіло з `io.Reader` |
| `Error(w, status, message)` | JSON-тіло помилки (див. нижче) |
| `Created(w, location, data)` | `201` + `Location` + JSON-тіло |
| `Accepted(w, data)`        | `202` + JSON-тіло |
| `NoContent(w)`             | `204` |
| `Redirect(w, url)`         | редирект (лише валідні 3xx) |
| `ServeFile(w, r, filename)`| віддати файл із диска |
| `ServeFileAsDownload(w, filename, data)` | завантаження з пам'яті |

Дефолтні типи вмісту застосовуються автоматично (`application/json;
charset=utf-8` тощо), а дефолтний статус — `200 OK`.

## Білдер Response

```go
func NewResponse(w http.ResponseWriter, opts ...Option) *Response
```

Білдер дзеркалить кожен помічник як метод і додає ланцюгову конфігурацію:

```go
resp.NewResponse(w).
    SetStatus(http.StatusOK).
    SetHeader("X-Custom-Header", "value").
    JSON(resp.R{"status": "success"})
```

Методи конфігурації (`SetStatus`, `SetHeader`/`AddHeader`/`DelHeader`,
`SetCookie`/`DelCookie`/`ExpiredCookie`, `SetJSONEncoder`, …) повертають
`*Response` для ланцюга; термінальні методи (`JSON`, `HTML`, `XML`, `String`,
`Blob`, `Stream`, `Error`, `Created`, `Accepted`, `NoContent`, `Redirect`,
`ServeFile`) пишуть відповідь і повертають `error`. Ґетери `Status()` і
`Header()` читають поточний стан.

## Опції

Опції — це значення `func(*Response) *Response`, які приймає кожен помічник і
`NewResponse`. Вони поділяються на кілька родин:

- **Статус** — `WithStatus(code)` плюс іменована константа для кожного статусу
  (`WithStatusCreated()`, `WithStatusNotFound()`, …).
- **Заголовки** — `WithHeader(key, values...)`, плюс типізовані помічники `Add…`
  для повного набору стандартних заголовків (`AddETag`, `AddLastModified`,
  `AddCacheControl`, `AddContentDisposition(type, filename)`, CORS
  `AddAccessControl…`, `AddLink`, `AddWarning` та багато інших).
- **Тип вмісту** — шорткати `As…` (`AsApplicationJSON()`, `AsTextHTML()`,
  `AsOctetStream()`, …).
- **Кукі** — `WithCookie(cookie)`.
- **Поведінка JSON** — `WithStreamJSON()`, `ApplyJSONEncoder(fn)`.
- **Помилки** — `WithErrorCode(n)`.
- **Безпека** — `SecurityHeaders()`, `WithCSP(policy)`.

```go
resp.JSON(w, data,
    resp.WithStatus(http.StatusOK),
    resp.WithHeader("X-Custom-Header", "value"),
    resp.AddETag(`"123456"`),
)
```

## Відповіді-помилки

```go
func Error(w http.ResponseWriter, status int, message string, opts ...Option) error
```

**Перший аргумент — завжди HTTP-статус.** Він пишеться в рядок статусу й
дзеркалиться в поле `code` JSON-тіла, тож вони ніколи не розходяться:

```go
// HTTP 404, тіло {"code":404,"message":"Not Found"}.
resp.Error(w, http.StatusNotFound, "")

// HTTP 404, тіло {"code":1042,"message":"User not found"}.
resp.Error(w, http.StatusNotFound, "User not found", resp.WithErrorCode(1042))
```

`WithErrorCode(n)` виставляє окремий бізнес-код у тілі, зберігаючи HTTP-статус.
Аргумент статусу має перевагу над опцією `WithStatus`, а значення поза діапазоном
відкочується до `500 Internal Server Error`.

## Безпечне за замовчуванням кодування

За замовчуванням `JSON`/`JSONP`/`XML` кодують спершу в пульований буфер, тож
помилка серіалізації повертається **до** того, як зафіксовано рядок статусу, —
хендлер усе ще може надіслати належну відповідь-помилку:

```go
if err := resp.JSON(w, data); err != nil {
    // Ще нічого не записано, тож запасний варіант спрацює.
    resp.Error(w, http.StatusInternalServerError, "Internal server error")
    return
}
```

Для дуже великих навантажень, де пам'ять важливіша за цю гарантію, увімкніть
потокове кодування без копіювання:

```go
resp.JSON(w, largeData, resp.WithStreamJSON())
```

Пульований буфер тримає дефолтний (безпечний) шлях на тій самій кількості
алокацій, що й пряме потокове кодування.

## JSONP

Ім'я колбека валідується проти суворого білого списку JavaScript-ідентифікаторів
**до** того, як щось записано; небезпечне ім'я повертає
[`ErrInvalidJSONPCallback`](#помилки) і нічого не пише, закриваючи вектор
відбитого XSS:

```go
callback := r.URL.Query().Get("callback")
if err := resp.JSONP(w, data, callback); err != nil {
    // errors.Is(err, resp.ErrInvalidJSONPCallback) для відхиленого імені.
    resp.Error(w, http.StatusBadRequest, "invalid callback")
}
```

## Кастомний JSON-енкодер

Під'єднайте будь-який енкодер, що задовольняє `JSONEncodeFunc`
(`func(w io.Writer, v any) error`):

```go
enc := func(w io.Writer, v any) error {
    e := json.NewEncoder(w)
    e.SetEscapeHTML(false)
    return e.Encode(v)
}

resp.NewResponse(w, resp.ApplyJSONEncoder(enc)).JSON(resp.R{"message": "ok"})
```

`JSONIndent(w, data, prefix, indent)` дає форматований JSON для орієнтованих на
людину чи дебаг-ендпоінтів.

## Проєкція полів

Формуйте навантаження до серіалізації:

```go
data := resp.OnlyFields(user, "ID", "Email") // лишити лише ці
data = resp.ExcludeFields(user, "Password")   // або прибрати ці
resp.JSON(w, data)
```

## Файли й завантаження

```go
// Віддати статичний файл, нав'язавши ім'я завантаження.
resp.ServeFile(w, r, "path/to/file.pdf",
    resp.AddContentDisposition("attachment", "document.pdf"))

// Віддати вміст із пам'яті як завантаження.
resp.ServeFileAsDownload(w, "file.txt", []byte("Dynamic content"))
```

`ServeFileAsDownload` виставляє `Content-Length`, `X-Content-Type-Options:
nosniff` і форму RFC 5987 `filename*` для не-ASCII імен.

## Безпекові заголовки

`SecurityHeaders()` застосовує консервативний, широко безпечний пресет
(`nosniff`, `X-Frame-Options: DENY`, `Referrer-Policy:
strict-origin-when-cross-origin`). CSP — явна опція, бо корисна політика залежить
від застосунку:

```go
resp.JSON(w, data,
    resp.SecurityHeaders(),
    resp.WithCSP("default-src 'self'"),
)
```

## Конкурентність

Як і `http.ResponseWriter`, відповідь `resp` **не** безпечна для конкурентного
використання. Пишіть з однієї горутини й лише поки хендлер не повернувся. Щоб
поєднати роботу кількох горутин, збирайте результати через канали й пишіть у
горутині хендлера:

```go
select {
case err := <-errChan:
    resp.Error(w, http.StatusInternalServerError, err.Error())
case result := <-resultChan:
    resp.JSON(w, result)
case <-time.After(5 * time.Second):
    resp.Error(w, http.StatusGatewayTimeout, "Processing timeout")
}
```

## Помилки

```go
var ErrInvalidJSONPCallback = errors.New("invalid JSONP callback name")
```

Повертається з `JSONP` для імені колбека, що не пройшло білий список
ідентифікаторів. Звіряйте через `errors.Is`. Усі помічники повертають помилку
кодування/запису (чи `nil`), тож ви вирішуєте, як реагувати.

## Рецепти й поради

**Завжди перевіряйте помилку — відновитися безпечно.** У дефолтному буферованому
режимі невдалий `JSON` нічого не записав, тож запасна відповідь `Error` спрацює.

**Тримайте статус і тіло синхронно.** Використовуйте `Error` (статус
дзеркалиться в `code`) замість ручного статусу плюс окремого коду тіла; додавайте
`WithErrorCode`, лише коли бізнес-код відрізняється від HTTP-статусу.

**Беріть білдер, коли конфігурації багато.** Ланцюжте `SetStatus`/`SetHeader`/
`SetCookie` на `NewResponse`, а не передавайте довгий список `...Option`.

**Потоково віддавайте велике.** Використовуйте `Stream` для великих файлів і
`WithStreamJSON()` для дуже великого JSON, обмінюючи буферовану страховку на
меншу пам'ять.

**Ховайте чутливі поля.** Запускайте `ExcludeFields`/`OnlyFields` перед `JSON`,
щоб пароль чи токен ніколи не дійшов до енкодера.
