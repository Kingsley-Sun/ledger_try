package main

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"fmt"
)

type Message struct {
	//Define the message type from node to node
	//Use gob to encode it and tcp to transfer it
	MessageType int
	Parameters  []byte
}

const (
	//recieve this kind message
	//means send all blocks after specify height
	M_getall int = iota
	//recieve this kind message
	//means sent block with specify height
	M_getblock
	//recieve this kind message
	//means some node send a block to me
	//Using M_getall or M_getblock
	M_sendblock

	//recieve this kind message
	//means send all my online peers to
	M_getonlinepeers

	//recieve this kind message
	//means some node send his block to me
	M_sendonlinepeers


	//recieve this kind message
	//means some node wants to join in the onlinepeers
	M_askforin

	//recieve this kind message
	//means some node wants to join in the onlinepeers
	M_returnforin

	//For SuperNode consensus
	//recieve this kind message
	//means a Leader node send me a new block
	C_sendblock
)

/// AskFor
type AskFor struct {
	Turn int
	Info *PeerInfo
	Sig  []byte
}

///

func (m *Message)EncodeMessage() (bytes.Buffer,error){
	var content bytes.Buffer
	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(m)
	return content,err
}

func Get_M_getall(config *SuperNodeConfig,height int) *Message {
	var askfor AskFor
	askfor.Turn = height
	askfor.Info = &PeerInfo{
		Province: config.Province,
		PublicKey: hex.EncodeToString(config.PublicKey),
		//Invanlid field
		Addr:      "",
		Rpcport:   config.RpcPort,
		Protocol:  "tcp",
	}
	askfor.Sig,_ = Signature(config.PrivateKey,[]byte("GetBlock"))

	var content bytes.Buffer
	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(askfor)
	if err != nil {
		fmt.Println("[*]Get_M_getall Encode Get Error",err)
		return nil
	}


	return &Message{
		MessageType: M_getall,
		Parameters:  content.Bytes(),
	}
}

func Get_M_getblock(config *SuperNodeConfig,height int) *Message {
	var askfor AskFor
	askfor.Turn = height
	askfor.Info = &PeerInfo{
		Province: config.Province,
		PublicKey: hex.EncodeToString(config.PublicKey),
		Addr:      "",
		Rpcport:   config.RpcPort,
		Protocol:  "tcp",
	}
	askfor.Sig,_ = Signature(config.PrivateKey,[]byte("GetBlock"))

	var content bytes.Buffer
	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(askfor)
	if err != nil {
		fmt.Println("[*]Get_M_getblock Encode Get Error",err)
		return nil
	}

	return &Message{
		MessageType: M_getblock,
		Parameters:  content.Bytes(),
	}

}

func Get_M_sendblock(blockstream []byte) *Message {
	return &Message{
		MessageType: M_sendblock,
		Parameters:  blockstream,
	}
}

func Get_M_getonlinepeers(config *SuperNodeConfig,turn int) *Message {
	var askfor AskFor
	askfor.Turn = turn
	askfor.Info = &PeerInfo{
		Province: config.Province,
		PublicKey: hex.EncodeToString(config.PublicKey),
		Addr:      "",
		Rpcport:   config.RpcPort,
		Protocol:  "tcp",
	}
	askfor.Sig,_ = Signature(config.PrivateKey,[]byte("What_is_new_onlinepeers"))

	var content bytes.Buffer
	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(askfor)
	if err != nil {
		fmt.Println("[*]Get_M_getonlinepeers Encode Get Error",err)
		return nil
	}

	return &Message{
		MessageType: M_getonlinepeers,
		Parameters:  content.Bytes(),
	}
}

func Get_M_sendonlinepeers(content []byte) *Message {
	return &Message{
		MessageType: M_sendonlinepeers,
		Parameters:  content,
	}
}

func Get_M_askforin(askfor AskFor) *Message {
	var content bytes.Buffer
	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(askfor)
	if err != nil {
		fmt.Println("[*]Get_M_askforin Encode Get Error",err)
		return nil
	}
	return &Message{
		MessageType: M_askforin,
		Parameters:  content.Bytes(),
	}
}

func Get_M_returnforin(config *SuperNodeConfig,turn int) *Message {
	var askfor AskFor
	askfor.Turn = turn
	askfor.Info = &PeerInfo{
		Province: config.Province,
		PublicKey: hex.EncodeToString(config.PublicKey),
		Addr:      "",
		Rpcport:   config.RpcPort,
		Protocol:  "tcp",
	}
	askfor.Sig,_ = Signature(config.PrivateKey,[]byte("YouAreIn"))

	var content bytes.Buffer
	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(askfor)

	if err != nil {
		fmt.Println("[*]Get_M_returnforin Encode Get Error",err)
		return nil
	}

	return &Message{
		MessageType: M_returnforin,
		Parameters:  content.Bytes(),
	}

}

func Get_C_sendblock(blockstream []byte) *Message {
	return &Message{
		MessageType: C_sendblock,
		Parameters:  blockstream,
	}
}