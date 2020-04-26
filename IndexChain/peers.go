package main

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

type Peer struct {
	Province  string
	PublicKey []byte
	Addr      string
	Rpcport   string
	Protocol  string
}

type Peers struct {
	//Save all legal pubkey actually
	Superpeers []*Peer

	Mutex       sync.Mutex
	OnlinePeers []*Peer

	Commits int

	//use for no leaders waiting block
	Chann chan bool
}

/////  For load peers from config.file
type PeerInfo struct {
	Province  string `json:"province"`
	PublicKey string `json:"publickey"`
	Addr      string `json:"address"`
	Rpcport   string `json:"rpcport"`
	Protocol  string `json:"protocol"`
}

type PeersInfo struct {
	PeersInfo []*PeerInfo `json:"peers"`
}

///// For AskForIn struct

func (p *Peers) PrintAllPeers() {
	//fmt.Println("Superpeers:")
	//for _, v := range p.Superpeers {
	//	fmt.Println("-----")
	//	fmt.Println(v.Province, v.Addr, v.Rpcport, hex.EncodeToString(v.PublicKey), v.Protocol)
	//	fmt.Println("-----")
	//}
	fmt.Println("--------   PeersOnline   --------")
	for _, v := range p.OnlinePeers {
		fmt.Println("+++++")
		fmt.Println(v.Province, v.Addr, v.Rpcport, hex.EncodeToString(v.PublicKey), v.Protocol)
		fmt.Println("+++++")
	}
}

func IsSameProvince(source string,cmp string) bool {
	index := strings.Index(cmp,source)
	if index == 0 {
		return true
	}
	return false
}

//check pub whether in my superpeers
func (p *Peers) IsLegal(pub string,province string) bool {
	for _, peer := range p.Superpeers {
		if hex.EncodeToString(peer.PublicKey) == pub && IsSameProvince(peer.Province, province) == true {
			return true
		}
	}
	return false
}

func IsSamePeer(p1 *Peer,p2 *Peer) bool {
	if p1.Province != p2.Province {
		return false
	}
	if hex.EncodeToString(p1.PublicKey) != hex.EncodeToString(p2.PublicKey) {
		return false
	}
	return true
}

//If two peer all the same return true
func Contains(ps []*Peer,p *Peer) bool {
	for _,peer := range ps{
		if IsSamePeer(peer,p) == true {
			return true
		}
	}
	return false
}

/////
func GetPeer(pro string, pub string, addr string, port string, proto string) *Peer {
	pubbytes, err := hex.DecodeString(pub)
	if err != nil {
		log.Panic(err)
	}
	return &Peer{
		Province:  pro,
		PublicKey: pubbytes,
		Addr:      addr,
		Rpcport:   port,
		Protocol:  proto,
	}
}

func GetDNSSeed(peersconfig string) *Peers {
	if _, err := os.Stat(peersconfig); os.IsNotExist(err) {
		log.Panic(err)
	}

	content, err := ioutil.ReadFile(peersconfig)
	if err != nil {
		log.Panic(err)
	}
	var p PeersInfo
	err = json.Unmarshal(content, &p)

	var superpeers []*Peer
	var onlinepeers []*Peer
	for _, peerinfo := range p.PeersInfo {
		peer := GetPeer(peerinfo.Province, peerinfo.PublicKey, peerinfo.Addr, peerinfo.Rpcport, peerinfo.Protocol)
		superpeers = append(superpeers, peer)
		if peer.Addr != "" {
			onlinepeers = append(onlinepeers, peer)
		}
	}

	return &Peers{
		Superpeers:  superpeers,
		OnlinePeers: onlinepeers,
		Commits:     0,
		Mutex:       sync.Mutex{},
		Chann:       make(chan bool),
	}
}

func (p *Peers) HandlePeers(config *SuperNodeConfig, c <-chan bool, wg *sync.WaitGroup) {

	defer wg.Done()

	go func() {
		for {
			p.SendMess_GetOnlinPeers(config, 0)
			time.Sleep(time.Second * 5)
		}
	}()

	go func() {
		for {
			// maintain the heartbeat to all online peers
			time.Sleep(time.Second * 6)
			var removelist []int
			var ww sync.WaitGroup
			p.Mutex.Lock()
			for index := range p.OnlinePeers {
				ww.Add(1)
				peer := p.OnlinePeers[index]
				go func(index int, thisp *Peer, ww *sync.WaitGroup) {
					defer ww.Done()
					if p.IsAlive(thisp) == false {
						removelist = append(removelist, index)
					}
				}(index, peer, &ww)
			}
			ww.Wait()

			sort.Ints(removelist)
			for i, v := range removelist {
				fmt.Println("[*]A OnlinePeer Dead. Remove it : ",p.OnlinePeers[v-i])
				p.OnlinePeers = append(p.OnlinePeers[:v-i], p.OnlinePeers[v+1-i:]...)
			}
			p.Mutex.Unlock()
		}
	}()

	<-c
}

func (p *Peers) IsAlive(peer *Peer) bool {
	//Use Rpc
	conn, err := net.Dial("tcp", peer.Addr+":10000")
	if err != nil {
		fmt.Println("[*]IsAlive net.Dial ",peer.Addr," Error ",err)
		return false
	}

	client := rpc.NewClientWithCodec(jsonrpc.NewClientCodec(conn))

	var reply string
	args := "Are_You_Ok?"
	err = client.Call("NodeRequestService.HeartBeats", args, &reply)
	if err != nil {
		fmt.Println("[*]IsAlive Client.Call ",peer.Addr," Error",err)
		return false
	}

	if reply == "I_Am_Ok" {
		return true
	}else {
		return false
	}
}

func (p *Peers) SendMess_GetOnlinPeers(config *SuperNodeConfig, turn int) {
	message := Get_M_getonlinepeers(config, turn)
	if message == nil {
		return
	}

	content, err := message.EncodeMessage()
	if err != nil {
		fmt.Println("[*]SendMess_GetOnlinPeers EncodeMessage Got Error ", err)
		return
	}

	var despeers []*Peer
	despeers = append(despeers, p.OnlinePeers...)
	//fmt.Println("[~]SendMess_GetOnlinePeers")
	SendToPeer(despeers, content.Bytes())
}

func (p *Peers) ReturnMess_GetOnlinePeers(ps []*Peer, pro string, pubkey string, ip string, rpcport string) {
	peers := []*Peer{
		GetPeer(pro,pubkey,ip,rpcport,"tcp"),
	}

	var content bytes.Buffer
	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(ps)
	if err != nil {
		fmt.Println("[*]ReturnMess_GetOnlinePeers Encode Got Error ", err)
		return
	}

	message := Get_M_sendonlinepeers(content.Bytes())
	if message == nil {
		return
	}

	contents,errs := message.EncodeMessage()
	if errs != nil {
		fmt.Println("[*]ReturnMess_GetOnlinePeers EncodeMessage Got Error ", err)
		return
	}

	//fmt.Println("[~]In function ReturnMess_GetOnliePeers")
	SendToPeer(peers, contents.Bytes())

}

func (p *Peers) SendMess_AskForIn(askfor AskFor) {
	message := Get_M_askforin(askfor)
	if message == nil {
		return
	}

	content,err := message.EncodeMessage()
	if err != nil {
		fmt.Println("[*]SendMess_AskForIn EncodeMessage Got Error ", err)
		return
	}

	var despeers []*Peer
	despeers = append(despeers, p.OnlinePeers...)
	//fmt.Println("[~]In peers.go function AskForIn")
	SendToPeer(despeers, content.Bytes())

}

func (p *Peers) ReturnMess_AskForIn(config *SuperNodeConfig,pro string, pubkey string,  rpcport string, turn int,ip string) {
	//fmt.Println("[~]In function ReturnMess_AskForIn")
	peers := []*Peer{
		GetPeer(pro,pubkey,ip,rpcport,"tcp"),
	}

	message := Get_M_returnforin(config,turn)
	if message == nil {
		return
	}

	content,err := message.EncodeMessage()
	if err != nil {
		fmt.Println("[*]ReturnMess_AskForIn EncodeMessage Got Error ", err)
		return
	}
	SendToPeer(peers, content.Bytes())
}

func (p *Peers) SendMess_GetBlock(config *SuperNodeConfig,height int) {
	//send GetBlock Message to SuperNode
	message := Get_M_getblock(config,height)
	if message == nil {
		return
	}

	content,err := message.EncodeMessage()
	if err != nil {
		fmt.Println("[*]SendMess_GetBlock EncodeMessage Got Error ", err)
		return
	}

	var despeers []*Peer
	despeers = append(despeers, p.OnlinePeers...)

	//fmt.Println("[~]In peers.go function SendMess_GetBlock")
	SendToPeer(despeers, content.Bytes())
}

func (p *Peers) SendMess_GetAll(config *SuperNodeConfig,height int) {
	//send GetAllBlocks Message to SuperNode
	message := Get_M_getall(config,height)
	if message == nil {
		return
	}

	content,err := message.EncodeMessage()
	if err != nil {
		fmt.Println("[*]SendMess_GetAll EncodeMessage Got Error ", err)
		return
	}

	//here can use random to select some node
	var despeers []*Peer
	despeers = append(despeers, p.OnlinePeers...)

	//fmt.Println("[~]In peers.go function SendMess")
	SendToPeer(despeers, content.Bytes())
}

func (p *Peers) ReturnMess_GetBlock(block *Block,pro string, pubkey string, ip string, rpcport string) {
	peers := []*Peer{
		GetPeer(pro,pubkey,ip,rpcport,"tcp"),
	}

	block_stream := block.Serialize()
	message := Get_M_sendblock(block_stream)
	if message == nil {
		return
	}

	content,err := message.EncodeMessage()
	if err != nil {
		fmt.Println("[*]ReturnMess_GetBlock EncodeMessage Got Error ", err)
		return
	}
	//fmt.Println("[~]In function ReturnMess_GetBlock")
	SendToPeer(peers, content.Bytes())
}

func (p *Peers) ReturnMess_GetAll(blocks []*Block,pro string, pubkey string, ip string, rpcport string) {
	peers := []*Peer{
		GetPeer(pro,pubkey,ip,rpcport,"tcp"),
	}
	var wg sync.WaitGroup
	for _, blockIndex := range blocks {
		block := blockIndex
		wg.Add(1)
		go func() {
			block_stream := block.Serialize()
			message := Get_M_sendblock(block_stream)
			if message == nil {
				return
			}
			content,err := message.EncodeMessage()
			if err != nil {
				fmt.Println("[*]ReturnMess_GetAll EncodeMessage Got Error ", err)
				return
			}
			//fmt.Println("[~]In function ReturnMess_GetAll")
			SendToPeer(peers, content.Bytes())
			wg.Done()
		}()
	}
	wg.Wait()
}

//stream is a seriliazation of Message
func SendToPeer(peers []*Peer, stream []byte) {
	var wg sync.WaitGroup
	for index, _ := range peers {
		wg.Add(1)
		i := index
		go func() {
			conn, err := net.Dial(peers[i].Protocol, peers[i].Addr+":"+peers[i].Rpcport)
			if err == nil {
				//fmt.Println("[~]Dial success to", peers[i].Addr+":"+peers[i].Rpcport)
				defer conn.Close()
				conn.Write(stream)
				//fmt.Println("[~]Send to ", peers[i].Addr+":"+peers[i].Rpcport+" succeccfully")
			} else {
				fmt.Println("[*]SendTOPeer Net.Dial ",peers[i].Addr," has error ", err)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
