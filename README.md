# im

## 代码生成
```shell
$ cd src/api

// 单独生成packet
$ protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative packet.proto

// packet 和 kitex一起生成
$ kitex broker.ptoto
$ kitex business.proto
$ kitex router.proto

```

## Starting on IDE
### build
```shell
go mod tidy
```

```shell
cd client
go mod tidy
```

### Start MysQL Redis Etcd and Kafka
```shell
docker compose up -d
```
### Start im-broker
```shell
go run im-broker.go
```

## Start client
```shell
cd client
go run im-client.go
```