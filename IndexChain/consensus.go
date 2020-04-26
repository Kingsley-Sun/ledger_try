package main

import (
	"encoding/hex"
	"fmt"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"sort"
)

func (s *SuperNode) handle_consensus(mess *Message,blockchain *Blockchain) {
	if s.ConsensusTurn < 0 {
		return
	}
	//fmt.Println("[~]In handle_consensus")
	switch mess.MessageType {
	case C_sendblock:
		s.handle_csendblock(mess,blockchain)
	}
}

func (s *SuperNode)RequestCommits() bool {
	if s.ConsensusTurn != 4 {
		return false
	}
	commitpeer := 0
	for _, peer := range s.Peers.OnlinePeers {

		conn, err := net.Dial("tcp", peer.Addr+":10000")
		if err != nil {
			fmt.Println("[*]RequestCommits Dial ",peer.Addr," error:",err)
			continue
		}

		client := rpc.NewClientWithCodec(jsonrpc.NewClientCodec(conn))

		var reply string

		args := hex.EncodeToString(s.ConsensusBlock.Serialize())

		//fmt.Println("[~]Send Block args")
		//fmt.Println(s.ConsensusBlock.Serialize())

		err = client.Call("NodeRequestService.RequestCommits", args, &reply)

		if err != nil {
			fmt.Println("[*]NodeRequestService.RequestCommits Call ",peer.Addr," error",err)
			continue
		}

		if reply == "Commit" {
			commitpeer += 1
		}
	}

	var num int
	if len(s.Peers.OnlinePeers) >= 3 {
		num = len(s.Peers.OnlinePeers) * 2 / 3
	} else {
		num = len(s.Peers.OnlinePeers)
	}

	if commitpeer >= num {
		//fmt.Println("[~]True Commitpeer is ",commitpeer)
		return true
	}else {
		//fmt.Println("[~]False Commitpeer is ",commitpeer)
		return false
	}

}

func (s *SuperNode) handle_csendblock(mess *Message,bc *Blockchain) {
	//fmt.Println("[~]In handle_csendblock")
	if s.ConsensusTurn != 30 {
		return
	}
	block := DeserializeBlock(mess.Parameters)
	if block.Header.Height != bc.GetBestHeight() + 1{
		fmt.Println("[*]Recieve an illegal Block With a wrong Height")
		return
	}
	if block.Verify(s.ConsensusPeer) == false {
		return
	}

	s.ConsensusBlock = block
	s.Peers.Chann <- true

}

func (p *Peers) SendMess_CsendBlock(blockstream []byte) {
	//fmt.Println("[~]In SendMess_CsendBlock")
	message := Get_C_sendblock(blockstream)

	content,err := message.EncodeMessage()
	if err != nil {
		fmt.Println("[*]SendMess_CsendBlock EncodeMessage Error",err)
		return
	}

	var despeers []*Peer
	despeers = append(despeers, p.OnlinePeers...)
	SendToPeer(despeers, content.Bytes())
}


func GetProvinceId(province string) int {
	result := 0
	for _,v := range province {
		result += int(v)
	}
	return result
}

func ChooseLeader(prev []byte, onlinepeer []*Peer) *Peer {
	var hashint int
	for _, v := range prev {
		hashint += int(v)
	}
	var idindex []int
	peersid := make(map[int]*Peer)
	for _, v := range onlinepeer {
		vid := GetProvinceId(v.Province)
		idindex = append(idindex, vid)
		peersid[vid] = v
	}
	//fmt.Println("[~]IDindex[] is ", idindex)
	sort.Ints(idindex)

	target := hashint % len(idindex)
	index := idindex[target]
	fmt.Println("Leader is ", peersid[index].Province)

	return peersid[index]
}

