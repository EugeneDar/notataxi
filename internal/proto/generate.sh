set -e

if [ ! -f LICENSE ]; then
    echo "Most likely you are running the script from wrong directory, run from the root of the repository"
    exit 1
fi

echo Generating all python protobufs in the internal/protobufs folder...
mkdir -p internal/protobufs
python3 -m grpc_tools.protoc \
    -I=internal/proto \
    --python_out=internal/protobufs \
    --grpc_python_out=internal/protobufs \
    internal/proto/*.proto

echo Generating protobufs for mocks in the appropriate folders...
python3 -m grpc_tools.protoc \
    -I=internal/proto \
    --python_out=cmd/mocks/config \
    --grpc_python_out=cmd/mocks/config \
    internal/proto/config.proto
python3 -m grpc_tools.protoc \
    -I=internal/proto \
    --python_out=cmd/mocks/executor_profile \
    --grpc_python_out=cmd/mocks/executor_profile \
    internal/proto/executor_profile.proto
python3 -m grpc_tools.protoc \
    -I=internal/proto \
    --python_out=cmd/mocks/order_data \
    --grpc_python_out=cmd/mocks/order_data \
    internal/proto/order_data.proto
python3 -m grpc_tools.protoc \
    -I=internal/proto \
    --python_out=cmd/mocks/toll_roads \
    --grpc_python_out=cmd/mocks/toll_roads \
    internal/proto/toll_roads.proto
python3 -m grpc_tools.protoc \
    -I=internal/proto \
    --python_out=cmd/mocks/zone \
    --grpc_python_out=cmd/mocks/zone \
    internal/proto/zone_data.proto

echo Generating all golang protobufs in the internal/protobufs folder...
GOPATH=$HOME/go
PATH=$PATH:$GOPATH/bin
mkdir -p internal/protobufs
protoc \
    -I=internal/proto \
    --go_out=internal/protobufs \
    --go-grpc_out=internal/protobufs \
    internal/proto/*.proto

echo Success!
