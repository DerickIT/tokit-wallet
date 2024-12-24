package wallet

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"os"

	common2 "tokit/wallet/pkg/common2"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
)

// Wallet manages multiple blockchain clients.
type Wallet struct {
	clients        map[string]common2.Blockchain
	clientFactory  func(chain string, rpcURL string) (common2.Blockchain, error)
	ethereumClient *ethclient.Client
}

// NewWallet creates a new wallet with the given client factory.
func NewWallet(clientFactory func(chain string, rpcURL string) (common2.Blockchain, error), ethereumClient *ethclient.Client) *Wallet {
	return &Wallet{
		clients:        make(map[string]common2.Blockchain),
		clientFactory:  clientFactory,
		ethereumClient: ethereumClient,
	}
}

// AddClient adds a new blockchain client to the wallet.
func (w *Wallet) AddClient(chain string) error {
	rpcURL := os.Getenv(chain + "_RPC_URL")
	if rpcURL == "" {
		rpcURL = "http://localhost:8545" // Default if not set
	}
	client, err := w.clientFactory(chain, rpcURL)
	if err != nil {
		return fmt.Errorf("failed to create client for %s: %w", chain, err)
	}
	w.clients[chain] = client
	return nil
}

// GetAddress returns the address for the given chain.
func (w *Wallet) GetAddress(chain string) (string, error) {
	client, ok := w.clients[chain]
	if !ok {
		return "", fmt.Errorf("client for %s not found", chain)
	}
	return client.GetAddress()
}

// Transfer transfers funds from the wallet to the given address on the specified chain.
func (w *Wallet) Transfer(chain string, privateKeyHex string, toAddress string, amount float64) (string, error) {
	client, ok := w.clients[chain]
	if !ok {
		return "", fmt.Errorf("client for %s not found", chain)
	}
	return client.Transfer(privateKeyHex, toAddress, amount)
}

// GetTransactionHistory returns the transaction history for the given address on the specified chain.
func (w *Wallet) GetTransactionHistory(chain string, address string) ([]common2.Transaction, error) {
	client, ok := w.clients[chain]
	if !ok {
		return nil, fmt.Errorf("client for %s not found", chain)
	}
	return client.GetTransactionHistory(address)
}

func (w *Wallet) deriveChildKey() (*ecdsa.PrivateKey, error) {
	// TODO: Implement actual key derivation logic based on the chain.
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

// SignTransaction signs an Ethereum transaction.
func (w *Wallet) SignTransaction(chain string, toAddress string, amount *big.Int) (*types.Transaction, error) {
	const ethereumChain = "ethereum"
	if chain != ethereumChain {
		return nil, fmt.Errorf("signing transaction is only supported for Ethereum")
	}

	privateKey, err := w.deriveChildKey()
	if err != nil {
		return nil, err
	}

	nonce, err := w.ethereumClient.PendingNonceAt(context.Background(), common.HexToAddress(toAddress))
	if err != nil {
		return nil, err
	}

	gasPrice, err := w.ethereumClient.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}

	msg := ethereum.CallMsg{
		From:     crypto.PubkeyToAddress(privateKey.PublicKey),
		To:       &common.Address{},
		Value:    amount,
		GasPrice: gasPrice,
		Data:     nil,
	}

	err = rlp.Encode(nil, msg)
	if err != nil {
		return nil, err
	}

	gasLimit := uint64(21000)
	tx := types.NewTransaction(nonce, common.HexToAddress(toAddress), amount, gasLimit, gasPrice, nil)

	signedTx, err := types.SignTx(tx, types.HomesteadSigner{}, privateKey)
	if err != nil {
		return nil, err
	}

	return signedTx, nil
}
