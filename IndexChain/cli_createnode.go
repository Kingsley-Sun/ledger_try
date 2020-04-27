package main

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func CreateSuperNode(configfile string,keystone string) error {
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
	s.Mempool = &Mempool{}
	if keystone == "" {
		//Create A new Peer
		private, pub := newAccount()
		s.PrivateKey = private
		s.PublicKey = pub
	}else {
		if _, err := os.Stat(keystone); os.IsNotExist(err) {
			return err
		}
		nodeContent, err := ioutil.ReadFile(keystone)
		if err != nil {
			return err
		}

		var keystone KeyStone
		gob.RegisterName("Curve", elliptic.P256())
		decoder := gob.NewDecoder(bytes.NewReader(nodeContent))
		err = decoder.Decode(&keystone)
		if err != nil {
			return err
		}

		s.PrivateKey = keystone.PrivateKey
		s.PublicKey = keystone.PublicKey
	}

	//Persistence
	s.SaveToFile()

	fmt.Println("Your Province is : ", s.Province)
	fmt.Println("Your Public key is : ")
	fmt.Println(hex.EncodeToString(s.PublicKey))
	fmt.Println("Your Private key is : ")
	fmt.Println(s.PrivateKey.D.Bytes())

	fmt.Println("Gather it with other province to peers.config")
	return nil
}
