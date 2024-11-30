set -e

if [ ! -f LICENSE ]; then
    echo "Most likely you are running the script from wrong directory, run from the root of the repository"
    exit 1
fi

echo Deleting all automatically generated protobuf files...
rm -rf internal/protobufs
find . -type f -name '**.pb.go' -delete
find . -type f -name '**_pb2.py' -delete
find . -type f -name '**_pb2_grpc.py' -delete

echo Success!
