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

var transferTokenAddress string
var transferFromAddress string

var transferCmd = &cobra.Command{
	Use:   "transfer [chain] [to] [amount]",
	Short: "Transfer funds to another address",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		chainName := args[0]
		toAddress := args[1]
		amountStr := args[2]

		svc, err := wallet.NewService()
		if err != nil {
			utils.Log.Fatalf("Failed to init wallet service: %v", err)
		}

		fromAccount, err := resolveLocalAccount(svc, transferFromAddress)
		if err != nil {
			utils.Log.Fatal("No accounts found. Please create or import a wallet.")
		}

		client, err := chain.NewClient(chainName, AppConfig)
		if err != nil {
			utils.Log.Fatalf("Failed to create client: %v", err)
		}
		defer client.Close()

		symbol := client.Config.Symbol
		amountBaseUnits := new(big.Int)
		unsignedTx := &types.Transaction{}

		if transferTokenAddress != "" {
			metadata, err := client.GetTokenMetadata(transferTokenAddress)
			if err != nil {
				utils.Log.Fatalf("Failed to get token metadata: %v", err)
			}
			parsedAmount, err := chain.ParseUnits(amountStr, metadata.Decimals)
			if err != nil {
				utils.Log.Fatalf("Invalid token amount: %v", err)
			}
			amountBaseUnits = parsedAmount
			symbol = metadata.Symbol
			unsignedTx, err = client.PrepareTokenTransfer(fromAccount.Address, transferTokenAddress, toAddress, amountBaseUnits)
			if err != nil {
				utils.Log.Fatalf("Failed to prepare token transfer: %v", err)
			}
		} else {
			parsedAmount, err := chain.ParseUnits(amountStr, 18)
			if err != nil {
				utils.Log.Fatalf("Invalid amount: %v", err)
			}
			amountBaseUnits = parsedAmount
			unsignedTx, err = client.PrepareNativeTransfer(fromAccount.Address, toAddress, amountBaseUnits)
			if err != nil {
				utils.Log.Fatalf("Failed to prepare transfer: %v", err)
			}
		}

		maxNetworkFee := estimateMaxNetworkFee(unsignedTx)
		fmt.Printf("\nCONFIRM TRANSACTION\n")
		fmt.Printf("Chain:            %s\n", chainName)
		fmt.Printf("From:             %s\n", fromAccount.Address.Hex())
		fmt.Printf("To:               %s\n", toAddress)
		fmt.Printf("Amount:           %s %s\n", amountStr, symbol)
		if transferTokenAddress != "" {
			fmt.Printf("Token:            %s\n", transferTokenAddress)
		}
		fmt.Printf("Gas limit:        %d\n", unsignedTx.Gas())
		fmt.Printf("Max fee per gas:  %s Gwei\n", chain.FormatUnits(unsignedTx.GasFeeCap(), 9, 4))
		fmt.Printf("Priority fee:     %s Gwei\n", chain.FormatUnits(unsignedTx.GasTipCap(), 9, 4))
		fmt.Printf("Max network fee:  %s %s\n", chain.FormatUnits(maxNetworkFee, 18, 6), client.Config.Symbol)
		if transferTokenAddress == "" {
			totalMaxCost := new(big.Int).Add(new(big.Int).Set(unsignedTx.Value()), maxNetworkFee)
			fmt.Printf("Max total spend:  %s %s\n", chain.FormatUnits(totalMaxCost, 18, 6), client.Config.Symbol)
		}
		fmt.Println(strings.Repeat("-", 48))

		fmt.Print("Enter password to confirm: ")
		bytePassword, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			utils.Log.Fatalf("Failed to read password: %v", err)
		}
		password := string(bytePassword)
		fmt.Println()

		signFn := func(a accounts.Account, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error) {
			return svc.SignTx(a, tx, chainID, password)
		}

		fmt.Println("\nSending transaction...")
		txHash, err := client.SendPreparedTransaction(fromAccount, unsignedTx, signFn)
		if err != nil {
			utils.Log.Fatalf("Failed to send transaction: %v", err)
		}

		fmt.Printf("\nTransaction sent.\nHash: %s\n", txHash)
		fmt.Printf("Explorer: %s/tx/%s\n", client.Config.Explorer, txHash)
	},
}

func estimateMaxNetworkFee(tx *types.Transaction) *big.Int {
	return new(big.Int).Mul(new(big.Int).SetUint64(tx.Gas()), tx.GasFeeCap())
}

func init() {
	rootCmd.AddCommand(transferCmd)
	transferCmd.Flags().StringVarP(&transferTokenAddress, "token", "t", "", "ERC20 token address")
	transferCmd.Flags().StringVar(&transferFromAddress, "from", "", "Local account address to send from")
}
