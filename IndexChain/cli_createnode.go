package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func CreateSuperNode(configfile string) error {
	//information
	s := SuperNodeConfig{}
	data_bytes, err := ioutil.ReadFile(configfile)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data_bytes, &s)
	if err != nil {
		return err
	}

	//ca
	private, pub := newAccount()
	s.PrivateKey = private
	s.PublicKey = pub
	s.Mempool = &Mempool{}

	//Persistence
	node := &SuperNode{
		Config: &s,
	}
	node.SaveToFile()
	fmt.Println("Your Province is : ", node.Config.Province)
	fmt.Println("Your Public key is : ")
	fmt.Println(hex.EncodeToString(node.Config.PublicKey))

	/*
		sss := hex.EncodeToString(node.Config.PublicKey)
		pubbytes,_ := hex.DecodeString(sss)
		fmt.Println(pubbytes)
	*/

	fmt.Println("Gather it with other province to peers.config")
	return nil
}
