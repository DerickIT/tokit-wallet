package chain

import (
	"context"
	"fmt"
	"math/big"

	"tokit/internal/config"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Client struct {
	EthClient *ethclient.Client
	ChainID   *big.Int
	Config    config.NetworkConfig
}

// NewClient creates a new client for the specified chain
func NewClient(chainName string, cfg *config.Config) (*Client, error) {
	networkCfg, ok := cfg.Networks[chainName]
	if !ok {
		return nil, fmt.Errorf("network configuration not found for: %s", chainName)
	}

	client, err := ethclient.Dial(networkCfg.RPCURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %w", chainName, err)
	}

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID: %w", err)
	}

	// Verify chain ID matches config
	if chainID.Int64() != networkCfg.ChainID {
		return nil, fmt.Errorf("chain ID mismatch for %s: expected %d, got %d", chainName, networkCfg.ChainID, chainID)
	}

	return &Client{
		EthClient: client,
		ChainID:   chainID,
		Config:    networkCfg,
	}, nil
}

func (c *Client) GetBalance(address string) (*big.Int, error) {
	if !common.IsHexAddress(address) {
		return nil, fmt.Errorf("invalid address: %s", address)
	}
	account := common.HexToAddress(address)
	return c.EthClient.BalanceAt(context.Background(), account, nil)
}

func (c *Client) Close() {
	c.EthClient.Close()
}
