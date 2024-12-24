package main

import (
	"fmt"
	"os"
	"tokit/wallet"
	"tokit/wallet/pkg/blockchain/evm"

	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	// Initialize wallet factory (for EVM chains)
	rpcURL := os.Getenv("RPC_URL") // You might want to use different RPC URLs for different chains
	if rpcURL == "" {
		rpcURL = "http://localhost:8545" // Default if not set
	}

	// Create Ethereum client
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		panic(err)
	}

	// Create wallet
	myWallet := wallet.NewWallet(evm.Create, client)
	err = myWallet.AddClient("ethereum")
	if err != nil {
		panic(err)
	}
	err = myWallet.AddClient("arbitrum")
	if err != nil {
		panic(err)
	}
	err = myWallet.AddClient("optimism")
	if err != nil {
		panic(err)
	}
	err = myWallet.AddClient("base")
	if err != nil {
		panic(err)
	}

	// Get Ethereum address
	ethAddress, err := myWallet.GetAddress("ethereum")
	if err != nil {
		fmt.Println("Error getting Ethereum address:", err)
	} else {
		fmt.Println("Ethereum Address:", ethAddress)
	}

	// Transfer funds (example)
	privateKey := "your_private_key_here" // Replace with an actual private key
	toAddress := "recipient_address_here" // Replace with a recipient address
	amount := 0.01

	txHash, err := myWallet.Transfer("ethereum", privateKey, toAddress, amount)
	if err != nil {
		fmt.Println("Error transferring funds:", err)
	} else {
		fmt.Println("Transaction Hash:", txHash)
	}

	// Get transaction history (example)
	history, err := myWallet.GetTransactionHistory("ethereum", ethAddress)
	if err != nil {
		fmt.Println("Error getting transaction history:", err)
	} else {
		fmt.Println("Transaction History:", history)
	}
}
