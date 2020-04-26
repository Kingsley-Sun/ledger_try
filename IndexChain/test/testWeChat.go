package main

import (
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "192.168.1.11:10050")
	if err == nil {
		defer conn.Close()
	} else {
		fmt.Println("net.Dial has error ", err)
	}

	conn, err = net.Dial("tcp", "192.168.1.22:10050")
	if err == nil {
		defer conn.Close()
	} else {
		fmt.Println("net.Dial has error ", err)
	}


	conn, err = net.Dial("tcp", "192.168.1.33:10050")
	if err == nil {
		defer conn.Close()
	} else {
		fmt.Println("net.Dial has error ", err)
	}

	conn, err = net.Dial("tcp", "192.168.1.44:10050")
	if err == nil {
		defer conn.Close()
	} else {
		fmt.Println("net.Dial has error ", err)
	}

}
