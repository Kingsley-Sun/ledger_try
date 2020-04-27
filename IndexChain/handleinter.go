package main

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"strings"
)

func (s *SuperNode) handle_getall(m *Message, b *Blockchain, ip string) {
	//read the database
	//send all my block after specify height
	var askfor AskFor
	decoder := gob.NewDecoder(bytes.NewReader(m.Parameters))
	err := decoder.Decode(&askfor)
	//fmt.Println("[~]Get Getall request from ", ip)
	if err != nil {
		fmt.Println("[*]In handle_getall Decode Get error : ", err)
		return
	}

	height := askfor.Turn
	bestheight := b.GetBestHeight()
	if height > bestheight {
		//ignore the message maybe a bad request
		return
	}

	var blocks []*Block
	for {
		height = height + 1
		block, err := b.GetBlockByHeight(height)
		if err != nil {
			fmt.Println("[*]In handle_getall GetBlockByHeight Error", err)
			return
		}
		blocks = append(blocks, block)
		if height == bestheight {
			break
		}
	}
	s.Peers.ReturnMess_GetAll(blocks, askfor.Info.Province, askfor.Info.PublicKey, ip, askfor.Info.Rpcport)
}

func (s *SuperNode) handle_getblock(m *Message, b *Blockchain, ip string) {
	//read the database
	//send my one block with specify height
	var askfor AskFor
	decoder := gob.NewDecoder(bytes.NewReader(m.Parameters))
	err := decoder.Decode(&askfor)

	//fmt.Println("[~]Get getblock request from ", ip)
	if err != nil {
		fmt.Println("[*]In handle_getblock Decode Get error : ", err)
		return
	}

	height := askfor.Turn
	bestheight := b.GetBestHeight()
	if height > bestheight {
		//ignore the message maybe a bad request
		return
	}
	block, err := b.GetBlockByHeight(height)
	if err != nil {
		//fmt.Println("[*]In handle_getblock GetBlockByHeight Error", err)
		return
	}
	s.Peers.ReturnMess_GetBlock(block, askfor.Info.Province, askfor.Info.PublicKey, ip, askfor.Info.Rpcport)
}

func (s *SuperNode) handle_sendblock(m *Message, b *Blockchain) {
	//recieve a block check it and
	//add it to the database
	block := DeserializeBlock(m.Parameters)
	if block.VerifyMBlock(s.Peers.Superpeers) == false {
		fmt.Println("[*]Handle_sendblock VerifMBlock Failed")
		return
	}
	b.Mutex.Lock()
	if block.Header.Height <= b.Latesheight {
		return
	}
	//fmt.Println("[~]Add A Block to Mempool")
	b.AddToMempool(block)
	b.Mutex.Unlock()
}

func (s *SuperNode) handle_getonlinepeers(m *Message, ip string) {
	if s.ConsensusTurn > 0 {
		return
	}

	//fmt.Println("[~]In function handle_getonlinepeers")
	var askfor AskFor
	decoder := gob.NewDecoder(bytes.NewReader(m.Parameters))
	err := decoder.Decode(&askfor)
	//fmt.Println("[~]Get getonlinepeers request from ", ip)

	if err != nil {
		fmt.Println("[*]In handle_getonlinepeers Decode Get error : ", err)
		return
	}

	if s.Peers.IsLegal(askfor.Info.PublicKey, askfor.Info.Province) == false {
		fmt.Println("[*]In handle_getonlinepeers Peers.IsLegal failed")
		return
	}

	pubbytes, errs := hex.DecodeString(askfor.Info.PublicKey)
	if errs != nil {
		fmt.Println("[*]In handle_getonlinepeers DecodeString Get error : ", err)
		return
	}

	if VerifySig(pubbytes, askfor.Sig, []byte("What_is_new_onlinepeers")) == false {
		fmt.Println("[*]In handle_getonlinepeers VerifySig failed")
		return
	}

	onlinepeers := s.Peers.OnlinePeers

	s.Peers.ReturnMess_GetOnlinePeers(onlinepeers, askfor.Info.Province, askfor.Info.PublicKey, ip, askfor.Info.Rpcport)
}

func (s *SuperNode) handle_sendonlinepeers(m *Message, ip string) {
	if s.ConsensusTurn > 0 {
		return
	}
	//fmt.Println("[~]In function handle_sendonlinepeers")
	var peers []*Peer
	decoder := gob.NewDecoder(bytes.NewReader(m.Parameters))
	err := decoder.Decode(&peers)
	if err != nil {
		fmt.Println("[*]In handle_sendonlinepeers Decode Get error : ", err)
		return
	}
	//fmt.Println("[~]Get onlinerpeers From ", ip)
	s.Peers.Mutex.Lock()
	for _, onpeer := range peers {
		//fmt.Println("[~]This is ", onpeer.Province)
		//fmt.Println("[~]Before IsLegal")
		if s.Peers.IsLegal(hex.EncodeToString(onpeer.PublicKey), onpeer.Province) == false {
			//fmt.Println("[~]" + onpeer.Province + " s.Peers.IsLegal failed")
			continue
		}

		//fmt.Println("[~]Before Contains")
		if Contains(s.Peers.OnlinePeers, onpeer) == true {
			//fmt.Println("[~]" + onpeer.Province + " Contains failed")
			continue
		}

		//fmt.Println("[~]Before s.Peers.IsAlive")
		if s.Peers.IsAlive(onpeer) == false {
			//fmt.Println("[~]" + onpeer.Province + " IsAlive failed")
			continue
		}

		//fmt.Println("[~]Before append")
		s.Peers.OnlinePeers = append(s.Peers.OnlinePeers, onpeer)
		fmt.Println("[*]Add A New OnlinePeers:", onpeer.Province)
	}
	s.Peers.Mutex.Unlock()
}

func (s *SuperNode) handle_askforin(m *Message,ip string) {
	if s.ConsensusTurn > 0 {
		return
	}

	var askfor AskFor
	decoder := gob.NewDecoder(bytes.NewReader(m.Parameters))
	err := decoder.Decode(&askfor)
	if err != nil {
		fmt.Println("[*]In handle_askforin Decode Get error : ", err)
		return
	}

	pubbytes, err := hex.DecodeString(askfor.Info.PublicKey)
	if err != nil {
		fmt.Println("[*]In handle_askforin DecodeString Get error : ", err)
		return
	}
	askforpeer := &Peer{
		Province: askfor.Info.Province,
		PublicKey: pubbytes,
		Addr:      ip,
		Rpcport:   askfor.Info.Rpcport,
		Protocol:  "tcp",
	}

	//fmt.Println("[~]In handle_askforin")
	if s.Peers.IsLegal(askfor.Info.PublicKey, askfor.Info.Province) == false {
		fmt.Println("[*]" + askfor.Info.Province + " s.Peers.IsLegal failed")
		return
	}

	if VerifySig(pubbytes, askfor.Sig, []byte("AskForIn")) == false {
		fmt.Println("[*]" + askfor.Info.Province + " VerifySig failed")
		return
	}

	s.Peers.ReturnMess_AskForIn(s.Config,askfor.Info.Province, askfor.Info.PublicKey,askfor.Info.Rpcport,askfor.Turn,ip)

	s.Peers.Mutex.Lock()
	if Contains(s.Peers.OnlinePeers,askforpeer) == false {
		fmt.Println("[*]Add A New OnlinePeers:", askfor.Info.Province)
		s.Peers.OnlinePeers = append(s.Peers.OnlinePeers, askforpeer)
	}
	s.Peers.Mutex.Unlock()
}

func (s *SuperNode) handle_returnforin(m *Message) {
	if s.ConsensusTurn > 0 {
		return
	}

	var replyask AskFor
	decoder := gob.NewDecoder(bytes.NewReader(m.Parameters))
	err := decoder.Decode(&replyask)
	if err != nil {
		fmt.Println("[*]In handle_returnforin Decode Get error : ", err)
		return
	}

	pubbytes, err := hex.DecodeString(replyask.Info.PublicKey)
	if err != nil {
		fmt.Println("[*]In handle_returnforin DecodeString Get error : ", err)
		return
	}

	//fmt.Println("[~]Before s.Peers.IsLegal")
	if s.Peers.IsLegal(replyask.Info.PublicKey, replyask.Info.Province) == false {
		fmt.Println("[*]" + replyask.Info.Province + " s.Peers.IsLegal")
		return
	}

	//fmt.Println("[~]Before VerifySig")
	if VerifySig(pubbytes, replyask.Sig, []byte("YouAreIn")) == false {
		fmt.Println("[*]" + replyask.Info.Province + " VerifySig failed")
		return
	}

	//Make Sure It's this turn
	if replyask.Turn != s.LeaderTurn {
		fmt.Println("[*]replyask.Turn is not this Turn")
		return
	}

	fmt.Println("[*]Replyask Get A new LeaderCommit")
	s.LeaderCommits += 1
}

func (s *SuperNode) handleInterConnection(conn net.Conn, b *Blockchain) {
	var buf bytes.Buffer
	buffer := make([]byte, 1024)
	for {
		//althouh the length of message always shorter than 1024
		//use loop to meets the large message in the future
		n, _ := conn.Read(buffer)
		//use the condition n == 0 to mark the end of message
		//means that the client should shutdown the socket
		if n == 0 {
			break
		} else {
			buf.Write(buffer[:n])
		}
	}
	//decode struct
	var message Message
	decoder := gob.NewDecoder(bytes.NewReader(buf.Bytes()))
	err := decoder.Decode(&message)
	if err != nil {
		log.Panic(err)
	}

	//ip address
	ipaddr := strings.Split(conn.RemoteAddr().String(), ":")[0]
	switch message.MessageType {
	case M_getall:
		s.handle_getall(&message, b, ipaddr)
	case M_getblock:
		s.handle_getblock(&message, b, ipaddr)
	case M_sendblock:
		s.handle_sendblock(&message, b)

	case M_getonlinepeers:
		s.handle_getonlinepeers(&message, ipaddr)
	case M_sendonlinepeers:
		s.handle_sendonlinepeers(&message, ipaddr)

	case M_askforin:
		s.handle_askforin(&message,ipaddr)
	case M_returnforin:
		s.handle_returnforin(&message)

	case C_sendblock:
		s.handle_consensus(&message,b)
	default:
		fmt.Println("message type wrong")
	}
}
