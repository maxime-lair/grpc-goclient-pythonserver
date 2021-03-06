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

table = {num:name[8:] for name,num in vars(socket).items() if name.startswith("IPPROTO")}

class SocketGuideServicer(server_pb2_grpc.SocketGuideServicer):
    """ Provides methods that implement functionality of route guide server """

    def __init__(self):
        logging.info("Creating SocketGuideServicer")
    
    def GetSocketFamilyList(self, request_iterator, context):
        logging.info("[%s] Entering GetSocketFamilyList" % (request_iterator.name))
        
        logging.info("[%s] Client wishes to receive socket family list for : %s" % (request_iterator.name, request_iterator.name))

        for socketFamily in socket.AddressFamily:
            yield server_pb2.SocketFamily(name=socketFamily._name_,
                                                value=int(socketFamily._value_))

    def GetSocketTypeList(self, request_iterator, context):
        logging.info("[%s] Entering GetSocketTypeList" % (request_iterator.client_id.name))
        logging.info("[%s] Testing possible type for %s" % (request_iterator.client_id.name, request_iterator.name))
        for socketType in socket.SocketKind:
            logging.debug('[%s] Testing out [%s] %s' % (request_iterator.client_id.name, socketType._value_,socketType._name_))                               
            try:                                                                      
                sock = socket.socket(request_iterator.value, socketType._value_)             
                logging.debug('[%s] Sending possible socketType %s for %s' % (request_iterator.client_id.name, socketType._name_, request_iterator.name))
                yield server_pb2.SocketType(name=socketType._name_, value=socketType._value_)
            except OSError as msg:                                                         
                sock = None
                continue

    def GetSocketProtocolList(self, request_iterator, context):
        logging.info("[%s] Entering GetSocketProtocolList" % (request_iterator.client_id.name))
        logging.info("[%s] Testing possible type for %s - %s" % (request_iterator.client_id.name, request_iterator.family.name, request_iterator.type.name))
        for i in range(144):                                                               
            try:                                                                           
                sock = socket.socket(request_iterator.family.value, request_iterator.type.value, i)          
                yield server_pb2.SocketProtocol(name=table[i], value=int(i))                                           
            except (KeyError, OSError) as msg:                                             
                sock = None                              
                continue            

    
def serve():
    logging.info("Starting to serve")
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    server_pb2_grpc.add_SocketGuideServicer_to_server(
        SocketGuideServicer(), server)
    server.add_insecure_port('[::]:50051')
    server.start()
    server.wait_for_termination()

if __name__ == "__main__":
    logging.basicConfig(level=logging.DEBUG)
    serve()