package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/gob"
	"io/ioutil"
	"log"
	"math/big"
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

type KeyStone struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
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
func (c *SuperNodeConfig) SaveToFile() {
	var content bytes.Buffer
	gob.RegisterName("Curve", elliptic.P256())
	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(*c)
	if err != nil {
		log.Panic(err)
	}
	fileName := "InitialNode.config"
	err = ioutil.WriteFile(fileName, content.Bytes(), 0644)
	if err != nil {
		log.Panic(err)
	}

	keystone := &KeyStone{
		PrivateKey: c.PrivateKey,
		PublicKey:  c.PublicKey,
	}

	var keycontent bytes.Buffer
	gob.RegisterName("Curve", elliptic.P256())
	keyencoder := gob.NewEncoder(&keycontent)
	err = keyencoder.Encode(*keystone)
	if err != nil {
		log.Panic(err)
	}
	keyfile := "keystone"
	err = ioutil.WriteFile(keyfile, content.Bytes(), 0644)
	if err != nil {
		log.Panic(err)
	}

}


func Signature(privatekey ecdsa.PrivateKey, hash []byte) ([]byte, error) {
	rr, ss, err := ecdsa.Sign(rand.Reader, &privatekey, hash)
	if err != nil {
		return nil, err
	}
	var sig []byte
	sig = append(sig, rr.Bytes()...)
	sig = append(sig, ss.Bytes()...)
	return sig, nil
}

func VerifySig(pubKey []byte, signature []byte, message []byte) bool {
	curve := elliptic.P256()
	X := hashToInt(pubKey[:32], curve)
	Y := hashToInt(pubKey[32:], curve)
	clientPub := &ecdsa.PublicKey{
		curve,
		X,
		Y,
	}

	r := signature[:32]
	s := signature[32:]

	rr := new(big.Int)
	ss := new(big.Int)
	rr.SetBytes(r)
	ss.SetBytes(s)

	var right = ecdsa.Verify(clientPub, message, rr, ss)
	return right
}

func hashToInt(hash []byte, c elliptic.Curve) *big.Int {
	orderBits := c.Params().N.BitLen()
	orderBytes := (orderBits + 7) / 8
	if len(hash) > orderBytes {
		hash = hash[:orderBytes]
	}

	ret := new(big.Int).SetBytes(hash)
	excess := len(hash)*8 - orderBits
	if excess > 0 {
		ret.Rsh(ret, uint(excess))
	}
	return ret
}
