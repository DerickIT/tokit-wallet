package chain

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// PrepareNativeTransfer builds an unsigned EIP-1559 transfer transaction.
func (c *Client) PrepareNativeTransfer(from common.Address, to string, amountWei *big.Int) (*types.Transaction, error) {
	toAddr := common.HexToAddress(to)
	return c.prepareDynamicFeeTx(from, toAddr, amountWei, nil, 21000)
}

// SendPreparedTransaction signs and broadcasts a prepared transaction.
func (c *Client) SendPreparedTransaction(
	from accounts.Account,
	tx *types.Transaction,
	signFn func(accounts.Account, *types.Transaction, *big.Int) (*types.Transaction, error),
) (string, error) {
	signedTx, err := signFn(from, tx, c.ChainID)
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %w", err)
	}

	if err := c.EthClient.SendTransaction(context.Background(), signedTx); err != nil {
		return "", fmt.Errorf("failed to send transaction: %w", err)
	}

	return signedTx.Hash().Hex(), nil
}

func (c *Client) prepareDynamicFeeTx(from common.Address, to common.Address, value *big.Int, data []byte, fallbackGas uint64) (*types.Transaction, error) {
	ctx := context.Background()

	nonce, err := c.EthClient.PendingNonceAt(ctx, from)
	if err != nil {
		return nil, fmt.Errorf("failed to get nonce: %w", err)
	}

	gasTipCap, gasFeeCap, err := c.suggestDynamicFees(ctx)
	if err != nil {
		return nil, err
	}

	gasLimit, err := c.EstimateGas(from, to, value, data)
	if err != nil {
		gasLimit = fallbackGas
	}

	return types.NewTx(&types.DynamicFeeTx{
		ChainID:   c.ChainID,
		Nonce:     nonce,
		GasTipCap: gasTipCap,
		GasFeeCap: gasFeeCap,
		Gas:       gasLimit,
		To:        &to,
		Value:     value,
		Data:      data,
	}), nil
}

func (c *Client) suggestDynamicFees(ctx context.Context) (*big.Int, *big.Int, error) {
	gasTipCap, err := c.EthClient.SuggestGasTipCap(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get gas tip cap: %w", err)
	}

	head, err := c.EthClient.HeaderByNumber(ctx, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get header: %w", err)
	}
	if head.BaseFee == nil {
		return nil, nil, fmt.Errorf("latest block did not include a base fee")
	}

	gasFeeCap := new(big.Int).Add(
		new(big.Int).Mul(head.BaseFee, big.NewInt(2)),
		gasTipCap,
	)

	return gasTipCap, gasFeeCap, nil
}
