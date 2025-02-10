# im

## Requirements
* Kafka
```shell
  cd $kafka_home
  bin/zookeeper-server-start.sh config/zookeeper.properties
  bin/kafka-server-start.sh config/server.properties
  bin/kafka-topics.sh --create --topic im-message-route --bootstrap-server localhost:9092

```
* Redis
```shell
cd $redis_home
redis-server
```
