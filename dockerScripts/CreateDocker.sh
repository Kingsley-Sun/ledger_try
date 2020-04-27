##get the docker images ledger-try
docker pull kingsleystc/ledger-try:lateset

docker network create --subnet=192.168.1.0/24 ledger_try_network

docker run -itd --name SuperNode1 --net ledger_try_network --ip 192.168.1.11 kingsleystc/ledger-try:latest /bin/bash
docker start SuperNode1
docker exec -d SuperNode1 mkdir /root/SuperNode
docker stop SuperNode1

docker run -itd --name SuperNode2 --net ledger_try_network --ip 192.168.1.22 kingsleystc/ledger-try:latest /bin/bash
docker start SuperNode2
docker exec -d SuperNode2 mkdir /root/SuperNode
docker stop SuperNode2

docker run -itd --name SuperNode3 --net ledger_try_network --ip 192.168.1.33 kingsleystc/ledger-try:latest /bin/bash
docker start SuperNode3
docker exec -d SuperNode3 mkdir /root/SuperNode
docker stop SuperNode3

docker run -itd --name SuperNode4 --net ledger_try_network --ip 192.168.1.44 kingsleystc/ledger-try:latest /bin/bash
docker start SuperNode4
docker exec -d SuperNode4 mkdir /root/SuperNode
docker stop SuperNode4

docker run -itd --name SuperNode5 --net ledger_try_network --ip 192.168.1.55 kingsleystc/ledger-try:latest /bin/bash
docker start SuperNode5
docker exec -d SuperNode5 mkdir /root/SuperNode
docker stop SuperNode5

docker run -itd --name SuperNode6 --net ledger_try_network --ip 192.168.1.66 kingsleystc/ledger-try:latest /bin/bash
docker start SuperNode6
docker exec -d SuperNode6 mkdir /root/SuperNode
docker stop SuperNode6
