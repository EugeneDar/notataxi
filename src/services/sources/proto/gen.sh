python3 -m grpc_tools.protoc \
    --proto_path=. \
    --python_out=../protobufs \
    --grpc_python_out=../protobufs \
    *.proto
