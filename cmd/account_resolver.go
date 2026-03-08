package cmd

import (
	"fmt"

	"tokit/internal/wallet"

	"github.com/ethereum/go-ethereum/accounts"
)

func resolveLocalAccount(svc *wallet.Service, preferredAddress string) (accounts.Account, error) {
	if preferredAddress != "" {
		return svc.GetAccount(preferredAddress)
	}

	if AppConfig != nil && AppConfig.DefaultAccount != "" {
		return svc.GetAccount(AppConfig.DefaultAccount)
	}

	accountsList := svc.ListAccounts()
	if len(accountsList) == 0 {
		return accounts.Account{}, fmt.Errorf("no local accounts found")
	}

	return accountsList[0], nil
}
