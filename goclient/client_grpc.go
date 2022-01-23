package main

import (
	"context"
	"fmt"
	"io"
	pb "main/pb_server"
	"time"
)

/*************
GRPC part
*************/

func socket_get_family_list(clientEnv clientEnv) ([]socketChoice, clientEnv) {

	clientEnv.logJournal = append(clientEnv.logJournal, fmt.Sprintf("[%s][GetFamilyList] Entering.", clientEnv.clientID.Name))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	socketFamilyStream, req_err := (*clientEnv.client).GetSocketFamilyList(ctx, clientEnv.clientID)
	if req_err != nil {
		clientEnv.err = errMsg{req_err}
		clientEnv.logJournal = append(clientEnv.logJournal, clientEnv.err.Error())
		return nil, clientEnv
	}
	var socketFamilyList []socketChoice
	for {
		family, stream_err := socketFamilyStream.Recv()
		if stream_err == io.EOF {
			break
		}
		if stream_err != nil {
			clientEnv.err = errMsg{req_err}
			clientEnv.logJournal = append(clientEnv.logJournal, clientEnv.err.Error())
			return nil, clientEnv
		}
		//clientEnv.logJournal = append(clientEnv.logJournal, fmt.Sprintf("[%s][GetFamilyList] Received family: %s", clientEnv.clientID.Name, family))
		socketFamilyList = append(socketFamilyList, socketChoice{
			Name:  family.Name,
			Value: family.Value,
		})
	}
	clientEnv.logJournal = append(clientEnv.logJournal, fmt.Sprintf("[%s][GetFamilyList] len=%d cap=%d", clientEnv.clientID.Name, len(socketFamilyList), cap(socketFamilyList)))

	return socketFamilyList, clientEnv
}

func socket_get_type_list(clientEnv clientEnv, clientChoice clientChoice) ([]socketChoice, clientEnv) {

	clientEnv.logJournal = append(clientEnv.logJournal, fmt.Sprintf("[%s][GetTypeList] Entering with family: %d --> %s", clientEnv.clientID.Name, clientChoice.selectedFamily.Value, clientChoice.selectedFamily.Name))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	socketTypeStream, req_err := (*clientEnv.client).GetSocketTypeList(ctx, &pb.SocketFamily{
		Name:     clientChoice.selectedFamily.Name,
		Value:    clientChoice.selectedFamily.Value,
		ClientId: clientEnv.clientID,
	})
	if req_err != nil {
		clientEnv.err = errMsg{req_err}
		clientEnv.logJournal = append(clientEnv.logJournal, clientEnv.err.Error())
		return nil, clientEnv
	}

	var socketTypeList []socketChoice
	for {
		socketType, stream_err := socketTypeStream.Recv()
		if stream_err == io.EOF {
			break
		}
		if stream_err != nil {
			clientEnv.err = errMsg{req_err}
			clientEnv.logJournal = append(clientEnv.logJournal, clientEnv.err.Error())
			return nil, clientEnv
		}
		clientEnv.logJournal = append(clientEnv.logJournal, fmt.Sprintf("[%s][GetTypeList] Received family: %s", clientEnv.clientID.Name, socketType))
		socketTypeList = append(socketTypeList, socketChoice{
			Name:  socketType.Name,
			Value: socketType.Value,
		})
	}
	clientEnv.logJournal = append(clientEnv.logJournal, fmt.Sprintf("[%s][GetTypeList] len=%d cap=%d", clientEnv.clientID.Name, len(socketTypeList), cap(socketTypeList)))

	return socketTypeList, clientEnv
}

func socket_get_protocol_list(clientEnv clientEnv, clientChoice clientChoice) ([]socketChoice, clientEnv) {

	clientEnv.logJournal = append(clientEnv.logJournal, fmt.Sprintf("[%s][GetProtocolList] Entering with family: [%d] %s -- [%d] %s",
		clientEnv.clientID.Name,
		clientChoice.selectedFamily.Value, clientChoice.selectedFamily.Name,
		clientChoice.selectedType.Value, clientChoice.selectedType.Name))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	socketProtocolStream, req_err := (*clientEnv.client).GetSocketProtocolList(ctx, &pb.SocketTypeAndFamily{
		Family: &pb.SocketFamily{
			Name:  clientChoice.selectedFamily.Name,
			Value: clientChoice.selectedFamily.Value,
		},
		Type: &pb.SocketType{
			Name:  clientChoice.selectedType.Name,
			Value: clientChoice.selectedType.Value,
		},
		ClientId: clientEnv.clientID,
	})
	if req_err != nil {
		clientEnv.err = errMsg{req_err}
		clientEnv.logJournal = append(clientEnv.logJournal, clientEnv.err.Error())
		return nil, clientEnv
	}

	var socketProtocolList []socketChoice
	for {
		socketProtocol, stream_err := socketProtocolStream.Recv()
		if stream_err == io.EOF {
			break
		}
		if stream_err != nil {
			clientEnv.err = errMsg{req_err}
			clientEnv.logJournal = append(clientEnv.logJournal, clientEnv.err.Error())
			return nil, clientEnv
		}
		clientEnv.logJournal = append(clientEnv.logJournal, fmt.Sprintf("[%s][GetProtocolList] Received protocol: %s", clientEnv.clientID.Name, socketProtocol))
		socketProtocolList = append(socketProtocolList, socketChoice{
			Name:  socketProtocol.Name,
			Value: socketProtocol.Value,
		})
	}
	clientEnv.logJournal = append(clientEnv.logJournal, fmt.Sprintf("[%s][GetProtocolList] len=%d cap=%d", clientEnv.clientID.Name, len(socketProtocolList), cap(socketProtocolList)))

	return socketProtocolList, clientEnv

}
