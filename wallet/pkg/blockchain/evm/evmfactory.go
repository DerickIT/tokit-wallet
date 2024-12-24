package evm

import (
	"fmt"
	"tokit/wallet/pkg/blockchain/arbitrum"
	"tokit/wallet/pkg/blockchain/base"
	"tokit/wallet/pkg/blockchain/ethereum"
	"tokit/wallet/pkg/blockchain/optimism"
	"tokit/wallet/pkg/common2"
)

// Create creates a new blockchain client based on the given chain name.
func Create(chainName string, rpcURL string) (common2.Blockchain, error) {
	switch chainName {
	case "ethereum":
		return ethereum.NewEthereum(rpcURL), nil
	case "arbitrum":
		return arbitrum.NewArbitrum(rpcURL), nil
	case "optimism":
		return optimism.NewOptimism(rpcURL), nil
	case "base":
		return base.NewBase(rpcURL), nil
	default:
		return nil, fmt.Errorf("unsupported chain: %s", chainName)
	}
}
