package chain

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

// ERC20 ABI Method IDs
var (
	transferMethodID  = crypto.Keccak256([]byte("transfer(address,uint256)"))[:4]
	balanceOfMethodID = crypto.Keccak256([]byte("balanceOf(address)"))[:4]
)

// GetTokenBalance returns the balance of an ERC20 token
func (c *Client) GetTokenBalance(tokenAddress, ownerAddress string) (*big.Int, error) {
	tokenAddr := common.HexToAddress(tokenAddress)
	ownerAddr := common.HexToAddress(ownerAddress)

	// Construct data: balanceOf(address)
	// methodID (4 bytes) + padded address (32 bytes)
	data := make([]byte, 0)
	data = append(data, balanceOfMethodID...)
	data = append(data, common.LeftPadBytes(ownerAddr.Bytes(), 32)...)

	msg := ethereum.CallMsg{
		To:   &tokenAddr,
		Data: data,
	}

	result, err := c.EthClient.CallContract(context.Background(), msg, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to call contract: %w", err)
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no result from contract call")
	}

	return new(big.Int).SetBytes(result), nil
}

// SendTokenTransaction sends an ERC20 token transfer
func (c *Client) SendTokenTransaction(
	from accounts.Account,
	tokenAddress string,
	to string,
	amount *big.Float,
	signFn func(accounts.Account, *types.Transaction, *big.Int) (*types.Transaction, error),
) (string, error) {
	ctx := context.Background()
	tokenAddr := common.HexToAddress(tokenAddress)
	toAddr := common.HexToAddress(to)

	// 1. Convert amount to Wei (assuming 18 decimals for now, ideally should fetch decimals)
	// TODO: Fetch decimals dynamically if needed, but standard is 18
	weiValue := new(big.Int)
	amount.Mul(amount, big.NewFloat(1e18)).Int(weiValue)

	// 2. Construct Data: transfer(address,uint256)
	data := make([]byte, 0)
	data = append(data, transferMethodID...)
	data = append(data, common.LeftPadBytes(toAddr.Bytes(), 32)...)
	data = append(data, common.LeftPadBytes(weiValue.Bytes(), 32)...)

	// 3. Get Nonce
	nonce, err := c.EthClient.PendingNonceAt(ctx, from.Address)
	if err != nil {
		return "", fmt.Errorf("failed to get nonce: %w", err)
	}

	// 4. Get Gas Tip Cap
	gasTipCap, err := c.EthClient.SuggestGasTipCap(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get gas tip cap: %w", err)
	}

	// 5. Get Header for Base Fee
	head, err := c.EthClient.HeaderByNumber(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("failed to get header: %w", err)
	}

	// 6. Calculate Gas Fee Cap
	gasFeeCap := new(big.Int).Add(
		new(big.Int).Mul(head.BaseFee, big.NewInt(2)),
		gasTipCap,
	)

	// 7. Estimate Gas
	gasLimit, err := c.EstimateGas(from.Address, tokenAddr, big.NewInt(0), data)
	if err != nil {
		// Fallback if estimation fails (though it shouldn't for standard ERC20)
		gasLimit = 100000 // Standard ERC20 transfer is usually ~65k
	}

	// 8. Create Transaction
	// Note: 'To' is the Token Address, 'Value' is 0 (ETH), 'Data' contains the transfer details
	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   c.ChainID,
		Nonce:     nonce,
		GasTipCap: gasTipCap,
		GasFeeCap: gasFeeCap,
		Gas:       gasLimit,
		To:        &tokenAddr,
		Value:     big.NewInt(0),
		Data:      data,
	})

	// 9. Sign Transaction
	signedTx, err := signFn(from, tx, c.ChainID)
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %w", err)
	}

	// 10. Send Transaction
	err = c.EthClient.SendTransaction(ctx, signedTx)
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %w", err)
	}

	return signedTx.Hash().Hex(), nil
}
