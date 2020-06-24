package model

type ChainResponse struct {
	CurrentChain  []Block
	CurrentLength int
}

type MainingResponse struct {
	Index             int
	Message           string
	PreviousBlockHash string
	Proof             int64
	Transactions      []Transaction
}

type RegisterNodeResponse struct {
	Message string
	Nodes   map[string]string
}

type ConsensusResponse struct {
	Message string
	Blocks  []Block
}
