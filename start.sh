#!/bin/sh

# Check for nounset and errors throwing
set -eu

pwd="$(pwd)"
cwd="$(dirname "$0")"

if [ "$pwd" == "$cwd" ]; then
    printf "Current working directory is the bash script, perfect.\n"
else
    printf "Current working directory is different than the one where the bash script is located, changing..\n"
    cd "$cwd"
fi

printf "Getting dependencies..\n"

if pip3 check grpcio >/dev/null; then
    printf "grpcio already installed\n"
else
    python3 -m pip install grpcio
fi

if pip3 check grpcio-tools >/dev/null; then
    printf "grpcio-tools already installed\n"
else
    python3 -m pip install grpcio-tools
fi

printf "Compiling protos files\n"
python3 -m grpc_tools.protoc -Iprotos/ --python_out=pythonserver/ --grpc_python_out=pythonserver/ protos/server.proto

printf "Starting python server\n"
# Start server
./pythonserver/server.py
server_pid=$!
printf "Server started with PID %s\n" "$server_pid"

printf "Returning to previous working directory from %s to %s\n" "$cwd" "$pwd"
cd "$pwd"
