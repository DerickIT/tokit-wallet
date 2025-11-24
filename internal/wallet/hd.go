package wallet

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

// GenerateMnemonic generates a new 12-word BIP39 mnemonic
func GenerateMnemonic() (string, error) {
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		return "", err
	}
	return bip39.NewMnemonic(entropy)
}

// ImportMnemonic derives a private key from a mnemonic and imports it into the keystore
// Uses BIP44 path: m/44'/60'/0'/0/0 (Standard Ethereum path for first account)
func (s *Service) ImportMnemonic(mnemonic, password string) (accounts.Account, error) {
	if !bip39.IsMnemonicValid(mnemonic) {
		return accounts.Account{}, errors.New("invalid mnemonic")
	}

	seed := bip39.NewSeed(mnemonic, "")
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return accounts.Account{}, err
	}

	// m/44'/60'/0'/0/0
	// 44'
	purpose, err := masterKey.NewChildKey(bip32.FirstHardenedChild + 44)
	if err != nil {
		return accounts.Account{}, err
	}
	// 60' (ETH)
	coinType, err := purpose.NewChildKey(bip32.FirstHardenedChild + 60)
	if err != nil {
		return accounts.Account{}, err
	}
	// 0' (Account)
	accountKey, err := coinType.NewChildKey(bip32.FirstHardenedChild + 0)
	if err != nil {
		return accounts.Account{}, err
	}
	// 0 (Change)
	change, err := accountKey.NewChildKey(0)
	if err != nil {
		return accounts.Account{}, err
	}
	// 0 (Index)
	addressKey, err := change.NewChildKey(0)
	if err != nil {
		return accounts.Account{}, err
	}

	privateKey, err := crypto.ToECDSA(addressKey.Key)
	if err != nil {
		return accounts.Account{}, err
	}

	return s.ks.ImportECDSA(privateKey, password)
}

// ImportPrivateKey imports a raw private key hex string
func (s *Service) ImportPrivateKey(privateKeyHex, password string) (accounts.Account, error) {
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return accounts.Account{}, fmt.Errorf("invalid private key: %w", err)
	}
	return s.ks.ImportECDSA(privateKey, password)
}
