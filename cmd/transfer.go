package cmd

import (
	"fmt"
	"math/big"
	"strings"
	"syscall"

	"tokit/internal/chain"
	"tokit/internal/utils"
	"tokit/internal/wallet"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var transferCmd = &cobra.Command{
	Use:   "transfer [chain] [to] [amount]",
	Short: "Transfer funds to another address",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		chainName := args[0]
		toAddress := args[1]
		amountStr := args[2]

		// Parse amount
		amount := new(big.Float)
		_, ok := amount.SetString(amountStr)
		if !ok {
			utils.Log.Fatal("Invalid amount")
		}

		// Init Wallet Service
		svc, err := wallet.NewService()
		if err != nil {
			utils.Log.Fatalf("Failed to init wallet service: %v", err)
		}

		// Get Sender Account (First one for now, or select via flag later)
		accountsList := svc.ListAccounts()
		if len(accountsList) == 0 {
			utils.Log.Fatal("No accounts found. Please create or import a wallet.")
		}
		fromAccount := accountsList[0]

		// Init Chain Client
		client, err := chain.NewClient(chainName, AppConfig)
		if err != nil {
			utils.Log.Fatalf("Failed to create client: %v", err)
		}
		defer client.Close()

		// Confirm Transaction
		fmt.Printf("\n⚠️  CONFIRM TRANSACTION\n")
		fmt.Printf("Chain:  %s\n", chainName)
		fmt.Printf("From:   %s\n", fromAccount.Address.Hex())
		fmt.Printf("To:     %s\n", toAddress)
		fmt.Printf("Amount: %s %s\n", amountStr, client.Config.Symbol)
		fmt.Println(strings.Repeat("-", 40))

		fmt.Print("Enter password to confirm: ")
		bytePassword, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			utils.Log.Fatalf("Failed to read password: %v", err)
		}
		password := string(bytePassword)
		fmt.Println()

		// Define Signer Function
		signFn := func(a accounts.Account, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error) {
			return svc.SignTx(a, tx, chainID, password)
		}

		// Send Transaction
		fmt.Println("\nSending transaction...")
		txHash, err := client.SendTransaction(fromAccount, toAddress, amount, signFn)
		if err != nil {
			utils.Log.Fatalf("Failed to send transaction: %v", err)
		}

		fmt.Printf("\n✅ Transaction Sent!\nHash: %s\n", txHash)
		fmt.Printf("Explorer: %s/tx/%s\n", client.Config.Explorer, txHash)
	},
}

func init() {
	rootCmd.AddCommand(transferCmd)
}
