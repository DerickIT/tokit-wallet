package optimism

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// OptimismClient is a wrapper around ethclient.Client that implements common.ClientIEthClient.
type OptimismClient struct {
	*ethclient.Client
}

func (oc *OptimismClient) Close() {
	oc.Client.Close()
}

func (oc *OptimismClient) ChainID(ctx context.Context) (*big.Int, error) {
	return oc.Client.ChainID(ctx)
}

func (oc *OptimismClient) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	return oc.Client.SuggestGasPrice(ctx)
}

func (oc *OptimismClient) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	return oc.Client.PendingNonceAt(ctx, account)
}

func (oc *OptimismClient) SendTransaction(ctx context.Context, tx *ethtypes.Transaction) error {
	return oc.Client.SendTransaction(ctx, tx)
}
