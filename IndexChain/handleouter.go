package main

import (
	"fmt"
	"net"
	"time"
)

func (s *SuperNode) handleOuterConnection(conn net.Conn) {
	//var buf bytes.Buffer
	//buffer := make([]byte, 1024)
	//for {
	//	//althouh the length of message always shorter than 1024
	//	//use loop to meets the large message in the future
	//	n, _ := conn.Read(buffer)
	//	//use the condition n == 0 to mark the end of message
	//	//means that the client should shutdown the socket
	//	if n == 0 {
	//		break
	//	} else {
	//		buf.Write(buffer[:n])
	//	}
	//}
	//
	//hash := buf.String()
	hash := "This is hash "+time.Now().String()

	stamp := time.Now().Format("20060102150405")
	note := &Note{
		HashID:    hash,
		Timestamp: stamp,
	}

	s.Mutex.Lock()
	s.Config.Mempool.AddNote(note)
	s.Mutex.Unlock()

	fmt.Println("Add a Note Hash: ", hash, " TimeStamp: ", stamp)
}
