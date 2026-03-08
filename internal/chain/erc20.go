package chain

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type TokenMetadata struct {
	Symbol   string
	Decimals uint8
}

var erc20ABI = mustParseABI(`[
	{"constant":true,"inputs":[{"name":"owner","type":"address"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"stateMutability":"view","type":"function"},
	{"constant":true,"inputs":[],"name":"decimals","outputs":[{"name":"","type":"uint8"}],"stateMutability":"view","type":"function"},
	{"constant":true,"inputs":[],"name":"symbol","outputs":[{"name":"","type":"string"}],"stateMutability":"view","type":"function"},
	{"constant":false,"inputs":[{"name":"to","type":"address"},{"name":"value","type":"uint256"}],"name":"transfer","outputs":[{"name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"}
]`)

func mustParseABI(definition string) abi.ABI {
	parsed, err := abi.JSON(strings.NewReader(definition))
	if err != nil {
		panic(err)
	}
	return parsed
}

// GetTokenMetadata returns basic ERC20 metadata used for balances and transfers.
func (c *Client) GetTokenMetadata(tokenAddress string) (*TokenMetadata, error) {
	tokenAddr := common.HexToAddress(tokenAddress)

	symbolValues, err := c.callERC20(tokenAddr, "symbol")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch token symbol: %w", err)
	}
	decimalsValues, err := c.callERC20(tokenAddr, "decimals")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch token decimals: %w", err)
	}

	symbol, ok := symbolValues[0].(string)
	if !ok || symbol == "" {
		return nil, fmt.Errorf("token symbol returned an unexpected value")
	}
	decimals, ok := decimalsValues[0].(uint8)
	if !ok {
		return nil, fmt.Errorf("token decimals returned an unexpected value")
	}

	return &TokenMetadata{Symbol: symbol, Decimals: decimals}, nil
}

// GetTokenBalance returns the balance of an ERC20 token.
func (c *Client) GetTokenBalance(tokenAddress, ownerAddress string) (*big.Int, error) {
	tokenAddr := common.HexToAddress(tokenAddress)
	ownerAddr := common.HexToAddress(ownerAddress)

	values, err := c.callERC20(tokenAddr, "balanceOf", ownerAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to call contract: %w", err)
	}

	balance, ok := values[0].(*big.Int)
	if !ok || balance == nil {
		return nil, fmt.Errorf("token balance returned an unexpected value")
	}

	return balance, nil
}

// SendTokenTransaction sends an ERC20 token transfer.
func (c *Client) SendTokenTransaction(
	from accounts.Account,
	tokenAddress string,
	to string,
	amountBaseUnits *big.Int,
	signFn func(accounts.Account, *types.Transaction, *big.Int) (*types.Transaction, error),
) (string, error) {
	ctx := context.Background()
	tokenAddr := common.HexToAddress(tokenAddress)
	toAddr := common.HexToAddress(to)

	data, err := erc20ABI.Pack("transfer", toAddr, amountBaseUnits)
	if err != nil {
		return "", fmt.Errorf("failed to encode token transfer: %w", err)
	}

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

	gasLimit, err := c.EstimateGas(from.Address, tokenAddr, big.NewInt(0), data)
	if err != nil {
		gasLimit = 100000
	}

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

	signedTx, err := signFn(from, tx, c.ChainID)
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %w", err)
	}

	if err := c.EthClient.SendTransaction(ctx, signedTx); err != nil {
		return "", fmt.Errorf("failed to send transaction: %w", err)
	}

	return signedTx.Hash().Hex(), nil
}

func (c *Client) callERC20(token common.Address, method string, args ...interface{}) ([]interface{}, error) {
	data, err := erc20ABI.Pack(method, args...)
	if err != nil {
		return nil, err
	}

	msg := ethereum.CallMsg{To: &token, Data: data}
	result, err := c.EthClient.CallContract(context.Background(), msg, nil)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("no result from contract call")
	}

	values, err := erc20ABI.Unpack(method, result)
	if err != nil {
		return nil, err
	}
	if len(values) == 0 {
		return nil, fmt.Errorf("contract call returned no values")
	}

	return values, nil
}
