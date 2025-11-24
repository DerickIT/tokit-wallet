package cmd

import (
	"fmt"
	"math/big"
	"os"

	"tokit/internal/chain"
	"tokit/internal/utils"
	"tokit/internal/wallet"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var balanceCmd = &cobra.Command{
	Use:   "balance [chain] [address]",
	Short: "Check account balance",
	Long:  `Check the balance of an account on a specific blockchain. If address is omitted, checks the first local wallet account.`,
	Args:  cobra.RangeArgs(0, 2),
	Run: func(cmd *cobra.Command, args []string) {
		chainName := AppConfig.Default
		if len(args) > 0 {
			chainName = args[0]
		}

		var address string
		if len(args) > 1 {
			address = args[1]
		} else {
			// Get first account from local wallet
			svc, err := wallet.NewService()
			if err != nil {
				utils.Log.Fatalf("Failed to init wallet service: %v", err)
			}
			accounts := svc.ListAccounts()
			if len(accounts) == 0 {
				utils.Log.Fatal("No local accounts found. Please provide an address or create a wallet.")
			}
			address = accounts[0].Address.Hex()
		}

		client, err := chain.NewClient(chainName, AppConfig)
		if err != nil {
			utils.Log.Fatalf("Failed to create client: %v", err)
		}
		defer client.Close()

		balance, err := client.GetBalance(address)
		if err != nil {
			utils.Log.Fatalf("Failed to get balance: %v", err)
		}

		// Convert Wei to Ether
		fBalance := new(big.Float)
		fBalance.SetString(balance.String())
		ethValue := new(big.Float).Quo(fBalance, big.NewFloat(1e18))

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Chain", "Address", "Balance", "Symbol"})
		table.SetBorder(false)
		table.Append([]string{
			chainName,
			address,
			fmt.Sprintf("%.6f", ethValue),
			client.Config.Symbol,
		})
		table.Render()
	},
}

func init() {
	rootCmd.AddCommand(balanceCmd)
}
