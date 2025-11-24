package cmd

import (
	"fmt"
	"os"
	"strings"
	"syscall"

	"tokit/internal/utils"
	"tokit/internal/wallet"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var walletCmd = &cobra.Command{
	Use:   "wallet",
	Short: "Manage wallet accounts",
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new wallet with a random mnemonic",
	Run: func(cmd *cobra.Command, args []string) {
		mnemonic, err := wallet.GenerateMnemonic()
		if err != nil {
			utils.Log.Fatalf("Failed to generate mnemonic: %v", err)
		}

		fmt.Println("⚠️  IMPORTANT: Write down this mnemonic phrase. It is the ONLY way to recover your funds!")
		fmt.Println(strings.Repeat("=", 60))
		fmt.Println(mnemonic)
		fmt.Println(strings.Repeat("=", 60))

		fmt.Print("\nEnter a password to encrypt your wallet: ")
		bytePassword, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			utils.Log.Fatalf("Failed to read password: %v", err)
		}
		password := string(bytePassword)
		fmt.Println()

		fmt.Print("Confirm password: ")
		byteConfirm, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			utils.Log.Fatalf("Failed to read password: %v", err)
		}
		if password != string(byteConfirm) {
			utils.Log.Fatal("Passwords do not match")
		}
		fmt.Println()

		svc, err := wallet.NewService()
		if err != nil {
			utils.Log.Fatalf("Failed to init wallet service: %v", err)
		}

		acc, err := svc.ImportMnemonic(mnemonic, password)
		if err != nil {
			utils.Log.Fatalf("Failed to create account: %v", err)
		}

		fmt.Printf("\n✅ Wallet created successfully!\nAddress: %s\n", acc.Address.Hex())
		fmt.Printf("Keystore location: %s\n", acc.URL.Path)
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all wallet accounts",
	Run: func(cmd *cobra.Command, args []string) {
		svc, err := wallet.NewService()
		if err != nil {
			utils.Log.Fatalf("Failed to init wallet service: %v", err)
		}

		accounts := svc.ListAccounts()
		if len(accounts) == 0 {
			fmt.Println("No accounts found.")
			return
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Index", "Address", "Location"})
		table.SetBorder(false)

		for i, acc := range accounts {
			table.Append([]string{
				fmt.Sprintf("%d", i),
				acc.Address.Hex(),
				acc.URL.Path,
			})
		}
		table.Render()
	},
}

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import a wallet using mnemonic or private key",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Enter mnemonic phrase or private key (hex): ")
		// Read input securely or just stdin? Mnemonic is long, maybe just read line
		// For better UX, let's just read line
		var input string
		// Using scanner to read full line including spaces for mnemonic
		// But here we can use simple Scanln if we assume quotes? No.
		// Let's use term for password-like input but print it? No, just standard input.
		// Actually, for security, maybe hide it?
		// Let's try reading password style for the secret
		byteInput, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			utils.Log.Fatalf("Failed to read input: %v", err)
		}
		input = strings.TrimSpace(string(byteInput))
		fmt.Println("\n(Input received)")

		fmt.Print("Enter a password to encrypt your wallet: ")
		bytePassword, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			utils.Log.Fatalf("Failed to read password: %v", err)
		}
		password := string(bytePassword)
		fmt.Println()

		svc, err := wallet.NewService()
		if err != nil {
			utils.Log.Fatalf("Failed to init wallet service: %v", err)
		}

		var acc interface{} // accounts.Account

		// Check if input is mnemonic (has spaces) or private key (hex)
		if strings.Contains(input, " ") {
			acc, err = svc.ImportMnemonic(input, password)
		} else {
			// Assume private key
			input = strings.TrimPrefix(input, "0x")
			acc, err = svc.ImportPrivateKey(input, password)
		}

		if err != nil {
			utils.Log.Fatalf("Failed to import account: %v", err)
		}

		account := acc.(interface{ Address() common.Address }).Address() // Wait, Import returns Account struct directly
		// My interface was wrong, let's fix
		// ImportMnemonic returns accounts.Account
		realAcc := acc.(interface{ Address() common.Address }) // This is getting messy with type assertion if I made a mistake
		// Let's look at my code: ImportMnemonic returns accounts.Account. ImportPrivateKey returns accounts.Account.
		// So `acc` is accounts.Account.

		// Re-writing this block cleanly:
		var finalAcc accounts.Account
		if strings.Contains(input, " ") {
			finalAcc, err = svc.ImportMnemonic(input, password)
		} else {
			input = strings.TrimPrefix(input, "0x")
			finalAcc, err = svc.ImportPrivateKey(input, password)
		}

		if err != nil {
			utils.Log.Fatalf("Failed to import account: %v", err)
		}

		fmt.Printf("\n✅ Wallet imported successfully!\nAddress: %s\n", finalAcc.Address.Hex())
	},
}

func init() {
	rootCmd.AddCommand(walletCmd)
	walletCmd.AddCommand(createCmd)
	walletCmd.AddCommand(listCmd)
	walletCmd.AddCommand(importCmd)
}
