# device-powerflow

## Trader service
Trader service is the main service. 
It relates to a identity that is participitating in market.

Every trader service will run a blockchain service 
(each blockchain instance will act as a Ledger of transactions in the market).

## Blockchain service
Blockchain is proxy for "view of market". 
It is a Ledger of the powerflow market.

To begin participating in the powerflow network they have to register as miner.

Blockchain service sends a request to registration server.

And Register server responds with peerlist of other blockchain holders 
on the network.

Then download the current state of chain from any of the fellow blockchain holder.
 
## Register service for Blockchain
Register service keeps track of nodes in the network.
1. Registration service Nodes in network
Registration service node begins by calling /registerSelfTo api.
2. Blockchain service Nodes in network
Blockchain service sends a request "register/blockchain" to a Register service.
And the register service responds with list of other nodes in the network.
This act as start point of execution, for every blockchain service.


##Run details
run command

``go run main.go <ip> <port> <label>``

example 

``go run main.go 127.0.0.1 6686 temper``

### How to initialze this application
1. Run Edgex devices
1. Run atleast one registration service
2. For a actor in marketplace (each actor/trader controls 1 EdgeX instance)
    - a. Run a trader service 
    - b. Run a blockchain service for this trader service (register req to a register service)
    
    
    

## Overview
### Each instance of powerflow contains the following:
1. edgex services for a powerflow instance
2. trader service to trade and conclude a transaction
3. blockchain service to maintain chain
