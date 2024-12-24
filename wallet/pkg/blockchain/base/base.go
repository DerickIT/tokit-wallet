package base

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

// Base implements the common.EVMClient interface for Base.
type Base struct {
	evmClient  common.EVMClient
	address    string
	privateKey *ecdsa.PrivateKey
}

// NewBase creates a new Base instance.
func NewBase(rpcURL string) *Base {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		// Handle error appropriately
		panic(err)
	}
	return &Base{evmClient: &common.EVMClientImpl{Client: client, RPCURL: rpcURL}}
}

// GetRPCURL returns the RPC URL of the client.
func (b *Base) GetRPCURL() string {
	return b.evmClient.GetRPCURL()
}

// GenerateKey generates a new Base private key and sets the address.
func (b *Base) GenerateKey() ([]byte, error) {
	key, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}
	b.privateKey = key
	b.address = crypto.PubkeyToAddress(key.PublicKey).Hex()
	return crypto.FromECDSA(key), nil
}

// GetAddress returns the Base address.
func (b *Base) GetAddress() (string, error) {
	return b.address, nil
}

// GetChainID returns the chain ID of the connected network.
func (b *Base) GetChainID() (int64, error) {
	return b.evmClient.GetChainID()
}

// Dial connects to the Base client.
func (b *Base) Dial(ctx context.Context) (*ethclient.Client, error) {
	rpcURL := b.evmClient.GetRPCURL()
	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// Transfer transfers funds on the Base network.
func (b *Base) Transfer(privateKeyHex string, toAddress string, amount float64) (string, error) {
	ctx := context.Background()
	client, err := b.evmClient.Dial(ctx)
	if err != nil {
		return "", err
	}
	defer client.Close()

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return "", err
	}
	b.privateKey = privateKey
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

	chainID, err := b.evmClient.GetChainID()
	if err != nil {
		return "", fmt.Errorf("failed to get base chain ID: %w", err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(chainID)), b.privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %w", err)
	}

	err = client.SendTransaction(ctx, signedTx)
	if err != nil {
		return "", err
	}

	return signedTx.Hash().Hex(), nil
}

// GetTransactionHistory retrieves the transaction history for the given Base address.
func (b *Base) GetTransactionHistory(address string) ([]common2.Transaction, error) {
	fmt.Println("Getting Base transaction history for address:", address)
	return []common2.Transaction{}, nil
}
