# Добро пожаловать в transaction-system !

Transaction-system - мое видение транзакционной системы с функциями выставления счетов и вывода средств, реализации статусов транзакций, актуального и замороженного баланса клиентов.

## 🚀 Доступные ручки

1. **Зачисление средств**
    
    - Метод: POST
    - Путь: `localhost:3000/invoice`
    - Описание: позволяет зачислять средства. Идентифицирует юзера по номеру кошелька `и/или` по номеру карты
```json
{
  "currency_code": 840,
  "amount": 100.50,
  "wallet_number" : 101234567,
  "card_number" : 5478396041568712
}
```
      
2. **Вывод средств**
    
    - Метод: POST
    - Путь: `localhost:3000/withdraw`
    - Описание: позволяет списывать средства. Идентифицирует юзера по номеру кошелька `и/или` по номеру карты
```json
{
  "currency_code": 643,
  "amount": 50.50,
  "wallet_number" : 789012345,
  "card_number" : 5267890123456789
}
```
      
3. **Получению актуального баланса**
    
    - Метод: GET
    - Путь: `localhost:3000/available-balance`
    - Обработчик: `NetServDB/controllers.AddUser`
    - Описание: выводит актуальный (Success) баланс юзера
```json
{
  "wallet_number" : 789012345,
  "card_number" : 5267890123456789
}
```
    
4. **Получению замороженного баланса**
   
    - Метод: GET
    - Путь: `localhost:3000/frozen-balance` (если при запуске не был указан в параметрах `-host` и `-port`)
    - Описание: выводит замороженный (Created) баланс юзера
```json
{
  "wallet_number" : 789012345,
  "card_number" : 5267890123456789
}
```
      

## 💡 Использованные технологии

Проект разработан с использованием следующих технологий:

- [**Gin**](https://github.com/gin-gonic/gin) - роутер
- [**zap**](https://github.com/uber-go/zap) - инструмент для эффективного логирования
- [**viper**](https://github.com/spf13/viper) - библиотека для конфигурирования приложения
- [**go-pg**](https://www.postgresql.org/) - миграции и ORM
- [**gocron**](https://github.com/go-co-op/gocron) - scheduler
- [**Postgresql**](https://www.postgresql.org/) 
- [**Kafka**](https://kafka.apache.org)
- [**NATS**](https://github.com/nats-io/nats.go)
- [**ZooKeeper**](https://zookeeper.apache.org)
