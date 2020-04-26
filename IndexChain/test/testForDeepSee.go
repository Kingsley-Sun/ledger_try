package main

import (
	"fmt"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

func main() {
	conn, err := net.Dial("tcp", "192.168.1.11"+":10000")
	if err != nil {
		fmt.Println("net.Dial error",err)
	}

	client := rpc.NewClientWithCodec(jsonrpc.NewClientCodec(conn))

	var reply string

	args := "3"
	err = client.Call("NodeRequestService.GetBlockNotes", args, &reply)
	if err != nil {
		fmt.Println("client.Call error",err)
	}
	fmt.Println(reply)

	args = ""
	err = client.Call("NodeRequestService.GetSuperNodes", args, &reply)
	if err != nil {
		fmt.Println("client.Call error",err)
	}
	fmt.Println(reply)

	args = ""
	err = client.Call("NodeRequestService.GetLatesetBlockCount", args, &reply)
	if err != nil {
		fmt.Println("client.Call error",err)
	}
	fmt.Println(reply)

}