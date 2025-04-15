# im

```shell
$ cd src/api

// 单独生成packet
$ protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative packet.proto

// packet 和 kitex一起生成
$ kitex broker.ptoto
$ kitex business.proto
$ kitex router.proto

```

