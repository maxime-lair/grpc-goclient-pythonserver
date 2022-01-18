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
        logging.info("[GetSocketFamilyList] %s -> %s" % (socketFamily.name, socketFamily.value))
        return socketFamily
    

def socket_get_type_list(socketFamilyChoice, stub):
    logging.info("GetSocketTypeList")
    
    if(socketFamilyChoice is not None):
        socketTypeList = stub.GetSocketTypeList(server_pb2.SocketFamily(
            name=socketFamilyChoice.name,
            value=socketFamilyChoice.value
        ))

        for socketType in socketTypeList:
            logging.info("[GetSocketTypeList] %s -> %s" % (socketType.name, socketType.value))
            return socketType
    else:
        logging.error("[GetSocketTypeList] Received NoneType in socketFamilyChoice")



def socket_get_protocol_list(socketTypeChoice, stub):
    logging.info("GetSocketProtocolList")

    if(socketTypeChoice is not None):
        socketProtocolList = stub.GetSocketProtocolList(server_pb2.SocketType(
        name=socketTypeChoice.name,
        value=socketTypeChoice.value
    ))
        for socketProtocol in socketProtocolList:
            logging.info("[GetSocketProtocolList] %s -> %s" % (socketProtocol.name, socketProtocol.value))
    else:
        logging.error("[GetSocketProtocolList] Received NoneType in socketTypeChoice")


def run():
    logging.info("Starting to run client")
    with grpc.insecure_channel('localhost:50051') as channel:
        stub = server_pb2_grpc.SocketGuideStub(channel)
        logging.info("-------------- SendSocketTree --------------")
        socketFamilyChoice = socket_get_family_list(stub)
        logging.info("-------------- GetSocketTypeList --------------")
        socketTypeChoice = socket_get_type_list(socketFamilyChoice, stub)
        logging.info("-------------- GetSocketProtocolList --------------")
        socket_get_protocol_list(socketTypeChoice, stub)

if __name__ == "__main__":
    logging.basicConfig(level=logging.DEBUG)
    run()