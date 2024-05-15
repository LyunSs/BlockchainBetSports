package bet

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

const DIFFICULTY int = 3

//Controller ...
type Controller struct {
	blockchain     *Blockchain
	currentNodeURL string
}

//ResponseToSend ...
type ResponseToSend struct {
	Note  string `json:"note,omitempty"`
	Error string `json:"error,omitempty"`
}

//GetBlockchain GET /blockchain
func (c *Controller) GetBlockchain(w http.ResponseWriter, r *http.Request) {
	data, _ := json.Marshal(c.blockchain)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	return
}

//RegisterBet POST /bet
func (c *Controller) RegisterBet(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body) // read the body of the request
	if err != nil {
		log.Println("Error RegisterBet reading body:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := r.Body.Close(); err != nil {
		log.Fatalln("Error RegisterBet", err)
	}
	var bet Bet
	if err := json.Unmarshal(body, &bet); err != nil { // unmarshall body contents as a type Candidate
		w.WriteHeader(http.StatusUnprocessableEntity) // unprocessable entity (statusCode 422)
		response := ResponseToSend{
			Error: "Invalid input data: " + err.Error(),
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Println("Error RegisterBet encoding response:", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	validate := validator.New()
	if err := validate.Struct(bet); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		var validationErrors string
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors += "Field " + err.StructField() + " is " + err.Tag() + "; "
		}
		response := ResponseToSend{
			Error: "Validation error: " + validationErrors,
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Println("Error RegisterBet encoding response:", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	

	success := c.blockchain.RegisterBet(bet) // registers the bet into the blockchain
	if !success {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	var resp ResponseToSend
	resp.Note = "Bet created successfully."
	data, _ := json.Marshal(resp)
	w.Write(data)
	return
}


//RegisterAndBroadcastBet POST /bet/broadcast
func (c *Controller) RegisterAndBroadcastBet(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body) // read the body of the request
	if err != nil {
		log.Println("Error RegisterBet", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := r.Body.Close(); err != nil {
		log.Println("Error RegisterBet", err)
	}
	var bet Bet
	if err := json.Unmarshal(body, &bet); err != nil { // unmarshall body contents as a type Candidate
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Println("Error RegisterBet unmarshalling data", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	success := c.blockchain.RegisterBet(bet) // registers the bet into the blockchain
	if !success {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// broadcast
	for _, node := range c.blockchain.NetworkNodes {
		if node != c.currentNodeURL {
			// call /register-node in node
			MakePostCall(node+"/bet", body)
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	var resp ResponseToSend
	resp.Note = "Bet created and broadcast successfully."
	data, _ := json.Marshal(resp)
	w.Write(data)
}

func (c *Controller) Mine(w http.ResponseWriter, r *http.Request){
	lastBlock := c.blockchain.GetLastBlock()
	previousBlockHash := lastBlock.Hash
	currentBlockData := BlockData{Index: strconv.Itoa(lastBlock.Index + 1), Bets: c.blockchain.PendingBets}
	currentBlockDataAsByteArray, _ := json.Marshal(currentBlockData)
	currentBlockDataAsStr := base64.URLEncoding.EncodeToString(currentBlockDataAsByteArray)

	nonce := c.blockchain.ProofOfWork(previousBlockHash, currentBlockDataAsStr, DIFFICULTY)
	blockHash := c.blockchain.HashBlock(previousBlockHash, currentBlockDataAsStr, nonce)
	newBlock := c.blockchain.CreateNewBlock(nonce, previousBlockHash, blockHash)
	blockToBroadcast, _ := json.Marshal(newBlock)

	for _, node := range c.blockchain.NetworkNodes { //Dans la fonction range, ca renvoie l'index et la valeur de la liste
		if node != c.currentNodeURL {
			// call /receive-new-block in node
			MakePostCall(node+"/receive-new-block", blockToBroadcast)
		}
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	var resp ResponseToSend
	resp.Note = "New block mined and broadcast successfully" 
	data, _ := json.Marshal(resp)
	w.Write(data)
	return

}


//BroadcastNode broadcasting node
func BroadcastNode(newNode string, nodes []string) {
	for _, node := range nodes {
		if node != newNode {
			var registerNodesJSON = []byte(`{"newnodeurl":"` + newNode + `"}`)

			// call /register-node in node
			MakePostCall(node+"/register-node", registerNodesJSON)
		}
	}
}

//MakeCall ...
func MakeCall(mode string, url string, jsonStr []byte) interface{} {
	// call url in node
	req, err := http.NewRequest(mode, url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error in call " + url)
		log.Println(err)
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	var returnValue interface{}
	if err := json.Unmarshal(respBody, &returnValue); err != nil { // unmarshal body contents as a type Candidate
		if err != nil {
			log.Fatalln("Error "+url+" unmarshalling data", err)
			return nil
		}
	}
	log.Println(returnValue)
	return returnValue
}

func (c *Controller) ReceivedNewBlock(w http.ResponseWriter, r *http.Request){
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error ReceivingNewBlock reading body:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := r.Body.Close(); err != nil {
		log.Println("Error ReceivingNewBlock")
	}
	var blockReceived Block
	if err := json.Unmarshal(body, &blockReceived); err != nil { // unmarshall body contents as a type Candidate
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Fatalln("Error RegisterNode unmarshalling data", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	var resp ResponseToSend
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	// append block to blockchain
	if c.blockchain.CheckNewBlockHash(blockReceived) {
		resp.Note = "New Block received and accepted."
		c.blockchain.PendingBets = Bets{}
		c.blockchain.Chain = append(c.blockchain.Chain, blockReceived)
	} else {
		resp.Note = "New Block rejected."
	}

	data, _ := json.Marshal(resp)
	w.Write(data)
	return

}

//MakePostCall ...
func MakePostCall(url string, jsonStr []byte) {
	// call url in POST
	MakeCall("POST", url, jsonStr)
}


func (c *Controller) RegisterNode(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body) // read the body of the request
	if err != nil {
		log.Println("Error RegisterNode reading body:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := r.Body.Close(); err != nil {
		log.Fatalln("Error RegisterNode", err)
	}
	var node struct {
		NewNodeURL string `json:"newNodeUrl"`
	}
	if err := json.Unmarshal(body, &node); err != nil { // unmarshall body contents as a type Candidate
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Fatalln("Error RegisterNode unmarshalling data", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	var resp ResponseToSend
	if node.NewNodeURL != c.currentNodeURL {
		success := c.blockchain.RegisterNode(node.NewNodeURL) // registers the node into the blockchain
		if !success {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	resp.Note = "Node registered successfully."
	data, _ := json.Marshal(resp)
	w.Write(data)
	return
}


func (c *Controller) RegisterAndBroadcastNode(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body) // read the body of the request
	if err != nil {
		log.Fatalln("Error RegisterNode", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := r.Body.Close(); err != nil {
		log.Fatalln("Error RegisterNode", err)
	}
	var node struct {
		NewNodeURL string `json:"newnodeurl"`
	}
	if err := json.Unmarshal(body, &node); err != nil { // unmarshall body contents as a type Candidate
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Fatalln("Error RegisterNode unmarshalling data", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	var resp ResponseToSend
	success := c.blockchain.RegisterNode(node.NewNodeURL) // registers the node into the blockchain
	if !success {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// broadcast
	BroadcastNode(node.NewNodeURL, c.blockchain.NetworkNodes)

	// register all nodes in new node
	allNodes := append(c.blockchain.NetworkNodes, c.currentNodeURL)
	payload, err := json.Marshal(allNodes)
	registerBulkJSON := []byte(payload)
	MakePostCall(node.NewNodeURL+"/register-nodes-bulk", registerBulkJSON)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	resp.Note = "Node registered successfully."
	data, _ := json.Marshal(resp)
	w.Write(data)
	return
}


//RegisterNodesBulk POST /register-nodes-bulk
func (c *Controller) RegisterNodesBulk(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body) // read the body of the request
	if err != nil {
		log.Fatalln("Error RegisterNodesBulk", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := r.Body.Close(); err != nil {
		log.Fatalln("Error RegisterNodesBulk", err)
	}
	var allNodes []string
	if err := json.Unmarshal(body, &allNodes); err != nil { // unmarshall body contents as a type Candidate
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Fatalln("Error RegisterNodesBulk unmarshalling data", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	for _, node := range allNodes {
		if node != c.currentNodeURL {
			success := c.blockchain.RegisterNode(node) // registers the node into the blockchain
			if !success {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
	}
	var resp ResponseToSend
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	resp.Note = "Bulk registration successful."
	data, _ := json.Marshal(resp)
	w.Write(data)
	return
}

func (c *Controller) Consensus(w http.ResponseWriter, r *http.Request) {
	maxChainLength := 0
	var longestChain *Blockchain
	var resp ResponseToSend
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	for _, node := range c.blockchain.NetworkNodes {
		if node != c.currentNodeURL {
			// call /blockchain in node
			// call url in node
			req, err := http.NewRequest("GET", node+"/blockchain", nil)
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)

			if err != nil {
				log.Println("Error retrieving blockchain")
				log.Println(err)
			}
			defer resp.Body.Close()
			respBody, err := io.ReadAll(resp.Body)
			var chain *Blockchain
			if err := json.Unmarshal(respBody, &chain); err != nil { // unmarshal body contents as a type Candidate
				if err != nil {
					log.Println("Error unmarshalling data", err)
				}
			}
			if chain != nil {
				chainLength := len(chain.Chain)
				if maxChainLength < chainLength {
					maxChainLength = chainLength
					longestChain = chain
				}
			}
		}
	}
	if maxChainLength > len(c.blockchain.Chain) && longestChain.ChainIsValid() {
		c.blockchain.Chain = longestChain.Chain
		c.blockchain.PendingBets = longestChain.PendingBets

		resp.Note = "This chain has been replaced."
	} else {
		resp.Note = "This chain has not been replaced."
	}

	w.WriteHeader(http.StatusOK)
	data, _ := json.Marshal(resp)
	w.Write(data)
	return

}
//ChainIsValid Used by consensus algorithm
func (b *Blockchain) ChainIsValid() bool {
	for i := 1; i < len(b.Chain); i++{
		currentBlock := b.Chain[i]
		prevBlock := b.Chain[i-1]
		// Création des données du bloc courant
		currentBlockData := BlockData{Index: strconv.Itoa(prevBlock.Index + 1), Bets: currentBlock.Bets}
		currentBlockDataAsByteArray, _ := json.Marshal(currentBlockData)
		currentBlockDataAsStr := base64.URLEncoding.EncodeToString(currentBlockDataAsByteArray)
		// Vérification du hash du bloc courant avec la difficulté
		blockHash := b.HashBlock(prevBlock.Hash, currentBlockDataAsStr, currentBlock.Nonce)

		if !isValidBlockHash(blockHash, DIFFICULTY) {
			return false
		}

		if currentBlock.PreviousBlockHash != prevBlock.Hash {
			return false
		}
 
	}

	genesisBlock := b.Chain[0]
	correctNonce := genesisBlock.Nonce == 100
	correctPreviousBlockHash := genesisBlock.PreviousBlockHash == "0"
	correctHash := genesisBlock.Hash == "0"
	correctBets := len(genesisBlock.Bets) == 0

	return (correctNonce && correctPreviousBlockHash && correctHash && correctBets)
}
func isValidBlockHash(hash string, difficulty int) bool {
	prefix := strings.Repeat("0", difficulty)
	return strings.HasPrefix(hash, prefix)
}

func (c *Controller) GetBetsForMatch(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) //pour extraire les variables de la requête HTTP 'r'
	matchID := strings.ToLower(vars["matchId"])

	bets := c.blockchain.GetBetsForMatch(matchID)
	w.WriteHeader(http.StatusOK)
	data, _ := json.Marshal(bets)
	w.Write(data)
	return
}

//GetBetsForPlayer GET /player/{playerName}
func (c *Controller) GetBetsForPlayer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	playerName := strings.ToLower(vars["playerName"])

	bets := c.blockchain.GetBetsForPlayer(playerName)
	w.WriteHeader(http.StatusOK)
	data, _ := json.Marshal(bets)
	w.Write(data)
	return
}