#!/usr/bin/env python3

import socket
from concurrent import futures
import logging

import grpc
import server_pb2
import server_pb2_grpc


## Require:
#  python3 -m pip install grpcio
#  python3 -m pip install grpcio-tools
#  python3 -m grpc_tools.protoc -Iprotos/ --python_out=pythonserver/ --grpc_python_out=pythonserver/ protos/server.proto

def socket_get_family_list(stub):
    logging.debug("SendSocketTree")
    socketFamilyList = stub.GetSocketFamilyList(server_pb2.SocketTree(choice="Alpha"))

    for socketFamily in socketFamilyList:
        logging.INFO("socket family received: %s : %s" % (socketFamily.name, socketFamily.value))


def socket_get_type_list(stub):
    logging.info("GetSocketTypeList")
    

def socket_get_protocol_list(stub):
    logging.info("GetSocketProtocolList")


def run():
    logging.info("Starting to run client")
    with grpc.insecure_channel('localhost:50051') as channel:
        stub = server_pb2_grpc.SocketGuideStub(channel)
        logging.info("-------------- SendSocketTree --------------")
        socket_get_family_list(stub)
        logging.info("-------------- GetSocketTypeList --------------")
        socket_get_type_list(stub)
        logging.info("-------------- GetSocketProtocolList --------------")
        socket_get_protocol_list(stub)

if __name__ == "__main__":
    logging.basicConfig()
    run()