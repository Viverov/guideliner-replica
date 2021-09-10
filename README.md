# Guideliner
___

Бекенд проекта для отображения дополнительной информации на втором экране - гайдов, статистики, изображений и так далее.


# Requirements:
___

[golang 1.16](https://golang.org/dl/)

[golangci-lint 1.40](https://golangci-lint.run/usage/install/)

# Commands
___

Команды для запуска описаны в makefile

+ `make lint`: запустить линтер
+ `make fmt`: запустить форматтер
+ `make generate`: запустить генерацию. Для системных нужд, вызывается другими командами.
+ `make tests-unit`: запустить юнит-тесты
+ `make tests-integration`: запустить интеграционные тесты. Вместе с этим, запускает все необходимые внешние зависимости (ex: DB, kafka) в докере (`docker-compose.test.dependencies`), сами тесты тоже запускаются в докере (`docker-compose.test.yml`)
+ `make run-guideliner-debug`: собрать и запустить бекенд в дебаг режиме.
+ `make run-guideliner-development`: собрать и запустить бекенд в режиме разработки
+ `make run-guideliner-production`: собрать и запустить проект в режиме продакшена
+ `make run-migrations-debug`: собрать и запустить миграции для дебаг режима
+ `make run-migrations-development`: собрать и запустить миграции для режима разработчика
+ `make run-migrations-production`: собрать и запустить миграции для продакшена
+ `make vendor`: собрать необходимые для сборки проекта пакеты в директории `{PROJECT_ROOT}/vendor`. Для системных нужд, вызывается другими командами.
+ `make build-guideliner`: собрать бекенд. Результат помещается в `{PROJECT_ROOT}/bin/guideliner`
+ `make build-migrations`: собрать команду, используемую для накатки миграций. Результат помещается в `{PROJECT_ROOT}/bin/migrations`.
+ `make build-clean-db`: собрать команду, пересоздающую схему public. Используется для тестов - чтобы получить чистый инстанс базы. Результат помещается в `{PROJECT_ROOT}/bin/clean_db`
+ `make build-all`: совмещает в себе `build-guideliner`, `build-migrations` и `build-clean-db`
+ `make create-migration`: создает файл-заглушку вида `{DATE_NOW}_tempfile.go` в директории `./internal/migrations/`. Нужна, чтобы не генерировать ID/файл миграции руками

#Config
___

`./bin/clean_db`, `./bin/guideliner`, `./bin/migrations` (и соотвествующие команды `make-{SOMETHING}-{MODE}`) могут конфигурироваться из переменных окружения (далее ENV) или из config.json.

Режим устанавливается только из ENV, с помощью `GUIDELINER_ENV`. Разрешенные значения: `DEBUG`, `TEST`, `DEVELOPMENT`, `PRODUCTION`. Стандартное значение: `PRODUCTION`

**Важно**: в случае использования `PRODUCTION` режима разрешена конфигурация только из переменных окружения. Это сделано во избежание деплоя в продакшен с дефолтной конфигурацией из `config.json`

Примеры конфигурационных файлов хранятся в директории `configs`.

Описание конфигурационных файлов:

### config.json
```json5
{
  "db": { // Настройки БД. Передаются в gorm.io/driver/postgres, https://gorm.io/docs/connecting_to_the_database.html#PostgreSQL
    "host": "127.0.0.1",
    "port": "5555", 
    "login": "dblogin", 
    "password": "dbpassword", 
    "name": "test_db", 
    "sslMode": "disable" 
  },
  "server": { // Настройки приложения
    "host": "127.0.0.1",
    "port": "8080"
  },
  "tokens": {
    "secretKey": "SECRET_KEY" // Секретный ключ для генерации JWT токенов
  }
}
```

### ENV

```dotenv
# Настройки приложения
GUIDELINER_SERVER_HOST=127.0.0.1
GUIDELINER_SERVER_PORT=8080

# Настройки базы
GUIDELINER_DB_HOST=db-test
GUIDELINER_DB_PORT=5432
GUIDELINER_DB_LOGIN=dblogin
GUIDELINER_DB_PASSWORD=dbpassword
GUIDELINER_DB_NAME=test_db
GUIDELINER_DB_SSL_MODE=disable

# Секретный ключ для генерации JWT токенов
GUIDELINER_TOKEN_SECRET=SECRET_KEY
```

# Примечания по GUIDELINER_ENV

На текущий момент, режимы "недореализованы". По большей части, это связано с тем, что логирование ещё не сделано.

Основные отличия на текущий момент:

+ `PRODUCTION` не позволяет инициализацию из `config.json`
+ `TEST` заглушает некоторые, особенно раздражающие логи. Например, инициализацию логов.
