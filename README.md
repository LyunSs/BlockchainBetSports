# Blockchain in GO for BET SPORT

This project is a hands-on guide to building a basic blockchain implementation in Go, following the tutorial provided [here](https://medium.com/coinmonks/my-blockchain-in-go-8e2d1a853a84)

The tutorial covers the creation of a simple blockchain with the following features:

1. Block Structure: Each block in the blockchain contains an index, timestamp, data, previous hash, and current hash..
2. Blockchain Structure: The blockchain itself is represented as a slice of blocks.
3. Proof of Work (PoW): The consensus mechanism used is PoW, where miners must find a hash value that meets certain criteria (difficulty level) to add a new block to the chain.
4. HTTP Server: The blockchain is exposed via a simple HTTP server that handles API requests.

You can run the code just by cloning the repo and using 

```
go run main.go 9000
```

Or you can use Docker to run the blockchain :

```
docker build -t blockchainbetsports_image .
```

Then,

```
docker run -p 9000:9000 --name <container_name> blockchainbetsports_image
```

# Credits

This project was developed by LyunSs following the tutorial by Mauricio M. Ribeiro
