# kafkaManager

Этот проект представляет собой JSON REST API сервер, написанный на языке Go. Сервер предоставляет две основные ручки:

### *Возможные функции:*
- *POST /orders/*
добавление сообщения с заказом для сохранения и дальнейшей обработки
- *GET /status/* - получение статуса обработки заказа
- 
### *База данных:*

*Для хранения данных используется СУБД PostgreSQL. В проекте определена таблица:*
- message_tracker - содержит информацию о поступивших сообщениях, включая уникальный идентификатор, статус и сериализованное сообщение.

## Документация Swagger

Проект включает документацию Swagger для облегчения работы с API. После запуска сервера вы можете получить доступ к документации: /swagger/*

## Окружение

### Требования
- Docker 20.x.x
- Docker Compose 1.29.x
- Go 1.19+

### Переменные окружения

| Параметр                           | Описание                                               | Значение по умолчанию                          |
|------------------------------------|--------------------------------------------------------|------------------------------------------------|
| `ENV_TYPE`                         | Режим окружения                                        | `development`                                  |
| `HTTP_SERVER_PORT`                 | Порт HTTP сервера                                      | `8081`                                         |
| `HTTP_SERVER_NAME`                 | Имя HTTP сервера                                       | `kafka-manager`                                |
| `HTTP_SERVER_TIMEOUT`              | Таймаут HTTP сервера                                   | `4s`                                           |
| `HTTP_SERVER_IDLE_TIMEOUT`         | Таймаут бездействия HTTP сервера                       | `60s`                                          |
| `DB_HOST`                          | Хост базы данных                                       | `localhost`                                    |
| `DB_PORT`                          | Порт базы данных                                       | `5432`                                         |
| `DB_NAME`                          | Имя базы данных                                        | `pkk_db`                                       |
| `DB_USER_NAME`                     | Имя пользователя базы данных                           | `user`                                         |
| `DB_PASSWORD`                      | Пароль пользователя базы данных                        | `secret`                                       |
| `DB_SSL_MODE`                      | Режим SSL для базы данных                              | `disable`                                      |
| `DB_DRIVER_NAME`                   | Имя драйвера базы данных                               | `postgres`                                     |
| `DB_MAX_CONNS`                     | Максимальное количество соединений БД                  | `20`                                           |
| `MIGRATION_URL`                    | URL для миграций                                       | `file://migrations`                            |
| `KAFKA_BOOTSTRAP_SERVERS`          | Список серверов Kafka                                  | `kafka:9092`                                   |
| `ZOOKEEPER_HOST`                   | Хост Zookeeper                                         | `zookeeper`                                    |
| `ZOOKEEPER_PORT`                   | Порт Zookeeper                                         | `2181`                                         |
| `KAFKA_PORT`                       | Порт Kafka                                             | `9092`                                         |
| `KAFKA_EXPOSE_PORT`                | Порт, на котором Kafka будет доступен извне            | `9093`                                         |
| `KAFKA_BROKER`                     | Адрес брокера Kafka                                    | `kafka:9092`                                   |
| `KAFKA_TOPIC`                      | Тема Kafka                                             | `order_status`                                 |
| `KAFKA_ADVERTISED_LISTENERS`       | Список рекламируемых слушателей Kafka                  | `INSIDE://kafka:9093,OUTSIDE://kafka:9092`     |
| `KAFKA_LISTENER_NAMES`             | Имена слушателей Kafka                                 | `INSIDE,OUTSIDE`                               |
| `KAFKA_LISTENERS`                  | Список слушателей Kafka                                | `INSIDE://0.0.0.0:9093,OUTSIDE://0.0.0.0:9092` |
| `KAFKA_ZOOKEEPER_CONNECT`          | Адрес подключения к Zookeeper                          | `zookeeper:2181`                               |
| `KAFKA_INTER_BROKER_LISTENER_NAME` | Имя слушателя для межброкерского общения               | `INSIDE`                                       |
| `KAFKA_RETRIES`                    | Количество попыток повторного подключения Kafka        | `5`                                            |
| `KAFKA_TIMEOUT`                    | Таймаут для Kafka                                      | `1s`                                           |


## Установка

**PROJECT**

- Создать новую директорию для проекта. В консоли перейти в созданную директорию и написать: git clone https://github.com/ShevelevEvgeniy/kafkaManager
- добавить сертификата и приватного ключа и прописать пути до них в переменные окружения например: HTTP_SERVER_CERT_FILE='/app/config/secrets/server.crt'
  HTTP_SERVER_KEY_FILE='/app/config/secrets/server.key'

**DOCKER**

*Сборка:*

Скопировать файл .env.dist и переименовать в .env, настроить параметры окружения cp .env.dist .env
Для развертывания, запустите установку, выполнив команду ниже: make install

*Служебное:*
- make migrate-up - Запуск миграций
- make migrate-down - Откат миграций
- make migrate-create name="$" - Создание новой миграции
- make run-tests - запуск тестов

Если у вас возникли вопросы или проблемы, вы можете связаться со мной по адресу Z_shevelev@mail.ru