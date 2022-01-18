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
python3 -m pip install grpcio
python3 -m pip install grpcio-tools
python3 -m grpc_tools.protoc -Iprotos/ --python_out=pythonserver/ --grpc_python_out=pythonserver/ protos/server.proto

printf "Starting python server\n"
# Start server
./pythonserver/server.py &
server_pid=$!
printf "Server started with PID %s\n" "$server_pid"

printf "Returning to previous working directory from %s to %s\n" "$cwd" "$pwd"
cd "$pwd"
