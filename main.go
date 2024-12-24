package main

import (
	"fmt"
	"os"
	"tokit/wallet"
	"tokit/wallet/pkg/blockchain/evm"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/tyler-smith/go-bip32"
)

func main() {
	// Initialize wallet factory (for EVM chains)
	const defaultRPCURL = "http://localhost:8545"
	rpcURL := os.Getenv("RPC_URL")
	if rpcURL == "" {
		fmt.Println("Warning: RPC_URL environment variable not set. Using default:", defaultRPCURL)
		rpcURL = defaultRPCURL
	}

	// Create Ethereum client
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		panic(err)
	}

	// Get wallet seed from environment variable
	seedPhrase := os.Getenv("WALLET_SEED")
	if seedPhrase == "" {
		panic("WALLET_SEED environment variable not set")
	}
	seed := []byte(seedPhrase)
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		panic(err)
	}

	// Create wallet
	myWallet := wallet.NewWallet(evm.Create, client, masterKey)

	ethereumChain := os.Getenv("ETHEREUM_CHAIN_NAME")
	if ethereumChain == "" {
		panic("ETHEREUM_CHAIN_NAME environment variable not set")
	}
	arbitrumChain := os.Getenv("ARBITRUM_CHAIN_NAME")
	if arbitrumChain == "" {
		panic("ARBITRUM_CHAIN_NAME environment variable not set")
	}
	optimismChain := os.Getenv("OPTIMISM_CHAIN_NAME")
	if optimismChain == "" {
		panic("OPTIMISM_CHAIN_NAME environment variable not set")
	}
	baseChain := os.Getenv("BASE_CHAIN_NAME")
	if baseChain == "" {
		panic("BASE_CHAIN_NAME environment variable not set")
	}

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
