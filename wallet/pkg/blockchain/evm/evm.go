package evm

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"tokit/wallet/pkg/common2"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Client implements the common2.Blockchain interface.
type Client struct {
	rpcURL     string
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
	address    string
}

// NewClient creates a new EVM client.
func NewClient(rpcURL string) *Client {
	return &Client{rpcURL: rpcURL}
}

// GenerateKey generates a new private key and sets the address.
func (c *Client) GenerateKey() ([]byte, error) {
	privateKeyECDSA, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}
	c.privateKey = privateKeyECDSA
	c.publicKey = &privateKeyECDSA.PublicKey
	c.address = crypto.PubkeyToAddress(*c.publicKey).Hex()
	privateKeyBytes := crypto.FromECDSA(privateKeyECDSA)
	return privateKeyBytes, nil
}

// GetAddress returns the EVM address.
func (c *Client) GetAddress() (string, error) {
	return c.address, nil
}

// GetChainID returns the chain ID of the connected network.
func (c *Client) GetChainID() (int64, error) {
	ctx := context.Background()
	client, err := ethclient.DialContext(ctx, c.rpcURL)
	if err != nil {
		return 0, err
	}
	defer client.Close()

	chainID, err := client.ChainID(ctx)
	if err != nil {
		return 0, err
	}
	return chainID.Int64(), nil
}

// Dial connects to the Ethereum client.
func (c *Client) Dial(ctx context.Context) (*ethclient.Client, error) {
	client, err := ethclient.DialContext(ctx, c.rpcURL)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// GetRPCURL returns the RPC URL of the client.
func (c *Client) GetRPCURL() string {
	return c.rpcURL
}

// Transfer transfers funds to the specified address.
func (c *Client) Transfer(privateKeyHex string, toAddress string, amount float64) (string, error) {
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return "", fmt.Errorf("invalid private key: %w", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", fmt.Errorf("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	ctx := context.Background()
	client, err := c.Dial(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to connect to RPC: %w", err)
	}
	defer client.Close()

	nonce, err := client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		return "", fmt.Errorf("failed to get nonce: %w", err)
	}

	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get suggested gas price: %w", err)
	}

	// Convert amount to wei
	amountWei := new(big.Int)
	amountEth := big.NewFloat(amount)
	power := new(big.Int).SetInt64(18)
	ten := big.NewInt(10)
	powerOfTen := new(big.Int).Exp(ten, power, nil)
	floatValue := new(big.Float).Mul(amountEth, new(big.Float).SetInt(powerOfTen))
	amountWei, _ = floatValue.Int(amountWei)

	toAddressCommon := common.HexToAddress(toAddress)
	var data []byte
	tx := types.NewTransaction(nonce, toAddressCommon, amountWei, uint64(21000), gasPrice, data)

	chainID, err := client.ChainID(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get chain ID: %w", err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %w", err)
	}

	err = client.SendTransaction(ctx, signedTx)
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %w", err)
	}

	return signedTx.Hash().Hex(), nil
}

// GetTransactionHistory retrieves the transaction history for the given address.
func (c *Client) GetTransactionHistory(address string) ([]common2.Transaction, error) {
	// Placeholder implementation
	return []common2.Transaction{}, nil
}
