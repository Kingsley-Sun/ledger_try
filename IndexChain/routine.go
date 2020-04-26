package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func (s *SuperNode) Routine(bc *Blockchain, cc <-chan bool, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Println("")
	fmt.Println("[*]Starting Routine")
	fmt.Println("")

	c := make(chan bool)

	var mutex sync.Mutex
	go func() {
		for {
			mutex.Lock()

			//phrase 1
			go Sync(c, 1)
			<-c

			s.Peers.Mutex.Lock()
			s.ConsensusTurn = 1

			//fmt.Println("[~]~~~~ Sync1 over ~~~~")

			//phrase 2
			go Sync(c, 2)
			<-c

			s.ConsensusTurn = 2

			var leaderpeer *Peer
			var isMe bool
			var prev *Block
			if bc.Latesheight == 0 {
				leaderpeer = ChooseLeader([]byte{0, 0, 0, 0}, s.Peers.OnlinePeers)
				prev = nil
			} else {
				prev, _ = bc.GetBlockByHeight(bc.GetBestHeight())
				hash := prev.Header.Hash[:4]
				leaderpeer = ChooseLeader(hash, s.Peers.OnlinePeers)
			}

			if s.Config.Province == leaderpeer.Province {
				isMe = true

				fmt.Println("[~]   Leader IS ME ")
				fmt.Println("[~]      Peer is    ")
				fmt.Println(leaderpeer.Province)
			} else {
				isMe = false

				fmt.Println("[~] Leader IS  NOT  ME")
				fmt.Println("[~]      Peer is    ")
				fmt.Println(leaderpeer.Province)
			}
			s.ConsensusPeer = leaderpeer

			//fmt.Println("[~]~~~~ Sync2 over ~~~~")

			go Sync(c, 3)
			<-c
			s.ConsensusTurn = 3


			if isMe {
				fmt.Println("[*]Generating A Block")

				//Waiting for other node prepare
				time.Sleep(time.Second * 1)

				var block *Block
				block = NewBlock(s, prev)

				s.ConsensusBlock = block
				s.Peers.SendMess_CsendBlock(block.Serialize())

			} else {
				//Waiting for block
				fmt.Println("[*]Waiting for block")
				s.ConsensusTurn = 30

				ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
				for {
					select {
					case <-ctx.Done():
						fmt.Println("[*]Waiting For Block Timeout. This turn can't build a block")
						goto RoundEnd
					case <-s.Peers.Chann:
						goto CommitPhrase
					}
				}
			}

		CommitPhrase:
			go Sync(c, 4)
			<-c
			s.ConsensusTurn = 4

			//Use Rpc
			if s.RequestCommits() {

				fmt.Println("[!!]Consensus Success")
				fmt.Println("[!!]Updating the Blockchain")

				//Update the BlockChain
				bc.Mutex.Lock()
				bc.UpdateDatabase(s.ConsensusBlock)
				bc.Mutex.Unlock()

				latesnotes, _ := bc.GetNotesByBlockHeight(bc.GetBestHeight())
				for _, note := range latesnotes {
					notehash := note.HashID
					if s.Config.Mempool.HasNote(notehash) == true {
						//delete it from note mempool
						s.Mutex.Lock()
						s.Config.Mempool.DeleteNote(notehash)
						s.Mutex.Unlock()
					}
				}

			}

			//fmt.Println("[~]~~~~ Sync4 over ~~~~")

		RoundEnd:
			go Sync(c,5)
			<-c
			//fmt.Println("[~]RoundEnd")

			s.ConsensusTurn = -1
			s.ConsensusBlock = nil
			s.ConsensusPeer = nil

			s.Peers.Mutex.Unlock()

			mutex.Unlock()
		}
	}()
	<-cc
	//quit safely
	mutex.Lock()
}

func Sync(c chan bool, phase int) {
	switch phase {
	case 1:
		for {
			if time.Now().Second() == 0 {
				//fmt.Println("Phase 1!!!!!!!!!")
				break
			}
		}
		c <- true
	case 2:
		for {
			if time.Now().Second() == 10 {
				//fmt.Println("Phsase 2!!!!!!!!!")
				break
			}
		}
		c <- true
	case 3:
		for {
			if time.Now().Second() == 13 {
				//fmt.Println("Phsase 3!!!!!!!!!")
				break
			}
		}
		c <- true
	case 4:
		for {
			if time.Now().Second() == 20 {
				//fmt.Println("Phsase 4!!!!!!!!!")
				break
			}
		}
		c <- true
	case 5:
		for {
			if time.Now().Second() == 40 {
				//fmt.Println("Phsase 5!!!!!!!!!")
				break
			}
		}
		c <- true
	}

}
