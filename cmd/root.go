package cmd

import (
	"fmt"
	"os"

	"tokit/internal/config"
	"tokit/internal/utils"

	"github.com/spf13/cobra"
)

var (
	cfgFile   string
	Verbose   bool
	AppConfig *config.Config
)

var rootCmd = &cobra.Command{
	Use:   "tokit",
	Short: "Tokit - A secure multi-chain wallet",
	Long: `Tokit is a CLI wallet that supports multiple EVM-compatible blockchains.
It provides secure key management using encrypted keystores and supports
standard wallet operations like transfers and balance checks.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		utils.InitLogger(Verbose)
		var err error
		AppConfig, err = config.LoadConfig()
		if err != nil {
			utils.Log.Fatalf("Failed to load config: %v", err)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
}
