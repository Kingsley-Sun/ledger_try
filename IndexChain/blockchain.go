package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"os"
	"strconv"
	"sync"
)

const dbFile = "blockchain_data.db"

type Blockchain struct {
	Latesheight int
	Database    *bolt.DB

	//control the block write
	Mutex sync.Mutex

	// key : height value : BlockCount
	Mempool map[int]*BlockCount
}

type BlockCount struct {
	block *Block
	count int
}

func (bc *Blockchain) PrintBlock(block *Block) {

	fmt.Println("------   Block is   ", "------")
	fmt.Println("block.Header.Timestamp", block.Header.Timestamp)
	fmt.Println("block.Header.Hash", block.Header.Hash)
	fmt.Println("block.Header.Height", block.Header.Height)
	fmt.Println("block.Header.Noteshash", block.Header.Noteshash)
	fmt.Println("block.Header.Prevhash", block.Header.Prevhash)
	fmt.Println("block.Header.Sig", block.Header.Sig)
	for _, note := range block.Notes {
		fmt.Println(" ")
		fmt.Println("     note.HashID", note.HashID)
		fmt.Println("     note.Timestamp", note.Timestamp)
		fmt.Println(" ")
	}
	fmt.Println("--------------------")
}

func (bc *Blockchain) PrintBlockchain() {
	height := bc.GetBestHeight()
	for i := 1; i <= height; i++ {
		block, err := bc.GetBlockByHeight(i)
		if err != nil {
			fmt.Println(err)
		}
		bc.PrintBlock(block)
	}
}

func (bc *Blockchain) GetBestHeight() int {
	return bc.Latesheight
}

func (bc *Blockchain) AddToMempool(block *Block) {
	if bc.Mempool[block.Header.Height] == nil {
		bc.Mempool[block.Header.Height] = &BlockCount{
			block: block,
			count: 1,
		}
	} else if bytes.Equal(bc.Mempool[block.Header.Height].block.Header.Hash[:], block.Header.Hash[:]) {
		bc.Mempool[block.Header.Height].count += 1
	} else {
		//temporily ingore this situation
		log.Panic("HARD FORK !!!!!!!")
	}
}

//this function only update dababase
//do not apply any verification
func (bc *Blockchain) UpdateDatabase(block *Block) {
	//do not mutex here
	//sequenced update only
	bytesBuffer := bytes.NewBuffer([]byte{})
	x := int64(block.Header.Height)
	binary.Write(bytesBuffer, binary.BigEndian, &x)
	height_byte := bytesBuffer.Bytes()
	err := bc.Database.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("blocks"))
		//blockInDb := b.Get(block.Header.Hash)
		//if blockInDb != nil {
		//	return nil
		//}
		err := b.Put(height_byte, block.Header.Hash[:])
		if err != nil {
			log.Panic(err)
		}
		blockdata := block.Serialize()
		err = b.Put(block.Header.Hash[:], blockdata)
		if err != nil {
			log.Panic(err)
		}
		err = b.Put([]byte("latest"), height_byte)
		if err != nil {
			log.Panic(err)
		}
		return nil
	})
	//after updating database successfully
	//update the blockchain.LatestHeight
	bc.Latesheight = block.Header.Height
	if err != nil {
		log.Panic(err)
	}
}

func (bc *Blockchain) GetNotesByBlockHeight(height int) ([]*Note, error) {
	block, err := bc.GetBlockByHeight(height)
	return block.Notes, err
}

func (bc *Blockchain) GetBlockByHeight(height int) (*Block, error) {
	var block Block
	bytesBuffer := bytes.NewBuffer([]byte{})
	x := int64(height)
	binary.Write(bytesBuffer, binary.BigEndian, &x)
	height_byte := bytesBuffer.Bytes()

	err := bc.Database.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("blocks"))
		blockhash := b.Get(height_byte)
		blockdata := b.Get(blockhash)
		if blockdata == nil {
			return errors.New("Block is not found.")
		}
		block = *DeserializeBlock(blockdata)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &block, nil
}

func (bc *Blockchain) GetNotesByBlockHash(blockhash []byte) []*Note {
	block, _ := bc.GetBlockByHash(blockhash)
	return block.Notes
}

// GetBlock finds a block by its hash and returns it
func (bc *Blockchain) GetBlockByHash(blockHash []byte) (*Block, error) {
	var block Block

	err := bc.Database.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("blocks"))

		blockData := b.Get(blockHash)

		if blockData == nil {
			return errors.New("Block is not found.")
		}

		block = *DeserializeBlock(blockData)

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &block, nil
}

func LoadLocalBlockChain() *Blockchain {
	var bc *Blockchain
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		//not exist
		db, err := bolt.Open(dbFile, 0600, nil)
		if err != nil {
			log.Panic(err)
		}

		err = db.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucket([]byte("blocks"))
			if err != nil {
				log.Panic(err)
			}
			return nil
		})

		bc = &Blockchain{
			Latesheight: 0,
			Database:    db,
			Mutex:       sync.Mutex{},
			Mempool:     make(map[int]*BlockCount),
		}
	} else {
		//exist
		db, err := bolt.Open(dbFile, 0600, nil)
		if err != nil {
			log.Panic(err)
		}
		var height int
		err = db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("blocks"))
			latestheight := b.Get([]byte("latest"))
			bytesBuffer := bytes.NewBuffer(latestheight)
			var latest int64
			binary.Read(bytesBuffer, binary.BigEndian, &latest)
			strInt64 := strconv.FormatInt(latest, 10)
			intInt, _ := strconv.Atoi(strInt64)
			height = intInt
			return nil
		})
		if err != nil {
			log.Panic(err)
		}
		bc = &Blockchain{
			Latesheight: height,
			Database:    db,
			Mutex:       sync.Mutex{},
			Mempool:     make(map[int]*BlockCount),
		}
	}
	return bc
}
