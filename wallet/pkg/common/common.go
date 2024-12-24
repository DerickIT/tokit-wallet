package common2

type Transaction struct {
	Hash   string
	From   string
	To     string
	Amount float64
}

type Blockchain interface {
	GetAddress() (string, error)
	Transfer(toAddress string, amount float64) (string, error)
	GetTransactionHistory(address string) ([]Transaction, error)
}
