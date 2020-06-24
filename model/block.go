package model

import "time"

type Block struct {
	Id                int
	Timestamp         time.Time
	Transactions      []Transaction
	Proof             int64
	PreviousBlockHash string
}
