package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"
)

func (c *CLI) RunSuper(nodeconfig string, peersconfig string) {
	superNode, err := GetSuper(nodeconfig, peersconfig)
	if err != nil {
		log.Panic(err)
	}
	superNode.Run()
}

func (s *SuperNode) Run() {
	fmt.Println("SuperNode Run()")

	var wg sync.WaitGroup

	blockchain := LoadLocalBlockChain()
	// Need to Generate Genius Block

	c := make(chan bool)
	//net server
	//handler all network communicating

	wg.Add(1)
	go s.NetServer(blockchain, c, &wg)


	wg.Add(1)
	go s.RpcServer(blockchain, c, &wg)


	wg.Add(1)
	//maintain the Peers state
	go s.Peers.HandlePeers(s.Config, c, &wg)

	s.Peers.SendMess_GetAll(s.Config,blockchain.GetBestHeight())

	go func() {
		for {
			fmt.Println("****************** Node And Chain State **************************")
			s.Peers.PrintAllPeers()
			fmt.Println("Database Height : ", blockchain.GetBestHeight())
			s.Config.Mempool.PrintMempool()
			blockchain.PrintBlockchain()
			fmt.Println("*******************         end         **************************")
			time.Sleep(time.Second * 10)
		}
	}()

	wg.Add(1)
	go s.SyncToLateset(blockchain,c,&wg)

	time.Sleep(time.Second * 9)
	s.AskForIn(c)

	wg.Add(1)
	go s.Routine(blockchain, c, &wg)


	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill)

	signal := <-interrupt
	fmt.Println("Got signal:", signal)
	close(c)

	//wait all server finish safely
	// wg.Wait()

	//Save blockchain state and mempool
	s.SaveToFile()

	fmt.Println("Exit Successfully!")
}

func (s *SuperNode) SyncToLateset(bc *Blockchain, c <-chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	var mutex sync.Mutex

	go func() {
		verifyMin := 2
		for {
			mutex.Lock()
			//start update block mempool
			for {
				nowheight := bc.GetBestHeight()
				if bc.Mempool[nowheight+1] == nil || bc.Mempool[nowheight+1].count < verifyMin {
					s.Peers.SendMess_GetBlock(s.Config, nowheight+1)
					break
				} else {
					bc.Mutex.Lock()
					bc.UpdateDatabase(bc.Mempool[nowheight+1].block)
					bc.Mutex.Unlock()
					delete(bc.Mempool, nowheight+1)

					//remove the local note in the block
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
			}

			mutex.Unlock()
			time.Sleep(time.Second * 1)
		}
	}()

	<-c
	mutex.Lock()
}

func (s *SuperNode) AskForIn(cc <-chan bool) {
	info := &PeerInfo{
		Province:  s.Config.Province,
		PublicKey: hex.EncodeToString(s.Config.PublicKey),
		Addr:      "",
		Rpcport:   s.Config.RpcPort,
		Protocol:  "tcp",
	}
	sig, _ := Signature(s.Config.PrivateKey, []byte("AskForIn"))
	askfor := AskFor{
		Turn: s.LeaderTurn,
		Info: info,
		Sig:  sig,
	}

	for {
		var num int

		s.Peers.Mutex.Lock()
		if len(s.Peers.OnlinePeers) >= 3 {
			num = len(s.Peers.OnlinePeers) * 2 / 3
		} else {
			num = len(s.Peers.OnlinePeers)
		}
		s.Peers.Mutex.Unlock()

		ctx, _ := context.WithTimeout(context.Background(), time.Second*15)
		turn_over := make(chan bool)
		go func(ctx context.Context, count int) {
			s.LeaderCommits = 0
			for {
				select {
				case <-ctx.Done():
					fmt.Println("[*]Ask For Being A Leader Timeout!!!")
					turn_over <- false
					return
				default:
					if s.LeaderCommits >= count {
						fmt.Println("[*]Ask For Being A Leader Success!!!")
						turn_over <- true
						return
					}
				}
			}
		}(ctx, num)

		time.Sleep(time.Second * 3)
		s.Peers.SendMess_AskForIn(askfor)
		for {
			select {
			case s := <-turn_over:
				if s == true {
					return
				} else {
					goto NewTurn
				}
			case <-cc:
				return
			}
		}

	NewTurn:
		s.LeaderTurn += 1
		askfor.Turn += 1
	}
}
