# BC-Parser 

Eth indexer and parser challenge 

## build 

Make sure you have go1.22 installed. 

And just run


```bash
make build
```

## run 

Run this 

```bash
./parser-api
```


## Usage 

The service creates a REST API to interact with the provided parser interface


### Get current block 

Returns the last block in decimal the indexer processed 

```bash
curl http://localhost:8081/api/v1/last_block
```

Response 

```json
{"last_block":19590123}
```

### Subscribe address

Add an address to the watch list for to/from transactions.
If the address is already on the list endpoint returns succesfully `ok`.
If is the address is invalid for eth an error is returned.


```bash
curl -X POST http://localhost:8081/api/v1/subscribe\?address\=0x1c9fce6dd765a22040d500019ada91acce65b5d2
```

response:

```json
{"status": "ok"}
```

in case of error:

http code: 400 
```json
{"error":"invalid address"}
```


### Get transactions 

Get a list of transactions for and address.
Returns error if the privided address is not subscribed.
Returns error if the provided address is invalid.


```bash
curl http://localhost:8081/api/v1/transactions\?address\=0x1c9fce6dd765a22040d500019ada91acce65b5d2
```

response:
```json
{
  "transactions": [
    {
      "blockHash": "0x602128504604c7369abb788836fbe7a1c3ab16a8b0d884b60c596db6711460aa",
      "blockNumber": "0x12aebe8",
      "from": "0x6907894f656b95d67e380349a5edc1f75bc45b8c",
      "gas": "0x5208",
      "gasPrice": "0x7c94f11fe",
      "maxFeePerGas": "0xa1d35fe21",
      "maxPriorityFeePerGas": "0x24bbc6a",
      "hash": "0x5363e5b048a841874afe8ce1389e9db4a4a662882ef9345270fd43b368198103",
      "input": "0x",
      "nonce": "0x4c",
      "to": "0x1c9fce6dd765a22040d500019ada91acce65b5d2",
      "transactionIndex": "0xaf",
      "value": "0x354a6ba7a18000",
      "type": "0x2",
      "chainId": "0x1",
      "v": "0x0",
      "r": "0xd089a65636f574fe4023f5c0da550d59f59a9ccce1b8a40935295846006eb6a4",
      "s": "0x14268bebbbde9c715b547b02aa09ef7e55759be2a2ee61e02504fe44ccdc168f"
    },
    {
      "blockHash": "0x9332e6354ba81948c46ecd17496c58b2f6895460be780be21528b2ff4a9abe1f",
      "blockNumber": "0x12aec0a",
      "from": "0x1c9fce6dd765a22040d500019ada91acce65b5d2",
      "gas": "0x5208",
      "gasPrice": "0x84ad299a8",
      "maxFeePerGas": "0xa3e9ab800",
      "maxPriorityFeePerGas": "0x3938700",
      "hash": "0xa7ea5928d5a65f89f70585b28dc0c70626b80c8265ddd7c6f9760113a83faed2",
      "input": "0x",
      "nonce": "0x32",
      "to": "0xf52605c7b778563a5a9144ef4dc53b57463ca2c7",
      "transactionIndex": "0x56",
      "value": "0x32d143c2fe6a00",
      "type": "0x2",
      "chainId": "0x1",
      "v": "0x1",
      "r": "0x45d1b6c46937b0d0b08c67dd78b25fb272b7a19f6971907c97aef843cff0c2b6",
      "s": "0x1a6bd8d046b38b72ecb0cde1036b8b261087ee25e5a7d237263cbc0650de0ebf",
      "yParity": "0x1"
    }
  ]
}
```

in case of error:

http code: 400 
```json
{"error":"invalid address"}
```

or 

```json
{"error":"address not subscribed"}
```



## Improvements
- read config from env (inital block, rpc endpoint, etc)
- add more test
- json rpc client could be improved (didnt want to wast time there)

## Packages 

### client 
Is the implementation of JSONRPC methods used for the indexer

### Indexer 
Is the process to index blocks form the blockchain 
Its a single thread process that sequencially read blocks untile its up-to-date to the last block. 
The output is sent to a channel to be used by other services. 
It could by improved to concurrently read many blocks to speed up the process. 

### Storage 
Is the in-memory storage, its modeled so i can also use it to hold the subscription addresses

### Observer 
Is the code that implements the proposed interface (Parser) for the challenge 

### App 
contains the handlers for the rest api






