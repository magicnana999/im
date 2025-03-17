# im

```shell
$ protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative heartbeat.proto

$ protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative message.proto  


```



## Requirements
* Kafka
```shell
  cd $kafka_home
  bin/zookeeper-server-start.sh config/zookeeper.properties
  bin/kafka-server-start.sh config/server.properties
  bin/kafka-topics.sh --create --topic msg-route --bootstrap-server localhost:9092

```
* Redis
```shell
cd $redis_home
redis-server
```


# Kafka

## Topic

### 创建
``` shell
bin/kafka-topics.sh --create --topic msg-route --bootstrap-server localhost:9092 --partitions 1 --replication-factor 1
bin/kafka-topics.sh --create --topic msg-store --bootstrap-server localhost:9092 --partitions 1 --replication-factor 1
bin/kafka-topics.sh --create --topic msg-push --bootstrap-server localhost:9092 --partitions 1 --replication-factor 1
bin/kafka-topics.sh --create --topic msg-offline --bootstrap-server localhost:9092 --partitions 1 --replication-factor 1

bin/kafka-topics.sh --create --topic 127.0.0.1-7539 --bootstrap-server localhost:9092 --partitions 1 --replication-factor 1

```
* partitions 分区8个
* replication-factor 副本1个

### 删除
``` shell
bin/kafka-topics.sh --delete --topic msg-route --bootstrap-server localhost:9092
bin/kafka-topics.sh --delete --topic msg-store --bootstrap-server localhost:9092
bin/kafka-topics.sh --delete --topic msg-push --bootstrap-server localhost:9092
bin/kafka-topics.sh --delete --topic msg-offline --bootstrap-server localhost:9092

```

### 查询
``` shell
bin/kafka-topics.sh --list --bootstrap-server localhost:9092
```

```shell
bin/kafka-topics.sh --describe --topic msg-route --bootstrap-server localhost:9092
```

```text
Topic: msg-route	TopicId: wvB01QusSg2964GQdrYowQ	PartitionCount: 8	ReplicationFactor: 1	Configs:
	Topic: msg-route	Partition: 0	Leader: 0	Replicas: 0	Isr: 0	Elr: N/A	LastKnownElr: N/A
	Topic: msg-route	Partition: 1	Leader: 0	Replicas: 0	Isr: 0	Elr: N/A	LastKnownElr: N/A
	Topic: msg-route	Partition: 2	Leader: 0	Replicas: 0	Isr: 0	Elr: N/A	LastKnownElr: N/A
	Topic: msg-route	Partition: 3	Leader: 0	Replicas: 0	Isr: 0	Elr: N/A	LastKnownElr: N/A
	Topic: msg-route	Partition: 4	Leader: 0	Replicas: 0	Isr: 0	Elr: N/A	LastKnownElr: N/A
	Topic: msg-route	Partition: 5	Leader: 0	Replicas: 0	Isr: 0	Elr: N/A	LastKnownElr: N/A
	Topic: msg-route	Partition: 6	Leader: 0	Replicas: 0	Isr: 0	Elr: N/A	LastKnownElr: N/A
	Topic: msg-route	Partition: 7	Leader: 0	Replicas: 0	Isr: 0	Elr: N/A	LastKnownElr: N/A
```

#### 输出解析  ```Topic: msg-route	TopicId: wvB01QusSg2964GQdrYowQ	PartitionCount: 8	ReplicationFactor: 1	Configs:```


1. Topic: msg-route 
   * 这是正在查询的 Kafka topic 名称。
2. TopicId: wvB01QusSg2964GQdrYowQ 
   * 每个 topic 在 Kafka 中有一个唯一的 ID，这里是 msg-route 的 TopicId。
3. PartitionCount: 8 
   * 该 topic 有 8 个分区。这意味着该 topic 中的消息会被分散到这 8 个分区中。
4. ReplicationFactor: 1 
   * 每个分区的副本数量为 1。这意味着该 topic 的每个分区只有一个副本，集群容错能力较低。如果一个 broker 宕机，可能会丢失数据。
5. Configs: 
   * 该行没有列出具体的配置项，表示该 topic 使用了 Kafka 的默认配置。

#### 每个分区输出解析  ```Topic: msg-route	Partition: 0	Leader: 0	Replicas: 0	Isr: 0	Elr: N/A	LastKnownElr: N/A```
1. Partition: 0 
   * 这是 msg-route topic 的第 0 号分区。Kafka 会将消息分配到这些分区。
2. Leader: 0 
   * 该分区的 Leader 位于 Broker 0。Kafka 的分区有一个 leader，负责处理所有的读写请求。leader 负责数据的写入和分发，其他副本只是进行数据同步。
3. Replicas: 0 
   * 该分区的副本（replica）存储在 Broker 0 上。由于 ReplicationFactor 为 1，所以只有一个副本。这里的 0 表示副本位置为 Broker 0，但没有显示其他副本（因为副本数为 1）。
4. Isr: 0 
   * ISR（In-Sync Replicas）是与 leader 保持同步的副本列表。0 表示没有副本与 leader 保持同步。这通常是不正常的状态，可能会导致数据丢失的风险。
5. Elr: N/A 和 LastKnownElr: N/A 
   * 这些字段表示 "Expired Leader Replica"（过期的领导副本），即那些已不再是 leader 的副本。此字段的 N/A 表示没有过期副本。


### 消费状态
```shell
bin/kafka-consumer-groups.sh --describe --group 127.0.0.1-7539-group --bootstrap-server localhost:9092
bin/kafka-consumer-groups.sh --describe --group msg-route-group --bootstrap-server localhost:9092
```
```textmate
Consumer group 'msg-route-group' has no active members.

GROUP                  TOPIC            PARTITION  CURRENT-OFFSET  LOG-END-OFFSET  LAG             CONSUMER-ID     HOST           CLIENT-ID
msg-route-group msg-route 3          43543           43543           0               -               -               -
msg-route-group msg-route 2          43336           43336           0               -               -               -
msg-route-group msg-route 1          43145           43145           0               -               -               -
msg-route-group msg-route 0          43749           43749           0               -               -               -
msg-route-group msg-route 7          43438           43438           0               -               -               -
msg-route-group msg-route 6          43529           43529           0               -               -               -
msg-route-group msg-route 5          43743           43743           0               -               -               -
msg-route-group msg-route 4          43635           43635           0               -               -               -%
```

##### 说明
* CURRENT-OFFSET：当前消费到的消息的偏移量。
* LOG-END-OFFSET：日志中最后一条消息的偏移量。
* LAG：滞后量，表示消费者组尚未消费的消息数（LOG-END-OFFSET - CURRENT-OFFSET）。