#!/bin/sh

# Check for nounset and errors throwing
set -eu

pwd="$(pwd)"
cwd="$(dirname "$0")/goclient"

if [ "$pwd" == "$cwd" ]; then
    printf "Current working directory is the bash script, perfect.\n"
else
    printf "Current working directory is different than the one where the bash script is located, changing..\n"
    cd "$cwd"
fi

go mod init main
go mod tidy

printf "Getting dependencies..\n"

if go list ./... | grep -i "protoc-gen-go@" >/dev/null; then
    printf "protoc-gen-go already installed\n"
else
    go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
fi

if go list ./... | grep -i "protoc-gen-go-grpc@" >/dev/null; then
    printf "protoc-gen-go-grpc already installed\n"
else
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1
fi

# Export go path to include GOPATH
if [[ *"$(go env GOPATH)/bin"* == "$PATH" ]]; then
    printf "PATH already have GOPATH defined\n"
else
    printf "PATH does NOT have GOPATH defined, adding $(go env GOPATH) to $PATH \n"
    export PATH="$PATH:$(go env GOPATH)/bin"
fi

printf "Compiling protos files\n"
# Require protobuf compiler installed
# PROTOC_ZIP=protoc-3.14.0-linux-x86_64.zip
# curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v3.14.0/$PROTOC_ZIP
# sudo unzip -o $PROTOC_ZIP -d /usr/local bin/protoc
# sudo unzip -o $PROTOC_ZIP -d /usr/local 'include/*'
# rm -f $PROTOC_ZIP
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ../protos/server.proto

printf "Starting go client\n"
printf "###########\n"
# Start client
go run goclient/client.go
printf "###########\n"
printf "Returning to previous working directory from %s to %s\n" "$cwd" "$pwd"
cd "$pwd"
