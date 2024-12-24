package common2

type Transaction struct {
	Hash      string
	From      string
	To        string
	Amount    float64
	Timestamp string
}

type Blockchain interface {
	GenerateKey() ([]byte, error)
	GetAddress() (string, error)
	Transfer(privateKeyHex string, toAddress string, amount float64) (string, error)
	GetTransactionHistory(address string) ([]Transaction, error)
}
