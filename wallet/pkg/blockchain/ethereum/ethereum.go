package ethereum

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

// Ethereum implements the common.EVMClient interface for Ethereum.
type Ethereum struct {
	evmClient  common.EVMClient
	address    string
	privateKey *ecdsa.PrivateKey
}

// NewEthereum creates a new Ethereum instance.
func NewEthereum(rpcURL string) *Ethereum {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		// Handle error appropriately
		panic(err)
	}
	return &Ethereum{evmClient: &common.EVMClientImpl{Client: client, RPCURL: rpcURL}}
}

// GetRPCURL returns the RPC URL of the client.
func (e *Ethereum) GetRPCURL() string {
	return e.evmClient.GetRPCURL()
}

// GenerateKey generates a new Ethereum private key and sets the address.
func (e *Ethereum) GenerateKey() ([]byte, error) {
	key, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}
	e.privateKey = key
	e.address = crypto.PubkeyToAddress(key.PublicKey).Hex()
	return crypto.FromECDSA(key), nil
}

// GetAddress returns the Ethereum address.
func (e *Ethereum) GetAddress() (string, error) {
	return e.address, nil
}

// GetChainID returns the chain ID of the connected network.
func (e *Ethereum) GetChainID() (int64, error) {
	return e.evmClient.GetChainID()
}

// Dial connects to the Ethereum client.
func (e *Ethereum) Dial(ctx context.Context) (*ethclient.Client, error) {
	rpcURL := e.evmClient.GetRPCURL()
	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// Transfer transfers funds on the Ethereum network.
func (e *Ethereum) Transfer(privateKeyHex string, toAddress string, amount float64) (string, error) {
	ctx := context.Background()
	client, err := e.evmClient.Dial(ctx)
	if err != nil {
		return "", err
	}
	defer client.Close()

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return "", err
	}
	e.privateKey = privateKey
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

	chainID, err := e.evmClient.GetChainID()
	if err != nil {
		return "", fmt.Errorf("failed to get ethereum chain ID: %w", err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(chainID)), e.privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %w", err)
	}

	err = client.SendTransaction(ctx, signedTx)
	if err != nil {
		return "", err
	}

	return signedTx.Hash().Hex(), nil
}

// GetTransactionHistory retrieves the transaction history for the given Ethereum address.
func (e *Ethereum) GetTransactionHistory(address string) ([]common2.Transaction, error) {
	fmt.Println("Getting Ethereum transaction history for address:", address)
	return []common2.Transaction{}, nil
}
