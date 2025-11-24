package cmd

import (
	"fmt"
	"math/big"
	"os"
	"text/tabwriter"
	"tokit/internal/chain"
	"tokit/internal/utils"
	"tokit/internal/wallet"

	"github.com/spf13/cobra"
)

var balanceTokenAddress string

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

		var balance *big.Int
		var symbol string

		if balanceTokenAddress != "" {
			// Check Token Balance
			balance, err = client.GetTokenBalance(balanceTokenAddress, address)
			if err != nil {
				utils.Log.Fatalf("Failed to get token balance: %v", err)
			}
			symbol = "TOKEN" // TODO: Fetch symbol from contract
		} else {
			// Check ETH Balance
			balance, err = client.GetBalance(address)
			if err != nil {
				utils.Log.Fatalf("Failed to get balance: %v", err)
			}
			symbol = client.Config.Symbol
		}

		// Convert Wei to Ether (assuming 18 decimals)
		fBalance := new(big.Float)
		fBalance.SetString(balance.String())
		ethValue := new(big.Float).Quo(fBalance, big.NewFloat(1e18))

		w := new(tabwriter.Writer)
		w.Init(os.Stdout, 0, 8, 2, '\t', 0)
		fmt.Fprintln(w, "Chain\tAddress\tBalance\tSymbol")
		fmt.Fprintf(w, "%s\t%s\t%.6f\t%s\n", chainName, address, ethValue, symbol)
		w.Flush()
	},
}

func init() {
	rootCmd.AddCommand(balanceCmd)
	balanceCmd.Flags().StringVarP(&balanceTokenAddress, "token", "t", "", "ERC20 token address")
}
