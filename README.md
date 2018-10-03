# goNooNet

This module is a base network implementation of Noosphere nodes in NZT-shard.
It provides REST endpoints for work with transactions, wallets and other blockchain stuff.
List of endpoints will be expanding in future.
Also goNooNet provides simple TCP-sockets for data exchange between nodes of shard.

For more information read Docs.

## REST endpoints

### Wallet endpoints
* **/wallet/transaction** Endpoint for making simple transactions   
* **/wallet/getBalance** Endpoint for requesting user balance
* **/wallet/tranStatus** Endpoint for checking transaction status  

### Blockhain endpoints
* **/blockchain/getHeight** Endpoint for requesting current blockchain height
* **/blockchain/getTran** Endpoint for getting transaction by key
* **/blockchain/getBlock** Endpoint for getting block by height
* **/blockchain/getVersion** Endpoint for getting version of node software

## Default ports

Default port for REST is 5000  
Default port for TCP is 3333
