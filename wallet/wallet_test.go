package wallet

import (
	"reflect"
	"testing"

	common2 "tokit/wallet/pkg/common2"

	"github.com/ethereum/go-ethereum/ethclient"
	bip32 "github.com/tyler-smith/go-bip32"
)

func TestNewWallet(t *testing.T) {
	clientFactory := func(chain string, rpcURL string) (common2.Blockchain, error) {
		// Mock client factory for testing
		return nil, nil
	}
	ethereumClient, _ := ethclient.Dial("http://localhost:8545")
	masterKey, _ := bip32.NewMasterKey([]byte("seed"))

	w := NewWallet(clientFactory, ethereumClient, masterKey)

	if w == nil {
		t.Errorf("NewWallet returned nil")
	}
	if reflect.TypeOf(w).String() != "*wallet.Wallet" {
		t.Errorf("NewWallet returned incorrect type: got %T, want *Wallet", w)
	}
}

func TestAddClient(t *testing.T) {
	w := &Wallet{
		clients: make(map[string]common2.Blockchain),
		clientFactory: func(chain string, rpcURL string) (common2.Blockchain, error) {
			return nil, nil // Mock client factory
		},
	}

	err := w.AddClient("testchain")
	if err != nil {
		t.Errorf("AddClient failed: %v", err)
	}

	if _, ok := w.clients["testchain"]; !ok {
		t.Errorf("Client for 'testchain' not added to wallet")
	}
}

func TestGetAddress(t *testing.T) {
	expectedAddress := "0x1234567890abcdef"
	clientFactory := func(chain string, rpcURL string) (common2.Blockchain, error) {
		return &MockBlockchain{address: expectedAddress, txHash: ""}, nil
	}
	w := &Wallet{
		clients:       make(map[string]common2.Blockchain),
		clientFactory: clientFactory,
	}
	w.AddClient("testchain")

	address, err := w.GetAddress("testchain")
	if err != nil {
		t.Errorf("GetAddress failed: %v", err)
	}
	if address != expectedAddress {
		t.Errorf("GetAddress got %s, want %s", address, expectedAddress)
	}

	_, err = w.GetAddress("nonexistent")
	if err == nil {
		t.Errorf("GetAddress with nonexistent chain should return an error")
	}
}

type MockBlockchain struct {
	address string
	txHash  string
}

func (m *MockBlockchain) GetAddress() (string, error) {
	return m.address, nil
}

func (m *MockBlockchain) Transfer(privateKeyHex string, toAddress string, amount float64) (string, error) {
	return m.txHash, nil
}

func (m *MockBlockchain) GetTransactionHistory(address string) ([]common2.Transaction, error) {
	return nil, nil
}

func (m *MockBlockchain) GenerateKey() ([]byte, error) {
	return []byte{}, nil
}

func TestTransfer(t *testing.T) {
	expectedTxHash := "0xabcdef1234567890"
	clientFactory := func(chain string, rpcURL string) (common2.Blockchain, error) {
		return &MockBlockchain{address: "", txHash: expectedTxHash}, nil
	}
	w := &Wallet{
		clients:       make(map[string]common2.Blockchain),
		clientFactory: clientFactory,
	}
	w.AddClient("testchain")

	txHash, err := w.Transfer("testchain", "0x0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef", "0xrecipient", 1.0)
	if err != nil {
		t.Errorf("Transfer failed: %v", err)
	}
	if txHash != expectedTxHash {
		t.Errorf("Transfer got %s, want %s", txHash, expectedTxHash)
	}

	_, err = w.Transfer("nonexistent", "0x0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef", "0xrecipient", 1.0)
	if err == nil {
		t.Errorf("Transfer with nonexistent chain should return an error")
	}
}
