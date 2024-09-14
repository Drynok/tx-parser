package model

type Transaction struct {
	Hash        string `json:"hash"`
	From        string `json:"from"`
	To          string `json:"to"`
	Value       string `json:"value"`
	Timestamp   int64  `json:"timestamp"`
	BlockNumber int64  `json:"block_number"`
}

type Block struct {
	Number       int           `json:"number"`
	Transactions []Transaction `json:"transactions"`
}
