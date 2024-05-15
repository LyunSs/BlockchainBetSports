package bet

//Bet represents a bet
type Bet struct {
	PlayerName   string `json:"playername" validate:"required"`
	MatchID      string `json:"matchid" validate:"required"`
	TeamOneScore int    `json:"teamonescore" validate:"required,min=0"`
	TeamTwoScore int    `json:"teamtwoscore" validate:"required,min=0"`
}

//Bets is an array of Bet
type Bets []Bet

//Block ...
type Block struct {
	Index             int    `json:"index"`
	Timestamp         int64  `json:"timestamp"`
	Bets              Bets   `json:"bets"`
	Nonce             int    `json:"nonce"`
	Hash              string `json:"hash"`
	PreviousBlockHash string `json:"previousblockhash"`
}

//Blocks is an array of Block
type Blocks []Block

//Blockchain ...
type Blockchain struct {
	Chain        Blocks   `json:"chain"`
	PendingBets  Bets     `json:"pending_bets"`
	NetworkNodes []string `json:"network_nodes"`
}

//BlockData is used in hash calculations
type BlockData struct {
	Index string
	Bets  Bets
}

