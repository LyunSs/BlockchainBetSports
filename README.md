# Blockchain in GO for BET SPORT

This project is a hands-on guide to building a basic blockchain implementation in Go, following the tutorial provided [here](https://medium.com/coinmonks/my-blockchain-in-go-8e2d1a853a84)

The tutorial covers the creation of a simple blockchain with the following features:

1. Block Structure: Each block in the blockchain contains
   - an index
   - timestamp
   - data
   - previous hash
   - current hash.
2. Blockchain Structure: The blockchain itself is represented as a slice of blocks.
3. Proof of Work (PoW): The consensus mechanism used is PoW, where miners must find a hash value that meets certain criteria (difficulty level) to add a new block to the chain.
4. HTTP Server: The blockchain is exposed via a simple HTTP server that handles API requests.

## API Endpoints

### POST Endpoints

#### 1. Register a new Bet

- **Endpoint:** `/bet`
- **Method:** `POST`
- **Description:** Adds a new bet to the pending bets pool.
- **Note:** If another node mines a block before this bet is included and they are on the same network, the bet will not be registered because it is not in the pending bets for everyone.
- **Example:**
  ```sh
  curl -X POST -H "Content-Type: application/json" -d '{
    "playername": "JohnDoe",
    "matchid": "match123",
    "teamonescore": 1,
    "teamtwoscore": 2
  }' http://localhost:9000/bet
  ```

#### 2. Register and Broadcast a new Bet

- **Endpoint:** `/bet/broadcast`
- **Method:** `POST`
- **Description:** Adds a new bet to the pending bets pool and broadcasts it to the network.
- **Example:**
  ```sh
  curl -X POST -H "Content-Type: application/json" -d '{
    "playername": "JohnDoe",
    "matchid": "match123",
    "teamonescore": 1,
    "teamtwoscore": 2
  }' http://localhost:9000/bet/broadcast
  ```

#### 3. Register a new Node

- **Endpoint:** `/register-node`
- **Method:** `POST`
- **Description:** Registers a new node in the network.
- **Note:** Registering a node without broadcasting creates a smaller network and is not shared with others. This endpoint is primarily for test purposes or for attempting to manipulate the blockchain.
- **Example:**
  ```sh
  curl -X POST -H "Content-Type: application/json" -d '{
    "newNodeUrl": "http://localhost:9001"
  }' http://localhost:9000/register-node
  ```

#### 4. Register and Broadcast a new Node

- **Endpoint:** `/register-and-broadcast-node`
- **Method:** `POST`
- **Description:** Registers a new node in the network and broadcasts it to other nodes.
- **Example:**
  ```sh
  curl -X POST -H "Content-Type: application/json" -d '{
    "newnodeurl": "http://localhost:9001"
  }' http://localhost:9000/register-and-broadcast-node
  ```

### GET Endpoints

#### 1. Get the entire Blockchain

- **Endpoint:** `/blockchain`
- **Method:** `GET`
- **Description:** Returns the entire blockchain.
- **Example:**
  ```sh
  curl -X GET http://localhost:9000/blockchain
  ```

#### 2. Mine a new Block

- **Endpoint:** `/mine`
- **Method:** `GET`
- **Description:** Mines a new block using the pending bets and adds it to the blockchain.
- **Example:**
  ```sh
  curl -X GET http://localhost:9000/mine
  ```

#### 3. Consensus Algorithm

- **Endpoint:** `/consensus`
- **Method:** `GET`
- **Description:** Implements the consensus algorithm to achieve blockchain consistency across the network.
- **Example:**
  ```sh
  curl -X GET http://localhost:9000/consensus
  ```

#### 4. Get Bets for a specific Match

- **Endpoint:** `/match/{matchId}`
- **Method:** `GET`
- **Description:** Retrieves all bets for a specific match.
- **Example:**
  ```sh
  curl -X GET http://localhost:9000/match/match123
  ```

#### 5. Get Bets for a specific Player

- **Endpoint:** `/player/{playerName}`
- **Method:** `GET`
- **Description:** Retrieves all bets made by a specific player.
- **Example:**
  ```sh
  curl -X GET http://localhost:9000/player/JohnDoe
  ```

You can run the code just by cloning the repo and using

```
go run main.go 9000
```

To multiple instance, you just have to change the port, and then you can try to add node :

```
go run main.go 9001
```

## With Docker (Optional)

You can also use Docker to run the blockchain.

**Please note: This setup currently works only with a single node, and multi-node support is coming in future updates.**

```
docker build -t blockchainbetsports_image .
```

Then, to forward the

```
docker run -p 9000:9000 --name <container_name> blockchainbetsports_image
```

# Credits

This project was developed by LyunSs following the tutorial by Mauricio M. Ribeiro
