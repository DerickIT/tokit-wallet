package wallet

import (
	"errors"
	"testing"

	"tokit/wallet/pkg/common2"
)

// MockBlockchainClient is a mock implementation of the Blockchain interface.
type MockBlockchainClient struct {
	address         string
	err             error
	transferSuccess bool
	history         []common2.Transaction
}

func (m *MockBlockchainClient) GetAddress() (string, error) {
	return m.address, m.err
}

func (m *MockBlockchainClient) Transfer(privateKeyHex string, toAddress string, amount float64) (string, error) {
	if m.transferSuccess {
		return "txHash", nil
	}
	return "", errors.New("transfer failed")
}

func (m *MockBlockchainClient) GetTransactionHistory(address string) ([]common2.Transaction, error) {
	return m.history, m.err
}

func (m *MockBlockchainClient) GenerateKey() ([]byte, error) {
	return nil, nil
}

func TestNewWallet(t *testing.T) {
	_ = NewWallet(nil, nil)
}

func TestWallet_AddClient(t *testing.T) {
	w := NewWallet(func(chain string, rpcURL string) (common2.Blockchain, error) {
		return &MockBlockchainClient{}, nil
	}, nil)

	err := w.AddClient("ethereum")
	if err != nil {
		t.Errorf("AddClient failed: %v", err)
	}

	if _, ok := w.clients["ethereum"]; !ok {
		t.Error("Client not added")
	}
}

func TestWallet_GetAddress(t *testing.T) {
	mockAddress := "0x123"
	w := NewWallet(func(chain string, rpcURL string) (common2.Blockchain, error) {
		return &MockBlockchainClient{address: mockAddress}, nil
	}, nil)
	w.clients["ethereum"] = &MockBlockchainClient{address: mockAddress}

	addr, err := w.GetAddress("ethereum")
	if err != nil {
		t.Fatalf("GetAddress failed: %v", err)
	}
	if addr != mockAddress {
		t.Errorf("Expected address %s, got %s", mockAddress, addr)
	}

	_, err = w.GetAddress("unknown")
	if err == nil {
		t.Error("Expected error for unknown chain")
	}
}

func TestWallet_GetAddress_ClientError(t *testing.T) {
	expectedErr := errors.New("client error")
	w := NewWallet(func(chain string, rpcURL string) (common2.Blockchain, error) {
		return &MockBlockchainClient{err: expectedErr}, nil
	}, nil)
	w.clients["ethereum"] = &MockBlockchainClient{err: expectedErr}

	_, err := w.GetAddress("ethereum")
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if !errors.Is(err, expectedErr) {
		t.Errorf("Expected error %v, got %v", expectedErr, err)
	}
}

func TestWallet_Transfer(t *testing.T) {
	w := NewWallet(func(chain string, rpcURL string) (common2.Blockchain, error) {
		return &MockBlockchainClient{transferSuccess: true}, nil
	}, nil)
	w.clients["ethereum"] = &MockBlockchainClient{transferSuccess: true}

	txHash, err := w.Transfer("ethereum", "0xprivateKey", "0xtoAddress", 1.0)
	if err != nil {
		t.Fatalf("Transfer failed: %v", err)
	}
	if txHash != "txHash" {
		t.Errorf("Expected txHash %s, got %s", "txHash", txHash)
	}

	w.clients["ethereum"] = &MockBlockchainClient{transferSuccess: false}
	_, err = w.Transfer("ethereum", "0xprivateKey", "0xtoAddress", 1.0)
	if err == nil {
		t.Error("Expected error for transfer failure")
	}

	_, err = w.Transfer("unknown", "0xprivateKey", "0xtoAddress", 1.0)
	if err == nil {
		t.Error("Expected error for unknown chain")
	}
}

func TestWallet_GetTransactionHistory(t *testing.T) {
	mockHistory := []common2.Transaction{{Hash: "0xabc"}}
	w := NewWallet(func(chain string, rpcURL string) (common2.Blockchain, error) {
		return &MockBlockchainClient{history: mockHistory}, nil
	}, nil)
	w.clients["ethereum"] = &MockBlockchainClient{history: mockHistory}

	history, err := w.GetTransactionHistory("ethereum", "0xaddress")
	if err != nil {
		t.Fatalf("GetTransactionHistory failed: %v", err)
	}
	if len(history) != 1 || history[0].Hash != "0xabc" {
		t.Errorf("Expected history %v, got %v", mockHistory, history)
	}

	_, err = w.GetTransactionHistory("unknown", "0xaddress")
	if err == nil {
		t.Error("Expected error for unknown chain")
	}

	expectedErr := errors.New("client error")
	w.clients["ethereum"] = &MockBlockchainClient{err: expectedErr}
	_, err = w.GetTransactionHistory("ethereum", "0xaddress")
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if !errors.Is(err, expectedErr) {
		t.Errorf("Expected error %v, got %v", expectedErr, err)
	}
}
