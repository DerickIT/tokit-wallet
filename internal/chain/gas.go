package chain

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
)

// EstimateGas calculates the gas limit for a transaction
func (c *Client) EstimateGas(from, to common.Address, value *big.Int, data []byte) (uint64, error) {
	msg := ethereum.CallMsg{
		From:     from,
		To:       &to,
		Gas:      0,
		GasPrice: nil,
		Value:    value,
		Data:     data,
	}

	gasLimit, err := c.EthClient.EstimateGas(context.Background(), msg)
	if err != nil {
		return 0, fmt.Errorf("failed to estimate gas: %w", err)
	}

	// Add a buffer (e.g., 10%) to be safe
	return gasLimit + (gasLimit / 10), nil
}
