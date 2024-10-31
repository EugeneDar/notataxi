#!/bin/bash

if [ "$(pwd)" != "$(cd "$(dirname "$0")" && pwd)" ]; then
    echo "Please navigate to the src/services/sources/proto directory and run the script again."
    exit 1
fi

python3 -m grpc_tools.protoc \
    --proto_path=. \
    --python_out=../protobufs \
    --grpc_python_out=../protobufs \
    *.proto
