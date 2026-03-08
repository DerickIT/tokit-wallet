package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"
	"tokit/internal/chain"
	"tokit/internal/utils"
	"tokit/internal/wallet"

	"github.com/spf13/cobra"
)

var balanceTokenAddress string
var balanceFromAddress string

var balanceCmd = &cobra.Command{
	Use:   "balance [chain] [address]",
	Short: "Check account balance",
	Long:  `Check the balance of an account on a specific blockchain. If address is omitted, checks the selected local wallet account.`,
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
			svc, err := wallet.NewService()
			if err != nil {
				utils.Log.Fatalf("Failed to init wallet service: %v", err)
			}
			account, err := resolveLocalAccount(svc, balanceFromAddress)
			if err != nil {
				utils.Log.Fatalf("Failed to resolve local account: %v", err)
			}
			address = account.Address.Hex()
		}

		client, err := chain.NewClient(chainName, AppConfig)
		if err != nil {
			utils.Log.Fatalf("Failed to create client: %v", err)
		}
		defer client.Close()

		displayBalance := "0"
		symbol := client.Config.Symbol

		if balanceTokenAddress != "" {
			metadata, err := client.GetTokenMetadata(balanceTokenAddress)
			if err != nil {
				utils.Log.Fatalf("Failed to get token metadata: %v", err)
			}
			balance, err := client.GetTokenBalance(balanceTokenAddress, address)
			if err != nil {
				utils.Log.Fatalf("Failed to get token balance: %v", err)
			}
			displayBalance = chain.FormatUnits(balance, metadata.Decimals, 6)
			symbol = metadata.Symbol
		} else {
			balance, err := client.GetBalance(address)
			if err != nil {
				utils.Log.Fatalf("Failed to get balance: %v", err)
			}
			displayBalance = chain.FormatUnits(balance, 18, 6)
		}

		w := new(tabwriter.Writer)
		w.Init(os.Stdout, 0, 8, 2, '\t', 0)
		fmt.Fprintln(w, "Chain\tAddress\tBalance\tSymbol")
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", chainName, address, displayBalance, symbol)
		w.Flush()
	},
}

func init() {
	rootCmd.AddCommand(balanceCmd)
	balanceCmd.Flags().StringVarP(&balanceTokenAddress, "token", "t", "", "ERC20 token address")
	balanceCmd.Flags().StringVar(&balanceFromAddress, "from", "", "Local account address to use when address is omitted")
}
