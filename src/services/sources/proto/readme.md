Build from parent directory:

For Python:
```
python3 -m grpc_tools.protoc -I=proto --python_out=protobufs --grpc_python_out=protobufs proto/*.proto
python3 -m grpc_tools.protoc -I=proto --python_out=../mocks/config --grpc_python_out=../mocks/config proto/config.proto
python3 -m grpc_tools.protoc -I=proto --python_out=../mocks/executor_profile --grpc_python_out=../mocks/executor_profile proto/executor_profile.proto
python3 -m grpc_tools.protoc -I=proto --python_out=../mocks/order_data --grpc_python_out=../mocks/order_data proto/order_data.proto
python3 -m grpc_tools.protoc -I=proto --python_out=../mocks/toll_roads --grpc_python_out=../mocks/toll_roads proto/toll_roads.proto
python3 -m grpc_tools.protoc -I=proto --python_out=../mocks/zone --grpc_python_out=../mocks/zone proto/zone_data.proto
```

For Go:
```
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin

protoc --go_out=./protobufs --go-grpc_out=./protobufs -Iproto proto/*.proto
```
