package factory

import (
	common2 "tokit/wallet/pkg/common2"
)

// BlockchainFactory defines an interface for creating Blockchain instances.
type BlockchainFactory interface {
	Create(chain string) (common2.Blockchain, error)
}
