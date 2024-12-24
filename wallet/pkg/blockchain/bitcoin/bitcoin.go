package bitcoin

import (
	"fmt"

	"tokit/wallet/pkg/blockchain"

	"crypto/ecdsa"
	"crypto/sha256"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"golang.org/x/crypto/ripemd160"
)

// Bitcoin implements the Blockchain interface for Bitcoin.
type Bitcoin struct {
	// Add any Bitcoin specific configurations here.
}

// NewBitcoin creates a new Bitcoin instance.
func NewBitcoin() *Bitcoin {
	return &Bitcoin{}
}

// GenerateKey generates a new Bitcoin private key.
func (b *Bitcoin) GenerateKey() ([]byte, error) {
	privateKey, err := btcec.NewPrivateKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %w", err)
	}
	return privateKey.Serialize(), nil
}

// GetAddress derives the Bitcoin address from the public key.
func (b *Bitcoin) GetAddress(key []byte) (string, error) {
	_, pubKey := btcec.PrivKeyFromBytes(key)
	serializedPubKey := pubKey.SerializeCompressed()

	hasherSHA256 := sha256.New()
	hasherSHA256.Write(serializedPubKey)
	hashSHA256 := hasherSHA256.Sum(nil)

	hasherRIPEMD160 := ripemd160.New()
	hasherRIPEMD160.Write(hashSHA256)
	hashRIPEMD160 := hasherRIPEMD160.Sum(nil)

	address, err := btcutil.NewAddressPubKeyHash(hashRIPEMD160, &chaincfg.MainNetParams)
	if err != nil {
		return "", fmt.Errorf("failed to create address: %w", err)
	}

	return address.EncodeAddress(), nil
}

// SignTransaction signs the Bitcoin transaction.
func (b *Bitcoin) SignTransaction(privateKey *ecdsa.PrivateKey, toAddress string, amount float64) (string, error) {
	fmt.Println("Signing Bitcoin transaction")
	// Placeholder for transaction signing logic.
	return "", nil
}

// Transfer transfers funds on the Bitcoin network.
func (b *Bitcoin) Transfer(fromKey *ecdsa.PrivateKey, toAddress string, amount float64) (string, error) {
	// Implementation for Bitcoin transfer
	fmt.Println("Transferring", amount, "BTC to", toAddress)
	// Need to implement the actual Bitcoin transaction creation and signing logic here.
	// This will involve using Bitcoin libraries to:
	// 1. Get the sender's UTXOs.
	// 2. Construct the transaction.
	// 3. Sign the transaction with the private key.
	// 4. Broadcast the transaction to the network.
	return "", fmt.Errorf("Bitcoin transfer functionality not yet implemented")
}

// GetTransactionHistory retrieves the transaction history for the given Bitcoin address.
func (b *Bitcoin) GetTransactionHistory(address string) ([]blockchain.Transaction, error) {
	fmt.Println("Getting Bitcoin transaction history")
	// Placeholder for transaction history retrieval logic.
	return nil, nil
}
