package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"strconv"
	"strings"
	"sync"
)

const NodeRequestServiceName = "NodeRequestService"

type NodeService struct {
	Node  *SuperNode
	Chain *Blockchain
}

type NodeRequestServiceInterface interface {
	//For IndexChain
	HeartBeats(args string,reply *string) error
	RequestCommits(args string,reply *string) error

	//For DeepSee
	GetBlockNotes(args string, reply *string) error
	GetSuperNodes(args string, reply *string) error
	GetLatesetBlockCount(args string, reply *string) error
}

func RegisterNodeRequestService(src NodeRequestServiceInterface) error {
	return rpc.RegisterName(NodeRequestServiceName, src)
}

func (n *NodeService) HeartBeats(args string,reply *string) error {
	if args == "Are_You_Ok?" {
		*reply = "I_Am_Ok"
	}else {
		*reply = "Bad"
	}
	return nil
}

func (n *NodeService) RequestCommits(args string,reply *string) error {
	//fmt.Println("[~]Recieve Block agrs")
	if n.Node.ConsensusTurn != 4 {
		*reply = "Reject"
		return nil
	}

	bss,_ := hex.DecodeString(args)
	//fmt.Println(bss)
	recvblock := *DeserializeBlock(bss)

	if IsSameBlock(recvblock,*(n.Node.ConsensusBlock)) == true {
		fmt.Println("[*]RequestCommits Recieve A Legal Block")
		*reply = "Commit"
	}else {
		fmt.Println("[*]RequestCommits Recieve A Illegal Block")
		*reply = "Reject"
	}
	return nil
}


func (n *NodeService) GetBlockNotes(args string, reply *string) error {
	blocknum, err := strconv.Atoi(args)
	if err != nil {
		*reply = "Parameters wrong! Can't transfer to Int"
		return err
	}
	result, err := n.Chain.GetNotesByBlockHeight(blocknum)
	if err != nil {
		*reply = "Parameters wrong! Maybe Too Big Or <0"
		return err
	}
	jsonbyte, _ := json.Marshal(result)
	*reply = string(jsonbyte)
	return nil
}

func (n *NodeService) GetSuperNodes(args string, reply *string) error {
	resultpeers := n.Node.Peers.OnlinePeers
	jsonbyte, _ := json.Marshal(resultpeers)
	*reply = string(jsonbyte)
	return nil
}

func (n *NodeService) GetLatesetBlockCount(args string, reply *string) error {
	result := n.Chain.GetBestHeight()
	*reply = strconv.Itoa(result)
	return nil
}

func (s *SuperNode) RpcServer(blockchain *Blockchain, c <-chan bool, wg *sync.WaitGroup) {
	defer wg.Done()

	err := RegisterNodeRequestService(&NodeService{
		Node:  s,
		Chain: blockchain,
	})
	if err != nil {
		log.Panic(err)
	}

	listener, err := net.Listen("tcp", "0.0.0.0:10000")
	//fmt.Println("[~]RPCServer Listen on ", ":10000")

	if err != nil {
		fmt.Println("[!!*!!]RPCSever Listen on 10000 FAILED")
		log.Panic(err)
	}
	defer listener.Close()

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				if strings.Contains(err.Error(), "use of closed network connection") {
					break
				}
				//fmt.Println("[~]Rpc connection Wrong")
				continue
			}
			//fmt.Println("[~] A new rpc connection accept from ",conn.RemoteAddr().String())
			go rpc.ServeCodec(jsonrpc.NewServerCodec(conn))
		}
	}()
	<-c
}

/*
#-*- coding:utf-8 -*-
#--Auther--:Kingsley

from socket import *

HOST = '127.0.0.1'
PORT = 10000
BUFSIZ = 1024
ADDRESS = (HOST, PORT)

tcpClientSocket = socket(AF_INET, SOCK_STREAM)
tcpClientSocket.connect(ADDRESS)

data = "{\"method\":\"NodeRequestService.GetBlockNotes\",\"params\":[\""
data += "1"
data += "\"],\"id\":1000}"

# 发送数据
tcpClientSocket.send(data)
# 接收数据
data, ADDR = tcpClientSocket.recvfrom(BUFSIZ)
print(data.encode("utf-8"))
*/
