Build from parent directory:

For Python:
```
python3 -m grpc_tools.protoc -I=proto --python_out=protobufs --grpc_python_out=protobufs proto/*.proto
```

For Go:
```
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin

protoc --go_out=./protobufs --go-grpc_out=./protobufs -Iproto proto/*.proto
```
