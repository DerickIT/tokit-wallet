package blockchain

import "crypto/ecdsa"

// BlockchainClient defines the interface for interacting with different blockchains.
type BlockchainClient interface {
	GetAddress(publicKey []byte) (string, error)
	Transfer(privateKey *ecdsa.PrivateKey, toAddress string, amount float64) (string, error)
	GetTransactionHistory(address string) ([]Transaction, error)
}

// Transaction represents a transaction on a blockchain.
type Transaction struct {
	Hash   string
	From   string
	To     string
	Amount float64
}
