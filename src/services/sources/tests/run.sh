#!/bin/bash


# Generate the protobuf files
cd ../proto
bash gen.sh
cd ../tests


# Check if the current directory is src/services/sources/tests
if [ "$(pwd)" != "$(cd "$(dirname "$0")" && pwd)" ]; then
    echo "Please navigate to the src/services/sources/tests directory and run the script again."
    exit 1
fi


# Run the tests
pytest functional.py
# pytest load.py


# Remove the generated files
# rm ../protobufs/*pb2*.py
