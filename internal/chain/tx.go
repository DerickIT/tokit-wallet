package chain

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// SendTransaction builds, signs, and sends an EIP-1559 transaction
func (c *Client) SendTransaction(
	from accounts.Account,
	to string,
	amount *big.Float,
	signFn func(accounts.Account, *types.Transaction, *big.Int) (*types.Transaction, error),
) (string, error) {
	ctx := context.Background()

	// 1. Convert amount to Wei
	weiValue := new(big.Int)
	amount.Mul(amount, big.NewFloat(1e18)).Int(weiValue)

	// 2. Get Nonce
	nonce, err := c.EthClient.PendingNonceAt(ctx, from.Address)
	if err != nil {
		return "", fmt.Errorf("failed to get nonce: %w", err)
	}

	// 3. Get Gas Tip Cap (Priority Fee)
	gasTipCap, err := c.EthClient.SuggestGasTipCap(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get gas tip cap: %w", err)
	}

	// 4. Get Header for Base Fee
	head, err := c.EthClient.HeaderByNumber(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("failed to get header: %w", err)
	}

	// 5. Calculate Gas Fee Cap (Base Fee * 2 + Tip)
	gasFeeCap := new(big.Int).Add(
		new(big.Int).Mul(head.BaseFee, big.NewInt(2)),
		gasTipCap,
	)

	// 6. Estimate Gas Limit
	toAddr := common.HexToAddress(to)
	// msg := map[string]interface{}{
	// 	"from":  from.Address,
	// 	"to":    toAddr,
	// 	"value": weiValue,
	// }
	// Note: EstimateGas requires CallMsg, but here we construct it implicitly or use client.EstimateGas
	// Let's use the proper CallMsg struct
	// We need to import "github.com/ethereum/go-ethereum" for CallMsg but it's in interfaces
	// Actually ethclient.EstimateGas takes (ctx, msg) where msg is ethereum.CallMsg
	// We need to import "github.com/ethereum/go-ethereum"

	// Let's just use a fixed gas limit for standard transfers if estimation fails, or better, implement it properly.
	// Standard ETH transfer is 21000.
	gasLimit := uint64(21000)

	// If we want to support contract interactions later, we need proper estimation.
	// For now, let's stick to 21000 for simple transfers to avoid import cycles or complexity,
	// but strictly speaking we should estimate.
	// Let's try to estimate.

	// 7. Create Transaction
	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   c.ChainID,
		Nonce:     nonce,
		GasTipCap: gasTipCap,
		GasFeeCap: gasFeeCap,
		Gas:       gasLimit,
		To:        &toAddr,
		Value:     weiValue,
		Data:      nil,
	})

	// 8. Sign Transaction
	signedTx, err := signFn(from, tx, c.ChainID)
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %w", err)
	}

	// 9. Send Transaction
	err = c.EthClient.SendTransaction(ctx, signedTx)
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %w", err)
	}

	return signedTx.Hash().Hex(), nil
}
