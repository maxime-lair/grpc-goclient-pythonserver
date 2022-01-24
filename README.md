# grpc-goclient-pythonserver
grpc implementation with go client and python server

This project was only done to showcase GRPC and Bubbletea TUI implementation in golang, with a python server. It has no real usecase besides being an example. The code needs to be factorized, and improved, as It is overly declared.

A client sends a clientID and receives a list of possible socket family from the server, you can then choose one family and request the server to fetch you the possible socket types associated with it. Lastly you can pick a type (and a family), and request a list of possible protocol for this socket. You can then print your final results.

It is one of my first projects in golang and bubbletea usage, so feel free to comment on how to improve its structure or code. Thanks

## Requirement

The script requires protobuf compiler, python 3.XX and golang to be installed on your system. The installation process will not be defined here.

The starting script just checks for missing library, but do not install anything that could require sudo rights

## Usage

Start the python server with `./start_server.sh`

Start a new terminal, and start the client:
`./start_client.sh`

If you want to simply test the server (without golang/TUI), you can use a test client done in python:
`./start_test_client.sh`

A run example:

![linux_socket](https://user-images.githubusercontent.com/72258375/150876976-5f2cb4ad-d43c-43d1-ac45-0bdbc4c70de3.gif)


## Architecture

![image](https://user-images.githubusercontent.com/72258375/150876077-1194d84e-f5f9-4c6b-a40d-f81203bb2e56.png)

The following chapters define each directory and how they operate

## Protobuf

The protobuf defines the methods and message type used in both clients and server. It is shared by both.

It is defined in the **protos** directory

```
    rpc GetSocketFamilyList(SocketTree) returns (stream SocketFamily) {}

    rpc GetSocketTypeList(SocketFamily) returns (stream SocketType) {}

    rpc GetSocketProtocolList(SocketTypeAndFamily) returns (stream SocketProtocol) {}
```

## Python server

The python server is a simple GRPC server, that will answer to three types of requests:
- GetSocketFamily: receives an ID (just to trace the request), and return a list of socket family
- GetSocketType: receives a socket family and return a list of possible types
- GetSocketProtocol: receives a socket family and type and return a list of possible protocols

## Goclient

The goclient is built with a TUI in *bubbletea, lipgloss, bubbles*

It has a few bugs when It prints, but It is totally useable, I tried to add as many "features" as possible to test its limits, such as :
- A spinner
- A list to select and return the selected value
- A progress bar depending on where you are

The client is split into many files, with the main ones being:
- client_grpc : *functions used to send GRPC to the server*
- client_update: *functions to update the TUI depending on event received*
- client_view: *functions to print the TUI depending on which state we are in*

The TUI is split into four parts:
- The header, containing a title and the current status and clientID
- The selection list (or result)
- A help section to help you navigate
- A log journal, that only prints the 5 previous entry

## Wordlist

The client generates a clientID (not unique, but as randomized as possible) made from a list of two files:
- color.txt: contains 118 colors
- animal.txt: contains 115 animals

It just picks a random tuple from these file and use it to track the request made to the server (in case you use multiple clients), It is nice to retrace what happened in the logs
