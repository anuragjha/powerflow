# **control-powerflow** 

## is a control mechanism for powerflow instance
Business logic for energy trade

domain of trader - 
1. Consume(Demand) Trader will read the devices requirement state (from device)
2. Create requirement tx and send to miner (Require tx) - Demand Tx

1. Supply Trader will read for requirement tx in the blockchain by polling chain on self instance
2. And will send offer tx for that requirement to miners (willing to Sell tx) - Supply Offer Tx

1. Consume(Demand) Trader will read chain to find seller(s) offer and then create a buy request for 1 offer
2. And will send buy tx for that offer to miners (willing to Buy Offer tx) - Buy Bid Tx

[A Device can only be part of 1 exchange activity at a time]
1. Supply Trader will read chain to find willing to Buy offer tx from consumer
2. And will send sell tx to miners - Supply Tx
