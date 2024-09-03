package bet

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

//RegisterBet registers a bet in our blockchain
func (b *Blockchain) RegisterBet(bet Bet) bool {
	log.Print("Registering following bet: ", bet)
	b.PendingBets = append(b.PendingBets, bet)
	return true
}

func contains(slice []string, element string) bool {
	for _, item := range slice {
		if item == element {
			return true
		}
	}
	return false
}


//RegisterNode registers a node in our blockchain
func (b *Blockchain) RegisterNode(nodeURL string) bool {
	if !contains(b.NetworkNodes, nodeURL){
		b.NetworkNodes = append(b.NetworkNodes, nodeURL)
	}
	return true
}

//CreateNewBlock ...
func (b *Blockchain) CreateNewBlock(nonce int, previousHash string, hash string) Block {
	newBlock := Block{
		Index:     len(b.Chain) + 1,
		Bets:      b.PendingBets,
		Timestamp: time.Now().UnixNano(),
		Nonce:     nonce,
		Hash:      hash, 
		PreviousBlockHash: previousHash}

	b.PendingBets = Bets{}
	b.Chain = append(b.Chain, newBlock)
	return newBlock
}

func (b *Blockchain) HashBlock(previousHash string, currentBlockData string, nonce int) string {

	h := sha256.New()
	strToHash := previousHash + currentBlockData +strconv.Itoa(nonce)
	h.Write([]byte(strToHash))
	hashed := base64.URLEncoding.EncodeToString(h.Sum(nil))
	return hashed
	
}

func (b *Blockchain) ProofOfWork(previousHash string, currentBlockData string, difficulty int) int {
	nonce := 0
	targetPrefix := strings.Repeat("0", difficulty) // Préfixe cible en fonction de la difficulté

	for {
		hash := b.HashBlock(previousHash, currentBlockData, nonce)

		// Vérifier si le hachage commence par le préfixe cible
		if strings.HasPrefix(hash, targetPrefix) {
			return nonce // Renvoyer le nonce trouvé
		}

		nonce++ // Incrémenter le nonce pour la prochaine itération
	}
}


func (b *Blockchain) GetLastBlock() Block {
	return b.Chain[len(b.Chain) - 1]
} 

func (b *Blockchain) CheckNewBlockHash(newBlock Block) bool {
	lastBlock := b.GetLastBlock()
	correctHash := lastBlock.Hash == newBlock.PreviousBlockHash
	correctIndex := (lastBlock.Index + 1) == newBlock.Index
	return	(correctHash && correctIndex)

}

//GetBetsForMatch ...
func (b *Blockchain) GetBetsForMatch(matchID string) Bets {
	matchBets := Bets{}
	i := 0
	chainLength := len(b.Chain)
	for i < chainLength {
		block := b.Chain[i]
		betsInBlock := block.Bets
		j := 0
		betsLength := len(betsInBlock)
		for j < betsLength {
			bet := betsInBlock[j]
			fmt.Println("les bets", bet)
			if bet.MatchID == matchID {
				matchBets = append(matchBets, bet)
			}
			j = j + 1
		}
		i = i + 1
	}
	return matchBets
}

//GetBetsForPlayer ...
func (b *Blockchain) GetBetsForPlayer(playerName string) Bets {
	matchBets := Bets{}
	i := 0
	chainLength := len(b.Chain)
	for i < chainLength {
		block := b.Chain[i]
		betsInBlock := block.Bets
		j := 0
		betsLength := len(betsInBlock)
		for j < betsLength {
			bet := betsInBlock[j]
			if bet.PlayerName == playerName {
				matchBets = append(matchBets, bet)
			}
			j = j + 1
		}
		i = i + 1
	}
	return matchBets
}