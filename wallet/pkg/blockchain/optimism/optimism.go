package optimism

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

// Optimism implements the common.EVMClient interface for Optimism.
type Optimism struct {
	evmClient  common.EVMClient
	address    string
	privateKey *ecdsa.PrivateKey
}

// NewOptimism creates a new Optimism instance.
func NewOptimism(rpcURL string) *Optimism {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		// Handle error appropriately, maybe return nil and an error
		panic(err)
	}
	return &Optimism{evmClient: &common.EVMClientImpl{Client: client, RPCURL: rpcURL}}
}

// GetRPCURL returns the RPC URL of the client.
func (o *Optimism) GetRPCURL() string {
	return o.evmClient.GetRPCURL()
}

// GenerateKey generates a new Optimism private key and sets the address.
func (o *Optimism) GenerateKey() ([]byte, error) {
	key, err := o.evmClient.GenerateKey()
	if err != nil {
		return nil, err
	}
	privateKey, err := crypto.ToECDSA(key)
	if err != nil {
		return nil, err
	}
	o.privateKey = privateKey
	o.address = crypto.PubkeyToAddress(privateKey.PublicKey).Hex()
	return key, nil
}

// GetAddress returns the Optimism address.
func (o *Optimism) GetAddress() (string, error) {
	return o.address, nil
}

// GetChainID returns the chain ID of the connected network.
func (o *Optimism) GetChainID() (int64, error) {
	return o.evmClient.GetChainID()
}

// Dial connects to the Optimism client.
func (o *Optimism) Dial(ctx context.Context) (*ethclient.Client, error) {
	rpcURL := o.evmClient.GetRPCURL()
	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// Transfer transfers funds on the Optimism network.
func (o *Optimism) Transfer(privateKeyHex string, toAddress string, amount float64) (string, error) {
	ctx := context.Background()
	client, err := o.evmClient.Dial(ctx)
	if err != nil {
		return "", err
	}
	defer client.Close()

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return "", err
	}
	o.privateKey = privateKey
	fromAddress := crypto.PubkeyToAddress(o.privateKey.PublicKey)

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

	chainID, err := o.evmClient.GetChainID()
	if err != nil {
		return "", fmt.Errorf("failed to get optimism chain ID: %w", err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(chainID)), o.privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %w", err)
	}

	err = client.SendTransaction(ctx, signedTx)
	if err != nil {
		return "", err
	}

	return signedTx.Hash().Hex(), nil
}

// GetTransactionHistory retrieves the transaction history for the given Optimism address.
func (o *Optimism) GetTransactionHistory(address string) ([]common2.Transaction, error) {
	fmt.Println("Getting Optimism transaction history for address:", address)
	return []common2.Transaction{}, nil
}
