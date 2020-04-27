## Introduction
* ```./dockerScipts``` is some script file for docker operation
* ```./IndexChain``` is the main code
  + ```./IndexChain/test``` is some API test for DeepSee and WeChat 
     - ```./IndexChain/test/testForDeepSee.go``` is some RPC call which used by DeepSee
     - ```./IndexChain/test/testWeChat.go``` is a RPC call to send a note to the IndexChain which used by WeChat
     
## Preparation
* Pull the image ledger-try
  + ```docker pull kingsleystc/ledger-try:lateset```
* run the script ```./dockerScripts/CreateDocker.sh```
  + This will create an isolate docker network at the subnet 192.168.1.0/24
  + Then Create 6 docker container where each one represent a SuperNode(peer) with specific container name and subnet
     - BeiJing ShangHai GuiZhou JiangXi ZheJiang XiZang
     - BeiJing and ShangHai are the DNSSeed which ipaddress and rpcport is specifid
  + Then mkdir a work directory for each container ```/root/SuperNode```
  
* ```./dockerScrips/StartNode.sh``` is the script for start 6 container
* ```./dockerScrips/attachX.sh``` is the script for attach the container

* ```./configX.json``` is the SuperNode Config for each node. ```createsuper``` will use it
  + This file will be copied to the container in the ```./bushujiedian.sh```

* ```./createNode.sh``` is the script used in the container , to teach how to create a SuperNode
  + There are two conditions , Init the Node First Time, Or Init the Node with a keystone
  
* ```./runNode.sh``` is the script used in the container , to start a SuperNode
  + ```./SuperNode runsuper -node InitialNode.config -peers Peers.json```

* ```bushujiedian.sh``` is the script for packing the codes and config files to the container, and build the project in the container

## How to Run

* [Out Container] Start the docker and container
* [Out Container] First running ```./bushujiedian.sh``` to deploy the codes and files 
* [In Each Container]
  + first run ```~/SuperNode/createNode``` script to create a new SuperNode
  + Then copy the publickey from the output to the ```./Peers.json```, make sure the province is specified correctly
  <font color="red">Tips: The complete Peers.json file is like that "All peer's province and pubkey is specified. As DNSSeed, BeiJing's and ShangHai's <address,rpcport,protocol> must be specified"</font>
* [Out Container] Running ```./bushupeers.sh``` to deploy the ```./Peers.json``` to each node
* [In Each Container] Running ```~/SuperNode/runNode.sh``` to start the Node
