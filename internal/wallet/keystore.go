package wallet

import (
	"fmt"
	"math/big"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type Service struct {
	ks *keystore.KeyStore
}

func NewService() (*Service, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	keystorePath := filepath.Join(home, ".tokit", "keystore")
	if err := os.MkdirAll(keystorePath, 0700); err != nil {
		return nil, err
	}

	// Use StandardScryptN and StandardScryptP for better security
	ks := keystore.NewKeyStore(keystorePath, keystore.StandardScryptN, keystore.StandardScryptP)

	return &Service{ks: ks}, nil
}

func (s *Service) CreateAccount(password string) (accounts.Account, error) {
	return s.ks.NewAccount(password)
}

func (s *Service) ImportAccount(keyJSON []byte, password, newPassword string) (accounts.Account, error) {
	return s.ks.Import(keyJSON, password, newPassword)
}

func (s *Service) ListAccounts() []accounts.Account {
	return s.ks.Accounts()
}

func (s *Service) SignTx(account accounts.Account, tx *types.Transaction, chainID *big.Int, password string) (*types.Transaction, error) {
	if err := s.ks.Unlock(account, password); err != nil {
		return nil, fmt.Errorf("failed to unlock account: %w", err)
	}
	// Ensure we lock it back even if signing fails, though SignTx usually doesn't need unlock if we use the keystore's SignTx method directly?
	// Actually keystore.SignTx requires unlock first or we can use SignTxWithPassphrase

	// Wait, go-ethereum's keystore.SignTxWithPassphrase is the easiest way
	return s.ks.SignTxWithPassphrase(account, password, tx, chainID)
}

func (s *Service) GetAccount(addressHex string) (accounts.Account, error) {
	if !common.IsHexAddress(addressHex) {
		return accounts.Account{}, fmt.Errorf("invalid address: %s", addressHex)
	}
	addr := common.HexToAddress(addressHex)

	for _, acc := range s.ks.Accounts() {
		if acc.Address == addr {
			return acc, nil
		}
	}
	return accounts.Account{}, fmt.Errorf("account not found: %s", addressHex)
}
