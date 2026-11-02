# Ledger и Gateway

## Эндпоинты

| Метод | Путь | Ответ при успехе |
| --- | --- | --- |
| POST | `/api/transactions` | 201, созданная транзакция |
| GET | `/api/transactions` | 200, массив транзакций |
| POST | `/api/budgets` | 201, созданный или обновлённый бюджет |
| GET | `/api/budgets` | 200, массив бюджетов |

Пустые списки возвращаются как `[]`. Для всех ответов установлен
`Content-Type: application/json; charset=utf-8`.
Ошибки имеют вид `{"error":"описание"}`:

- 400 — некорректный JSON или ошибка доменной валидации;
- 409 — превышение бюджета, `{"error":"budget exceeded"}`;
- 500 — непредвиденная внутренняя ошибка;
- 404 — неизвестный путь;
- 405 — неподдерживаемый метод, с заголовком `Allow: GET, POST`.

Сумма и лимит должны быть конечными положительными числами,
категория не должна быть пустой или состоять из пробелов.
`description` и строковая `date` сохраняются и возвращаются клиенту;
проверка формата даты опущена, как в предыдущем задании.
Тело ограничено 1 МиБ, допускается один JSON-объект без неизвестных полей.
Бюджет действует на всю историю категории. Достижение лимита разрешено,
категории без бюджета не ограничены. Повторный POST бюджета обновляет лимит.

Создать бюджет (201 Created):

```cmd
curl.exe -i -X POST http://localhost:8080/api/budgets -H "Content-Type: application/json" -d "{\"category\":\"еда\",\"limit\":500}"
```

Добавить транзакцию (201 Created):

```cmd
curl.exe -i -X POST http://localhost:8080/api/transactions -H "Content-Type: application/json" -d "{\"amount\":450,\"category\":\"еда\",\"description\":\"ланч\",\"date\":\"2025-09-10\"}"
```

Превысить бюджет (409 Conflict); транзакция не сохранится:

```cmd
curl.exe -i -X POST http://localhost:8080/api/transactions -H "Content-Type: application/json" -d "{\"amount\":100,\"category\":\"еда\",\"description\":\"кофе\",\"date\":\"2025-09-10\"}"
```

Ожидаемое тело ошибки:

```json
{"error":"budget exceeded"}
```

Получить списки (200 OK). В списке транзакций должна остаться только
успешная транзакция на 450 рублей:

```cmd
curl.exe -i http://localhost:8080/api/transactions
curl.exe -i http://localhost:8080/api/budgets
```

Некорректный бюджет (400 Bad Request); прежний лимит 500 сохранится:

```cmd
curl.exe -i -X POST http://localhost:8080/api/budgets -H "Content-Type: application/json" -d "{\"category\":\"еда\",\"limit\":0}"
```

Обновить бюджет (201 Created):

```cmd
curl.exe -i -X POST http://localhost:8080/api/budgets -H "Content-Type: application/json" -d "{\"category\":\"еда\",\"limit\":1000}"
```

Та же команда с переносами строк для CMD:

```cmd
curl.exe -i -X POST http://localhost:8080/api/budgets ^
  -H "Content-Type: application/json" ^
  -d "{\"category\":\"еда\",\"limit\":1000}"
```

Если появляется `invalid JSON: EOF`, проверьте, что запрос содержит `-d`
с JSON-телом. Ошибки `Bad hostname`, `Could not resolve host: application`
и `invalid character '\''` в показанных ранее командах связаны с кавычками
или переносами строк: используйте приведённые выше команды для CMD.

Middleware логирует метод, путь и длительность каждого запроса в консоль.
HTTP DTO и преобразователи находятся в `gateway/internal/api/dto.go`;
доменные модели и бизнес-проверка бюджета остаются в Ledger.

## Тестирование и покрытие

Все unit-тесты Ledger и интеграционные тесты Gateway успешно проходят.

| Область | Покрытие операторов |
| --- | ---: |
| Доменный пакет `ledger` | 98,2% |
| HTTP-пакет `gateway/internal/api` | 94,6% |
| Модуль Ledger целиком, включая `cmd/demo` | 70,9% |
| Модуль Gateway целиком, включая `main` | 88,8% |

Итоговое покрытие модулей ниже покрытия основных пакетов, поскольку
консольные точки входа `main` не покрыты тестами.

Проверяются:
- валидация транзакций и бюджетов;
- добавление и обновление бюджетов;
- достижение лимита и отказ при его превышении;
- отсутствие сохранения отклонённых транзакций;
- загрузка бюджетов из JSON и обработка ошибок;
- HTTP-статусы, JSON-ответы и заголовок Content-Type;
- логирование и обработка параллельных запросов.

Каждый сценарий использует отдельный экземпляр Ledger.
Очистка выполняется через `t.Cleanup(book.Reset)`.
Независимые доменные подтесты запускаются параллельно.

### Запуск в Windows CMD или PowerShell

Выполните из корня проекта:

```text
cd ledger
go test ./... -v -cover -coverprofile cover.out
go tool cover -html cover.out -o cover.html
go tool cover -func cover.out

cd ../gateway
go test ./... -v -cover -coverprofile cover.out
go tool cover -html cover.out -o cover.html
go tool cover -func cover.out
```

В каждой папке будут созданы:
- `cover.out` — профиль покрытия;
- `cover.html` — HTML-отчёт для просмотра в браузере.

