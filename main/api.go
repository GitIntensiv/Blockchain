package main

import (
	"blockchain/alghoritm"
	"blockchain/model"
	"encoding/json"
	"github.com/carlescere/scheduler"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

var name string
var blockchain model.Blockchain = model.Blockchain{}

func Maining(w http.ResponseWriter, r *http.Request) {
	lastBlock := blockchain.Chain[len(blockchain.Chain)-1]
	lastProof := lastBlock.Proof
	nextProof := alghoritm.FindProof(lastProof)
	blockchain.CreateNewTransaction("0", name, 1)
	hashBlock, err := alghoritm.HashBlock(lastBlock)
	if err != nil {
		log.Println("Can Not Calculate a Hash Of Block")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	newBlock := blockchain.CreateNewBlock(hashBlock, nextProof)
	response := model.MainingResponse{
		Index:             newBlock.Id,
		Message:           "New Block Forged",
		PreviousBlockHash: newBlock.PreviousBlockHash,
		Proof:             newBlock.Proof,
		Transactions:      newBlock.Transactions,
	}
	_ = json.NewEncoder(w).Encode(response)
}

func CreateTransactions(w http.ResponseWriter, r *http.Request) {
	var transaction model.Transaction
	request, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Incorrect Input Data")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_ = json.Unmarshal(request, &transaction)
	index := blockchain.AddNewTransaction(transaction)
	var response = "Transaction Will Be Added To Block " + string(index)
	_ = json.NewEncoder(w).Encode(response)
}

func GetChain(w http.ResponseWriter, r *http.Request) {
	info := model.ChainResponse{CurrentChain: blockchain.Chain, CurrentLength: len(blockchain.CurrentTransaction)}
	_ = json.NewEncoder(w).Encode(info)
}

func RegisterAsNode(w http.ResponseWriter, r *http.Request) {
	var nodes = make(map[string]string)
	request, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Incorrect Input Data")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_ = json.Unmarshal(request, &nodes)
	if nodes != nil {
		for k, v := range nodes {
			blockchain.RegisterNode(k, v)
		}
	} else {
		response := model.RegisterNodeResponse{Message: "No Blocks Were Added", Nodes: blockchain.Nodes}
		_ = json.NewEncoder(w).Encode(response)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func ResolveConflicts(w http.ResponseWriter, r *http.Request) {
	result, err := alghoritm.Consensus(&blockchain)
	var response model.ConsensusResponse
	if err != nil {
		log.Println("Some Errors Was Happened In Algorithm Of Consensus")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if result {
		response = model.ConsensusResponse{Message: "Chain Was Replaced", Blocks: blockchain.Chain}
	} else {
		response = model.ConsensusResponse{Message: "Chain Is Authoritative", Blocks: blockchain.Chain}
	}
	_ = json.NewEncoder(w).Encode(response)
}

func main() {
	blockchain.Init()
	_, _ = scheduler.Every(10).Seconds().Run(blockchain.CheckNodesState)
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/maining", Maining).Methods("GET")
	router.HandleFunc("/create/transaction", CreateTransactions).Methods("POST")
	router.HandleFunc("/chain", GetChain).Methods("GET")
	router.HandleFunc("/nodes/register", RegisterAsNode).Methods("POST")
	router.HandleFunc("/consensus", ResolveConflicts).Methods("POST")
	log.Fatal(http.ListenAndServe("localhost:8080", router))
}
