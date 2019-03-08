#!/bin/bash
#
# Copyright IBM Corp All Rights Reserved
#
# SPDX-License-Identifier: Apache-2.0
#
# Exit on first error
set -e
IMAGETAG="1.3.0"
THIRDPARTYIMAGETAG="0.4.13"
# don't rewrite paths for Windows Git Bash users
export MSYS_NO_PATHCONV=1
starttime=$(date +%s)
LANGUAGE=${1:-"golang"}
CC_SRC_PATH=github.com/kyc/go
if [ "$LANGUAGE" = "node" -o "$LANGUAGE" = "NODE" ]; then
	CC_SRC_PATH=/opt/gopath/src/github.com/kyc/node
fi

# clean the keystore
rm -rf ./hfc-key-store

# clean the keystore
# rm -rf ./config

# launch network; create channel and join peer to channel
#cd ../hlf
IMAGE_TAG=$IMAGETAG THIRDPARTYIMAGE_TAG=$THIRDPARTYIMAGETAG  ./start.sh

# Now launch the CLI container in order to install, instantiate chaincode
# and prime the ledger with our 10 cars
IMAGE_TAG=$IMAGETAG THIRDPARTYIMAGE_TAG=$THIRDPARTYIMAGETAG docker-compose -f ./docker-compose.yml up -d cli

printf "\n chain code install and instanciation in bnk1 start\n"
docker exec -e "CORE_PEER_LOCALMSPID=bnk1msp" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/bnk1.kyc.com/users/Admin@bnk1.kyc.com/msp" cli peer chaincode install -n kyccc -v 1.0 -p "$CC_SRC_PATH" -l "$LANGUAGE"
printf "\n chain code install in bnk1 end \n"
docker exec -e "CORE_PEER_LOCALMSPID=bnk1msp" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/bnk1.kyc.com/users/Admin@bnk1.kyc.com/msp" cli peer chaincode instantiate -o orderer.kyc.com:7050 -C commonchannel -n kyccc -l "$LANGUAGE" -v 1.0 -c '{"Args":[""]}' -P "OR ('bnk1msp.member','bnk2msp.member','rgl1msp.member')"
printf "\n chain code install and instanciation in bnk1 end. \n"

# sleep 4
# # docker exec -e "CORE_PEER_LOCALMSPID=bnk1msp" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/bnk1.kyc.com/users/Admin@bnk1.kyc.com/msp" cli peer chaincode invoke -o orderer.kyc.com:7050 -C commonchannel -n kyccc -c '{"function":"initLedger","Args":[""]}'

# printf "\n chain code install and instanciation in bnk2 start \n"
# docker exec -e "CORE_PEER_LOCALMSPID=bnk2msp" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/bnk2.kyc.com/users/Admin@bnk2.kyc.com/msp" -e "CORE_PEER_ADDRESS=peer0.bnk2.kyc.com:7051" cli peer chaincode install -n kyccc -v 1.0 -p "$CC_SRC_PATH" -l "$LANGUAGE"
# printf "\n chain code install in bnk2 end \n"
# docker exec -e "CORE_PEER_LOCALMSPID=bnk2msp" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/bnk2.kyc.com/users/Admin@bnk2.kyc.com/msp" -e "CORE_PEER_ADDRESS=peer0.bnk2.kyc.com:7051" cli peer chaincode instantiate -o orderer.kyc.com:7050 -C commonchannel -n kyccc -l "$LANGUAGE" -v 1.0 -c '{"Args":[""]}' -P "OR ('bnk1msp.member','bnk2msp.member','rgl1msp.member')"
# printf "\n chain code install and instanciation in bnk2 end \n"

# sleep 4

# printf "\n chain code install and instanciation in rgl start \n"
# docker exec -e "CORE_PEER_LOCALMSPID=rgl1msp" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/rgl1.kyc.com/users/Admin@rgl1.kyc.com/msp" -e "CORE_PEER_ADDRESS=peer0.rgl1.kyc.com:7051" cli peer chaincode install -n kyccc -v 1.0 -p "$CC_SRC_PATH" -l "$LANGUAGE"
# printf "\n chain code install in rgl end \n"
# docker exec -e "CORE_PEER_LOCALMSPID=rgl1msp" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/rgl1.kyc.com/users/Admin@rgl1.kyc.com/msp" -e "CORE_PEER_ADDRESS=peer0.rgl1.kyc.com:7051" cli peer chaincode instantiate -o orderer.kyc.com:7050 -C commonchannel -n kyccc -l "$LANGUAGE" -v 1.0 -c '{"Args":[""]}' -P "OR ('bnk1msp.member','bnk2msp.member','rgl1msp.member')"
# printf "\n chain code install and instanciation in rgl end \n"

printf "\n Admins enrollment start \n"

 node ./../Integration/enrollAdmin.js bnk1
# node ./../Integration/enrollAdmin.js bnk2
# node ./../Integration/enrollAdmin.js rgl1

printf "\n Admins enrollment end \n"


printf "\n User registrations  start \n"

  node ./../Integration/registerUser bnk1
#  node ./../Integration/registerUser bnk2
#  node ./../Integration/registerUser rgl1

printf "\n User registrations end \n"


printf 'Block Fabric network ready to connect'
printf 'One admin and one user created for each org.'

#to look at the data
#cat  /var/hyperledger/production/ledgersData/chains/chains/commonchannel/blockfile_000000