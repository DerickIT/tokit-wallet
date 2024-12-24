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
	rpcURL := os.Getenv("RPC_URL")
	if rpcURL == "" {
		const defaultRPCURL = "http://localhost:8545"
		// Consider logging a warning here
		fmt.Println("Warning: RPC_URL environment variable not set. Using default.")
		rpcURL = defaultRPCURL
	}

	// Create Ethereum client
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		panic(err)
	}

	// Create wallet
	myWallet := wallet.NewWallet(evm.Create, client)
	const (
		ethereumChain = "ethereum"
		arbitrumChain = "arbitrum"
		optimismChain = "optimism"
		baseChain     = "base"
	)

	err = myWallet.AddClient(ethereumChain)
	if err != nil {
		panic(err)
	}
	err = myWallet.AddClient(arbitrumChain)
	if err != nil {
		panic(err)
	}
	err = myWallet.AddClient(optimismChain)
	if err != nil {
		panic(err)
	}
	err = myWallet.AddClient(baseChain)
	if err != nil {
		panic(err)
	}

	// Get Ethereum address
	ethAddress, err := myWallet.GetAddress(ethereumChain)
	if err != nil {
		fmt.Println("Error getting Ethereum address:", err)
	} else {
		fmt.Println("Ethereum Address:", ethAddress)
	}

	// Get transaction history (example)
	history, err := myWallet.GetTransactionHistory(ethereumChain, ethAddress)
	if err != nil {
		fmt.Println("Error getting transaction history:", err)
	} else {
		fmt.Println("Transaction History:", history)
	}
}
