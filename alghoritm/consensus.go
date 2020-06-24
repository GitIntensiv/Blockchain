package alghoritm

import (
	"blockchain/model"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func CheckChain(chain []model.Block) bool {
	lastBlock := chain[0]
	for index := 1; index < len(chain); index++ {
		block := chain[index]
		hash, err := HashBlock(lastBlock)
		if err != nil {
			log.Println("Error in calculating hash")
			return false
		}
		if block.PreviousBlockHash != hash || !ProofOfWork(lastBlock.Proof, block.Proof) {
			return false
		}
		lastBlock = block
	}
	return true
}

func Consensus(blockchain *model.Blockchain) (bool, error) {
	var neighbors = blockchain.Nodes
	maxLength := len(blockchain.Chain)
	var chainResponse model.ChainResponse
	var newChain []model.Block = nil
	for _, url := range neighbors {
		response, err := http.Get(string(url) + "/chainResponse")
		if err != nil {
			log.Println("Error in creating connection")
			return false, err
		}
		if response.StatusCode == 200 {
			body, err := ioutil.ReadAll(response.Body)
			if err != nil {
				log.Println("Error in read answer")
				return false, err
			}
			_ = json.Unmarshal(body, &chainResponse)
			curLength := chainResponse.CurrentLength
			curChain := chainResponse.CurrentChain
			if curLength > maxLength && CheckChain(curChain) {
				maxLength = curLength
				newChain = curChain
			}
		}
		if newChain != nil {
			blockchain.Chain = newChain
			return true, nil
		}
	}
	return false, nil
}
