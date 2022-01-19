#!/usr/bin/env python3

import socket
from concurrent import futures
import logging
import random

import grpc
import server_pb2
import server_pb2_grpc


## Require:
#  python3 -m pip install grpcio
#  python3 -m pip install grpcio-tools
#  python3 -m grpc_tools.protoc -Iprotos/ --python_out=pythonserver/ --grpc_python_out=pythonserver/ protos/server.proto

def define_client_id():
    color_picked = random.choice(list(open('../wordlist/color.txt')))
    animal_picked = random.choice(list(open('../wordlist/animal.txt')))
    return color_picked+"_"+animal_picked

def socket_get_family_list(client_id, stub):
    logging.debug("[%s] SendSocketTree" % (client_id.name))
    socketFamilyList = stub.GetSocketFamilyList(client_id)
    
    for socketFamily in socketFamilyList:
        logging.info("[%s][GetSocketFamilyList] %s -> %s" % (client_id.name, socketFamily.name, socketFamily.value))
        return socketFamily
    

def socket_get_type_list(client_id, socketFamilyChoice, stub):
    logging.info("[%s] GetSocketTypeList" % (client_id.name))
    
    if(socketFamilyChoice is not None):
        socketTypeList = stub.GetSocketTypeList(server_pb2.SocketFamily(
            name=socketFamilyChoice.name,
            value=socketFamilyChoice.value,
            client_id=client_id
        ))

        for socketType in socketTypeList:
            logging.info("[%s][GetSocketTypeList] %s -> %s" % (client_id.name, socketType.name, socketType.value))
            return socketType
    else:
        logging.error("[%s][GetSocketTypeList] Received NoneType in socketFamilyChoice" % (client_id.name))



def socket_get_protocol_list(client_id, socketTypeChoice, stub):
    logging.info("[%s] GetSocketProtocolList" % (client_id.name))

    if(socketTypeChoice is not None):
        socketProtocolList = stub.GetSocketProtocolList(server_pb2.SocketType(
        name=socketTypeChoice.name,
        value=socketTypeChoice.value,
        client_id=client_id
    ))
        for socketProtocol in socketProtocolList:
            logging.info("[%s][GetSocketProtocolList] %s -> %s" % (client_id.name, socketProtocol.name, socketProtocol.value))
    else:
        logging.error("[%s][GetSocketProtocolList] Received NoneType in socketTypeChoice" % (client_id.name))


def run():
    client_id = server_pb2.SocketTree(client_id=define_client_id())
    logging.info("Starting to run client id:%s" % (client_id.name))
    with grpc.insecure_channel('localhost:50051') as channel:
        stub = server_pb2_grpc.SocketGuideStub(channel)
        logging.info("-------------- SendSocketTree --------------")
        socketFamilyChoice = socket_get_family_list(client_id, stub)
        logging.info("-------------- GetSocketTypeList --------------")
        socketTypeChoice = socket_get_type_list(client_id, socketFamilyChoice, stub)
        logging.info("-------------- GetSocketProtocolList --------------")
        socket_get_protocol_list(client_id, socketTypeChoice, stub)

if __name__ == "__main__":
    logging.basicConfig(level=logging.DEBUG)
    run()