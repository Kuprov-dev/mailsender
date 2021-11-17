# Запуск Kafka

```bash
docker-compose up -d
```

**Примечание**:

если вы перезапускали кластер не дожидаясь graceful termination, то вы можете столкнуться с ошибкой `Error while creating ephemeral at /brokers/ids/1001, node already exists...`. Для того, чтобы исправить ее выполните:

```bash
docker-compose rm --stop --force -v
```

# Создание Topic

Создадим топик `quickstart-events`:

```bash
docker exec \
    kafka_kafka_1 kafka-topics.sh \
    --create --topic quickstart-events --bootstrap-server :9092
```

Запишем сообщения в этот топик. Для этого подключимся к контейнеру:

```bash
docker exec -it kafka_kafka_1 /bin/bash
```

И выполним:

```bash
kafka-console-producer.sh --topic quickstart-events --bootstrap-server localhost:9092
> First message
> Second message
```

Ввод сообщений заканчивается нажатием `Ctrl-D`

Прочитаем сообщения:

```bash
kafka-console-consumer.sh --topic quickstart-events --from-beginning --bootstrap-server localhost:9092
```

Остановить консьюмер можно с помощью `Ctrl-C`

# Пишем приложение на Go

Создадим топик для тестирования:

```bash
docker exec simplemailsender_kafka_1 kafka-topics.sh --create --topic test-topic --bootstrap-server :9092
```

Запустим тест:

```bash
KAFKA_BROKER=127.0.0.1:9093 go test ./... -tags=integration_tests
```

Прочтем сообщения из топика. Подключимся к контейнеру:

```bash
docker exec -it kafka_kafka_1 /bin/bash
```

и выполним:

```bash
kafka-console-consumer.sh --topic test-topic --from-beginning --bootstrap-server localhost:9092
```