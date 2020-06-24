package alghoritm

import (
	"log"
	"strconv"
)

func ProofOfWork(lastProof int64, currentProof int64) bool {
	str := strconv.FormatInt(lastProof, 16) + strconv.FormatInt(currentProof, 16)
	hash, err := HashString(str)
	if err != nil {
		log.Fatal("Can not check proof")
	}
	return strconv.FormatInt(hash, 16)[:2] == "01"
}

func FindProof(lastProof int64) int64 {
	proof := 0
	for ; ProofOfWork(lastProof, int64(proof)); proof++ {
	}
	return int64(proof)
}
