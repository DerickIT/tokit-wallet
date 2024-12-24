package client

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math"
	"math/big"

	common2 "tokit/wallet/pkg/common2"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Client represents an Ethereum client.
type Client struct {
	rpcURL  string
	address string
}

// NewClient creates a new Ethereum client.
func NewClient(rpcURL string) *Client {
	return &Client{rpcURL: rpcURL}
}

// Dial connects to the Ethereum network.
func (c *Client) Dial(ctx context.Context) (*ethclient.Client, error) {
	ethClient, err := ethclient.DialContext(ctx, c.rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to ethereum network: %w", err)
	}
	return ethClient, nil
}

// GenerateKey generates a new Ethereum private key.
func (c *Client) GenerateKey() ([]byte, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %w", err)
	}
	c.address = crypto.PubkeyToAddress(privateKey.PublicKey).Hex()
	return crypto.FromECDSA(privateKey), nil
}

// GetAddress derives the Ethereum address from the public key.
func (c *Client) GetAddress() (string, error) {
	return c.address, nil
}

// SignTransaction signs the Ethereum transaction.
func (c *Client) SignTransaction(ctx context.Context, ethClient *ethclient.Client, privateKey *ecdsa.PrivateKey, toAddress string, amount float64) ([]byte, error) {
	fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
	nonce, err := ethClient.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get nonce: %w", err)
	}

	value := new(big.Int)
	value.SetString(fmt.Sprintf("%.0f", amount*math.Pow10(18)), 10)
	gasLimit := uint64(21000) // TODO: make gas limit dynamic
	gasPrice, err := ethClient.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get gas price: %w", err)
	}

	to := common.HexToAddress(toAddress)
	var data []byte
	tx := types.NewTransaction(nonce, to, value, gasLimit, gasPrice, data)

	chainID, err := c.GetChainID()
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID: %w", err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(chainID)), privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %w", err)
	}

	txData, err := signedTx.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transaction: %w", err)
	}
	return txData, nil
}

// GetChainID retrieves the chain ID of the connected Ethereum network.
func (c *Client) GetChainID() (int64, error) {
	ctx := context.Background()
	ethClient, err := ethclient.DialContext(ctx, c.rpcURL)
	if err != nil {
		return 0, fmt.Errorf("failed to connect to ethereum network: %w", err)
	}
	defer ethClient.Close()

	chainID, err := ethClient.ChainID(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get chain ID: %w", err)
	}
	return chainID.Int64(), nil
}

// GetTransactionHistory retrieves the transaction history for the given address.
func (c *Client) GetTransactionHistory(address string) ([]common2.Transaction, error) {
	fmt.Println("Getting transaction history for address:", address)
	return []common2.Transaction{}, nil
}

// Transfer transfers funds on the Ethereum network.
func (c *Client) Transfer(privateKey *ecdsa.PrivateKey, toAddress string, amount float64) (string, error) {
	ctx := context.Background()
	ethClient, err := c.Dial(ctx)
	if err != nil {
		return "", err
	}
	defer ethClient.Close()

	fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
	nonce, err := ethClient.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		return "", err
	}

	value := new(big.Int)
	value.SetString(fmt.Sprintf("%.0f", amount*math.Pow10(18)), 10)
	gasLimit := uint64(21000) // TODO: make gas limit dynamic
	gasPrice, err := ethClient.SuggestGasPrice(ctx)
	if err != nil {
		return "", err
	}

	to := common.HexToAddress(toAddress)
	var data []byte
	tx := types.NewTransaction(nonce, to, value, gasLimit, gasPrice, data)

	chainID, err := c.GetChainID()
	if err != nil {
		return "", fmt.Errorf("failed to get chain ID: %w", err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(chainID)), privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %w", err)
	}

	err = ethClient.SendTransaction(ctx, signedTx)
	if err != nil {
		return "", err
	}

	return signedTx.Hash().Hex(), nil
}
