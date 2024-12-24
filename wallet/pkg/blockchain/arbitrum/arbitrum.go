package arbitrum

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math"
	"math/big"

	"tokit/wallet/pkg/blockchain/common"
	common2 "tokit/wallet/pkg/common2"

	"github.com/ethereum/go-ethereum"
	ethereumcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Arbitrum implements the common.EVMClient interface for Arbitrum.
type Arbitrum struct {
	evmClient  common.EVMClient
	address    string
	privateKey *ecdsa.PrivateKey
}

// NewArbitrum creates a new Arbitrum instance.
func NewArbitrum(rpcURL string) *Arbitrum {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		// Handle error appropriately
		panic(err)
	}
	return &Arbitrum{evmClient: &common.EVMClientImpl{Client: client, RPCURL: rpcURL}}
}

// GetRPCURL returns the RPC URL of the client.
func (a *Arbitrum) GetRPCURL() string {
	return a.evmClient.GetRPCURL()
}

// GenerateKey generates a new Arbitrum private key and sets the address.
func (a *Arbitrum) GenerateKey() ([]byte, error) {
	key, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}
	a.privateKey = key
	a.address = crypto.PubkeyToAddress(key.PublicKey).Hex()
	return crypto.FromECDSA(key), nil
}

// GetAddress returns the Arbitrum address.
func (a *Arbitrum) GetAddress() (string, error) {
	return a.address, nil
}

// GetChainID returns the chain ID of the connected network.
func (a *Arbitrum) GetChainID() (int64, error) {
	return a.evmClient.GetChainID()
}

// Dial connects to the Arbitrum client.
func (a *Arbitrum) Dial(ctx context.Context) (*ethclient.Client, error) {
	rpcURL := a.evmClient.GetRPCURL()
	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// Transfer transfers funds on the Arbitrum network.
func (a *Arbitrum) Transfer(privateKeyHex string, toAddress string, amount float64) (string, error) {
	ctx := context.Background()
	client, err := a.evmClient.Dial(ctx)
	if err != nil {
		return "", err
	}
	defer client.Close()

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return "", err
	}
	a.privateKey = privateKey
	fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)

	nonce, err := client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		return "", err
	}

	value := new(big.Int)
	value.SetString(fmt.Sprintf("%.0f", amount*math.Pow10(18)), 10)

	to := ethereumcommon.HexToAddress(toAddress)
	var data []byte

	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return "", err
	}

	gasLimit, err := client.EstimateGas(ctx, ethereum.CallMsg{
		From:     fromAddress,
		To:       &to,
		GasPrice: gasPrice,
		Value:    value,
		Data:     data,
	})
	if err != nil {
		return "", err
	}

	tx := types.NewTransaction(nonce, to, value, gasLimit, gasPrice, data)

	chainID, err := a.evmClient.GetChainID()
	if err != nil {
		return "", fmt.Errorf("failed to get arbitrum chain ID: %w", err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(chainID)), a.privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %w", err)
	}

	err = client.SendTransaction(ctx, signedTx)
	if err != nil {
		return "", err
	}

	return signedTx.Hash().Hex(), nil
}

// GetTransactionHistory retrieves the transaction history for the given Arbitrum address.
func (a *Arbitrum) GetTransactionHistory(address string) ([]common2.Transaction, error) {
	fmt.Println("Getting Arbitrum transaction history for address:", address)
	return []common2.Transaction{}, nil
}
