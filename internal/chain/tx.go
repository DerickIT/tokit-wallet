package chain

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// SendTransaction builds, signs, and sends an EIP-1559 transaction.
func (c *Client) SendTransaction(
	from accounts.Account,
	to string,
	amountWei *big.Int,
	signFn func(accounts.Account, *types.Transaction, *big.Int) (*types.Transaction, error),
) (string, error) {
	ctx := context.Background()

	nonce, err := c.EthClient.PendingNonceAt(ctx, from.Address)
	if err != nil {
		return "", fmt.Errorf("failed to get nonce: %w", err)
	}

	gasTipCap, err := c.EthClient.SuggestGasTipCap(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get gas tip cap: %w", err)
	}

	head, err := c.EthClient.HeaderByNumber(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("failed to get header: %w", err)
	}

	gasFeeCap := new(big.Int).Add(
		new(big.Int).Mul(head.BaseFee, big.NewInt(2)),
		gasTipCap,
	)

	toAddr := common.HexToAddress(to)
	gasLimit := uint64(21000)

	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   c.ChainID,
		Nonce:     nonce,
		GasTipCap: gasTipCap,
		GasFeeCap: gasFeeCap,
		Gas:       gasLimit,
		To:        &toAddr,
		Value:     amountWei,
		Data:      nil,
	})

	signedTx, err := signFn(from, tx, c.ChainID)
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %w", err)
	}

	if err := c.EthClient.SendTransaction(ctx, signedTx); err != nil {
		return "", fmt.Errorf("failed to send transaction: %w", err)
	}

	return signedTx.Hash().Hex(), nil
}
