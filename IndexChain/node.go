package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/gob"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

type SuperNode struct {
	Peers *Peers

	//For Consensus
	ConsensusTurn  int
	ConsensusPeer  *Peer
	ConsensusBlock *Block

	//For Ask For in Onlinepeers
	LeaderTurn    int
	LeaderCommits int

	//control Mempool write and read
	Mutex sync.Mutex

	Config *SuperNodeConfig
}

type SuperNodeConfig struct {
	//note pool
	//get note from NormalNode and store it for building a block
	Mempool *Mempool

	//ca
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte

	//information
	Province string `json:"province"`
	IpAddr   string `json:"address"`
	RpcPort  string `json:"rpcport"`
}

func newAccount() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

	return *private, pubKey
}

//load node from file
func GetSuper(nodeconfig string, peersconfig string) (*SuperNode, error) {
	if _, err := os.Stat(nodeconfig); os.IsNotExist(err) {
		return nil, err
	}

	nodeContent, err := ioutil.ReadFile(nodeconfig)
	if err != nil {
		return nil, err
	}

	var supernodeconfig SuperNodeConfig
	gob.RegisterName("Curve", elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(nodeContent))
	err = decoder.Decode(&supernodeconfig)
	if err != nil {
		log.Panic(err)
	}

	//for inializing
	if supernodeconfig.Mempool.NotesMap == nil {
		supernodeconfig.Mempool.NotesMap = make(map[string]*Note)
	}

	supernode := &SuperNode{
		Peers:          GetDNSSeed(peersconfig),
		ConsensusTurn:  -1,
		ConsensusPeer:  nil,
		ConsensusBlock: nil,
		LeaderCommits:  0,
		LeaderTurn:     0,
		Mutex:          sync.Mutex{},
		Config:         &supernodeconfig,
	}

	return supernode, nil
}

//save node to file
func (n *SuperNode) SaveToFile() {
	var content bytes.Buffer
	gob.RegisterName("Curve", elliptic.P256())
	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(*n.Config)
	if err != nil {
		log.Panic(err)
	}
	fileName := "InitialNode.config"
	err = ioutil.WriteFile(fileName, content.Bytes(), 0644)
	if err != nil {
		log.Panic(err)
	}
}

func Signature(privatekey ecdsa.PrivateKey, hash []byte) ([]byte, error) {
	rr, ss, err := ecdsa.Sign(rand.Reader, &privatekey, hash[:])
	if err != nil {
		return nil, err
	}
	var sig []byte
	sig = append(sig, rr.Bytes()...)
	sig = append(sig, ss.Bytes()...)
	return sig, nil
}

func VerifySig(pubKey []byte, signature []byte, message []byte) bool {
	//TODO
	return true
}
