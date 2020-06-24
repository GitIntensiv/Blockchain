package model

import (
	"log"
	"net/http"
	"time"
)

type Blockchain struct {
	Chain              []Block
	CurrentTransaction []Transaction
	Nodes              map[string]string
}

func (blockchain *Blockchain) Init() {
	blockchain.CreateNewBlock("1", 200)
}

func (blockchain *Blockchain) CreateNewBlock(previousBlockHash string, proof int64) Block {
	var block = Block{Id: len(blockchain.Chain) + 1, Timestamp: time.Now(),
		Transactions: blockchain.CurrentTransaction, Proof: proof, PreviousBlockHash: previousBlockHash}
	blockchain.CurrentTransaction = nil
	blockchain.Chain = append(blockchain.Chain, block)
	return block
}

func (blockchain *Blockchain) CreateNewTransaction(sender string, recipient string, amount int) int {
	var transaction = Transaction{
		Sender: sender, Recipient: recipient, Amount: amount,
	}
	blockchain.AddNewTransaction(transaction)
	return len(blockchain.Chain)
}

func (blockchain *Blockchain) AddNewTransaction(transaction Transaction) int {
	blockchain.CurrentTransaction = append(blockchain.CurrentTransaction, transaction)
	return len(blockchain.Chain)
}

func (blockchain *Blockchain) GetLastBlock() Block {
	return blockchain.Chain[len(blockchain.Chain)-1]
}

func (blockchain *Blockchain) RegisterNode(name string, address string) {
	blockchain.Nodes[name] = address
}

func (blockchain *Blockchain) CheckNodesState() {
	for key, v := range blockchain.Nodes {
		response, err := http.Get(v)
		if err != nil {
			log.Println("Error in creating connection")
			return
		}
		if response.StatusCode == 404 {
			delete(blockchain.Nodes, key)
		}
	}
}
