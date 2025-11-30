# Wallet Service

Сервис управления кошельками с REST API для депозитов и выводов средств с учётом конкурентной нагрузки.

## Стек

- Go 1.25.1
- PostgreSQL 16
- pgx (драйвер)
- goose v3 (миграции)
- Docker, docker-compose

## Функциональность 
Модели выполнены в соответствии с требованиями ТЗ:
- POST /api/v1/wallets
  - Request:
    - walletId: UUID
    - operationType: "DEPOSIT" или "WITHDRAW"
    - amount: строка c двумя знаками после запятой, например "100.00"
  - Поведение:
    - Если кошелёк не существует, он создаётся с нулевым балансом.
    - Для DEPOSIT баланс увеличивается.
    - Для WITHDRAW при недостатке средств возвращается HTTP 402 Payment Required.
  - Response 200:
    - walletId: UUID
    - balance: строка в рублях с двумя знаками после запятой
    - operationStatus: "SUCCESS"
  - Response 402:
    - walletId
    - balance (текущий баланс)
    - operationStatus: "INSUFFICIENT_BALANCE"
    - error: "not_enough_money"

- GET /api/v1/wallets/{walletId}
  - Response 200:
    - walletId
    - balance
  - Response 404:
    - "wallet not found"

## Запуск

1. Создайте файл __config.env__ и скопируйте туда настройки __config.env.example__
2. Соберите и запустите систему:```docker-compose up -d --build``` (Сервис будет доступен на http://localhost:8080)
3. Для запуска тестов: ```go test ./...``` (Для запуска тестов необходимо наличие установленного языка Golang)
4. Для остановки работы сервиса необходимо ввести следующую команду: ```docker-compose down  ```

