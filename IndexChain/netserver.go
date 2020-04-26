package main

import (
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

func (s *SuperNode) NetServer(b *Blockchain, c <-chan bool, wg *sync.WaitGroup) {

	defer wg.Done()
	// Listen from Other Node
	ln, err := net.Listen("tcp", ":"+s.Config.RpcPort)
	if err != nil {
		fmt.Println("[!!*!!]Listen :"+s.Config.RpcPort+" Failed!")
		log.Panic(err)
	}
	defer ln.Close()

	go func() {
		for {
			//fmt.Println("[~]NetServer-Inter Listen on :", s.Config.RpcPort)
			conn, err := ln.Accept()
			if err != nil {
				if strings.Contains(err.Error(), "use of closed network connection") {
					break
				}
				fmt.Println("[*]NetServer-Inter connection Wrong")
				continue
			}
			go s.handleInterConnection(conn, b)
		}
	}()

	//Listen from WeChat
	lnn, err := net.Listen("tcp", ":10050")
	if err != nil {
		fmt.Println("[!!*!!]Listen :10050 Failed!")
		log.Panic(err)
	}
	defer lnn.Close()

	go func() {
		for {
			//fmt.Println("[~]NetServer-Outer Listen on ", ":10050")
			conn, err := lnn.Accept()
			if err != nil {
				if strings.Contains(err.Error(), "use of closed network connection") {
					break
				}
				fmt.Println("[*]NetServer-Outer connection Wrong")
				continue
			}
			go s.handleOuterConnection(conn)
		}
	}()
	<-c
}
