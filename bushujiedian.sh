tar -zcvf ./ledger.tar ./IndexChain

#SuperNode1
#1
#192.168.1.11
docker cp ledger.tar SuperNode1:/root/go/src/ledger.tar
docker cp config1.json SuperNode1:/root/SuperNode/config.json
docker cp runNode.sh SuperNode1:/root/SuperNode/runNode.sh
docker cp createNode.sh SuperNode1:/root/SuperNode/createNode.sh

docker exec -d SuperNode1 chmod 777 /root/SuperNode/runNode.sh
docker exec -d SuperNode1 chmod 777 /root/SuperNode/createNode.sh

docker exec -d SuperNode1 rm -rf /root/go/src/IndexChain
docker exec -d SuperNode1 tar -zxvf /root/go/src/ledger.tar -C /root/go/src/
docker exec -d SuperNode1 rm /root/go/src/ledger.tar
docker exec -d SuperNode1 go build -o /root/SuperNode/SuperNode IndexChain
docker exec -d SuperNode1 rm /root/SuperNode/blockchain_data.db

#SuperNode2
#2
#192.168.1.22
docker cp ledger.tar SuperNode2:/root/go/src/ledger.tar
docker cp config2.json SuperNode2:/root/SuperNode/config.json
docker cp runNode.sh SuperNode2:/root/SuperNode/runNode.sh
docker cp createNode.sh SuperNode2:/root/SuperNode/createNode.sh

docker exec -d SuperNode2 chmod 777 /root/SuperNode/runNode.sh
docker exec -d SuperNode2 chmod 777 /root/SuperNode/createNode.sh

docker exec -d SuperNode2 rm -rf /root/go/src/IndexChain/
docker exec -d SuperNode2 tar -zxvf /root/go/src/ledger.tar -C /root/go/src/
docker exec -d SuperNode2 rm /root/go/src/ledger.tar
docker exec -d SuperNode2 go build -o /root/SuperNode/SuperNode IndexChain
docker exec -d SuperNode2 rm /root/SuperNode/blockchain_data.db

#SuperNode3
#3
#192.168.1.33
docker cp ledger.tar SuperNode3:/root/go/src/ledger.tar
docker cp config3.json SuperNode3:/root/SuperNode/config.json
docker cp runNode.sh SuperNode3:/root/SuperNode/runNode.sh
docker cp createNode.sh SuperNode3:/root/SuperNode/createNode.sh

docker exec -d SuperNode3 chmod 777 /root/SuperNode/runNode.sh
docker exec -d SuperNode3 chmod 777 /root/SuperNode/createNode.sh

docker exec -d SuperNode3 rm -rf /root/go/src/IndexChain/
docker exec -d SuperNode3 tar -zxvf /root/go/src/ledger.tar -C /root/go/src/
docker exec -d SuperNode3 rm /root/go/src/ledger.tar
docker exec -d SuperNode3 go build -o /root/SuperNode/SuperNode IndexChain
docker exec -d SuperNode3 rm /root/SuperNode/blockchain_data.db

#SuperNode4
#4
#192.168.1.44
docker cp ledger.tar SuperNode4:/root/go/src/ledger.tar
docker cp config4.json SuperNode4:/root/SuperNode/config.json
docker cp runNode.sh SuperNode4:/root/SuperNode/runNode.sh
docker cp createNode.sh SuperNode4:/root/SuperNode/createNode.sh

docker exec -d SuperNode4 chmod 777 /root/SuperNode/runNode.sh
docker exec -d SuperNode4 chmod 777 /root/SuperNode/createNode.sh

docker exec -d SuperNode4 rm -rf /root/go/src/IndexChain/
docker exec -d SuperNode4 tar -zxvf /root/go/src/ledger.tar -C /root/go/src/
docker exec -d SuperNode4 rm /root/go/src/ledger.tar
docker exec -d SuperNode4 go build -o /root/SuperNode/SuperNode IndexChain
docker exec -d SuperNode4 rm /root/SuperNode/blockchain_data.db

#SuperNode5
#5
#192.168.1.55
docker cp ledger.tar SuperNode5:/root/go/src/ledger.tar
docker cp config5.json SuperNode5:/root/SuperNode/config.json
docker cp runNode.sh SuperNode5:/root/SuperNode/runNode.sh
docker cp createNode.sh SuperNode5:/root/SuperNode/createNode.sh

docker exec -d SuperNode5 chmod 777 /root/SuperNode/runNode.sh
docker exec -d SuperNode5 chmod 777 /root/SuperNode/createNode.sh

docker exec -d SuperNode5 rm -rf /root/go/src/IndexChain/
docker exec -d SuperNode5 tar -zxvf /root/go/src/ledger.tar -C /root/go/src/
docker exec -d SuperNode5 rm /root/go/src/ledger.tar
docker exec -d SuperNode5 go build -o /root/SuperNode/SuperNode IndexChain
docker exec -d SuperNode5 rm /root/SuperNode/blockchain_data.db

#SuperNode6
#6
#192.168.1.66
docker cp ledger.tar SuperNode6:/root/go/src/ledger.tar
docker cp config6.json SuperNode6:/root/SuperNode/config.json
docker cp runNode.sh SuperNode6:/root/SuperNode/runNode.sh
docker cp createNode.sh SuperNode6:/root/SuperNode/createNode.sh

docker exec -d SuperNode6 chmod 777 /root/SuperNode/runNode.sh
docker exec -d SuperNode6 chmod 777 /root/SuperNode/createNode.sh

docker exec -d SuperNode6 rm -rf /root/go/src/IndexChain/
docker exec -d SuperNode6 tar -zxvf /root/go/src/ledger.tar -C /root/go/src/
docker exec -d SuperNode6 rm /root/go/src/ledger.tar
docker exec -d SuperNode6 go build -o /root/SuperNode/SuperNode IndexChain
docker exec -d SuperNode6 rm /root/SuperNode/blockchain_data.db

rm ./ledger.tar
