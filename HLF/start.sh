#!/bin/bash
#
# Copyright IBM Corp All Rights Reserved
#
# SPDX-License-Identifier: Apache-2.0
#
# Exit on first error, print all commands.
set -ev

# don't rewrite paths for Windows Git Bash users
export MSYS_NO_PATHCONV=1

docker-compose -f docker-compose.yml down

docker-compose -f docker-compose.yml up -d
# ca.kyc.com orderer.kyc.com peer0.bnk1.kyc.com couchdb

# wait for Hyperledger Fabric to start
# incase of errors when running later commands, issue export FABRIC_START_TIMEOUT=<larger number>
export FABRIC_START_TIMEOUT=20
#echo ${FABRIC_START_TIMEOUT}
sleep ${FABRIC_START_TIMEOUT}


#channel creating
docker exec cli peer channel create -o orderer.kyc.com:7050 -c commonchannel -f /etc/hyperledger/configtx/commonchannel.tx
#channel created

#bnk1 joining channel
docker exec cli peer channel join -b commonchannel.block
#bnk1 joined channel
#bnk2 joining channel
docker exec -e "CORE_PEER_LOCALMSPID=bnk2msp" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/bnk2.kyc.com/users/Admin@bnk2.kyc.com/msp" -e "CORE_PEER_ADDRESS=peer0.bnk2.kyc.com:7051" cli peer channel join -b commonchannel.block
#bnk2 joined channel
#rgl1 joining channel
docker exec -e "CORE_PEER_LOCALMSPID=rgl1msp" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/rgl1.kyc.com/users/Admin@rgl1.kyc.com/msp" -e "CORE_PEER_ADDRESS=peer0.rgl1.kyc.com:7051" cli peer channel join -b commonchannel.block
#rgl1 joined channel

#all peers joined the commonchannel
