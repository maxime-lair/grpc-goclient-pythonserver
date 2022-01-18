#!/usr/bin/env python3                                                                 

import socket                                                                          

## Require:
#  python -m pip install grpcio
#  python -m pip install grpcio-tools

table = {num:name[8:] for name,num in vars(socket).items() if name.startswith("IPPROTO")}

def socketProtocol(socketFamily, socketType):                                          
    # IP Protocol number go up to 143 (up to 255 but other are for testing)            
    for i in range(144):                                                               
        try:                                                                           
            sock = socket.socket(socketFamily._value_, socketType._value_, i)          
            print('\t\t\t\t', table[i])                                                
        except (KeyError, OSError) as msg:                                             
            sock = None                                                                
            continue                                                                   


def socketSupportedType(socketFamily):                                                 
    # for each type, print the available protocol                                      
    print(socketFamily._name_)                                                         
    for socketType in socket.SocketKind:                                               
        try:                                                                           
            sock = socket.socket(socketFamily._value_, socketType._value_)             
            print('\t\t', socketType._name_)                                           
            socketProtocol(socketFamily, socketType)                                   
        except OSError as msg:                                                         
            sock = None                                                                
            continue                                                                   

def socketTree():                                                                      
    # for each socket family, print the available type                                 
        # for each type, print the available protocol                                  
    print("Family\t\tType\t\tProtocols")                                               
    for socketFamily in socket.AddressFamily:                                          
        socketSupportedType(socketFamily)                                              


def main():                                                                            
    socketTree()                                                                       

if __name__ == "__main__":                                                             
    main()         