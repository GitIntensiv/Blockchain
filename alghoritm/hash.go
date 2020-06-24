package alghoritm

import (
	"blockchain/model"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"log"
)

func HashBlock(block model.Block) (string, error) {
	str, err := json.Marshal(block)
	if err != nil {
		log.Println("Error in creating JSON")
		return string(-1), err
	}
	shaResult := sha512.Sum512(str)
	src := shaResult[:]
	dst := make([]byte, hex.EncodedLen(len(src)))
	return string(hex.Encode(dst, src)), nil
}

func HashString(str string) (int64, error) {
	shaResult := sha512.Sum512([]byte(str))
	src := shaResult[:]
	dst := make([]byte, hex.EncodedLen(len(src)))
	return int64(hex.Encode(dst, src)), nil
}
