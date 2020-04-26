package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"time"
)

type SuperNodeInBlock struct {
	PublicKey []byte
	Province  string
}

type Blockheader struct {
	Miner *SuperNodeInBlock
	//hash of Sig Prevhash Timestamp Height
	Hash [20]byte
	//sigature for NotesHash
	Sig       []byte
	Prevhash  [20]byte
	Noteshash [20]byte
	Timestamp int64
	Height    int
}

type Block struct {
	Header *Blockheader
	//information
	Notes []*Note
}

func (block *Block)Verify(peer *Peer) bool {
	if block.Header.Miner.Province != peer.Province {
		return false
	}
	if hex.EncodeToString(block.Header.Miner.PublicKey) != hex.EncodeToString(peer.PublicKey) {
		return false
	}
	//TODO
	return true
}

func (block *Block)VerifyMBlock() bool {
	//Used For Sendblock Message
	//TODO
	return true
}

// Serialize a block to stream
func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)
	if err != nil {
		fmt.Println("Block Serialize Wrong")
		return []byte{}
	}

	return result.Bytes()
}

// Deserializes a block from stream
func DeserializeBlock(d []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
		return nil
	}

	return &block
}

func NewBlock(s *SuperNode, previous *Block) *Block {
	var prevhash [20]byte
	var height int
	if previous != nil {
		prevhash = previous.Header.Hash
		height = previous.Header.Height + 1
	} else {
		prevhash = [20]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		height = 1
	}
	timeNow := time.Now().Unix()
	blockheader := &Blockheader{
		Miner: &SuperNodeInBlock{
			PublicKey: s.Config.PublicKey,
			Province:  s.Config.Province,
		},
		Prevhash:  prevhash,
		Timestamp: timeNow,
		Height:    height,
	}
	s.Mutex.Lock()
	notes := s.Config.Mempool.GetBlockNotes()
	s.Mutex.Unlock()

	noteshash := s.Config.Mempool.CalNotesHash(notes)
	signature, err := Signature(s.Config.PrivateKey,noteshash[:])
	if err != nil {
		return nil
	}

	blockheader.Noteshash = noteshash
	blockheader.Sig = signature

	blockhash := CalBlockHash(signature, prevhash[:])
	blockheader.Hash = blockhash

	block := &Block{
		Header: blockheader,
		Notes:  notes,
	}
	return block
}

func CalBlockHash(signature []byte, prevhash []byte) [20]byte {
	var content []byte
	content = append(content, signature...)
	content = append(content, prevhash...)
	return sha1.Sum(content)
}

func IsSameHeader(h1 Blockheader,h2 Blockheader) bool {
	//Verify Miner
	m1 := *(h1.Miner)
	m2 := *(h2.Miner)
	if hex.EncodeToString(m1.PublicKey) != hex.EncodeToString(m2.PublicKey) {
		return false
	}
	if m1.Province != m2.Province {
		return false
	}

	//Verify Hash
	if string(h1.Hash[:]) != string(h2.Hash[:]) {
		return false
	}

	//Verify Prevhash
	if string(h1.Prevhash[:]) != string(h2.Prevhash[:]) {
		return false
	}

	//Verify Noteshash
	if string(h1.Noteshash[:]) != string(h2.Noteshash[:]) {
		return false
	}

	//Verify Sig
	if string(h1.Sig[:]) != string(h2.Sig[:]) {
		return false
	}

	//Verify Timestamp
	if h1.Timestamp != h2.Timestamp {
		return false
	}

	//Verify Height
	if h1.Height != h2.Height {
		return false
	}
	return true
}

func IsSameBlock(b1 Block,b2 Block) bool {
	if IsSameHeader(*b1.Header,*b2.Header) == false {
		fmt.Println("Is Not Same Header")
		return false
	}

	if IsSameNotes(b1.Notes,b2.Notes) == false {
		fmt.Println("Is Not Same Notes")
		return false
	}
	return true
}


