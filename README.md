# grpc-common

```shell
buf dep update

buf push
```

## test

```shell
protoc --proto_path=example  -I common/ --go_out=:example --go-grpc_out=:example example/database.proto

protoc-go-inject-tag -input=./example/database.pb.go 

protoc \
  --proto_path=./example \
  -I common/ \
  --plugin=protoc-gen-database-query=./bin/linux-amd64/protoc-gen-database-query \
  --database-query_out=:./example \
  example/database.proto
```
