# Tokit Wallet

A secure, modular, and feature-rich CLI wallet for EVM-compatible blockchains.

## Features

*   **🔐 Secure Key Management**:
    *   Uses `go-ethereum/accounts/keystore` (Scrypt N/P) for encrypted key storage.
    *   BIP39 Mnemonic generation and import (12/24 words).
    *   BIP44 HD Key Derivation (`m/44'/60'/0'/0/0`).
    *   Interactive password prompts (no passwords in history).

*   **⛓️ Multi-Chain Support**:
    *   Pre-configured for Ethereum, Arbitrum, Optimism, and Base.
    *   Easily extensible via `~/.tokit/config.yaml`.
    *   Supports any EVM-compatible network.

*   **💸 Transaction Management**:
    *   **EIP-1559** support (Dynamic Fee Transactions).
    *   **Smart Gas Estimation** for accurate fee calculation.
    *   **ERC20 Token Support**: Transfer and check balances of any ERC20 token.
    *   Secure signing with local keystore.

*   **🛠️ Developer Friendly**:
    *   Built with `Cobra` for a robust CLI experience.
    *   Configuration management via `Viper`.
    *   Structured logging and error handling.

## Installation

```bash
go build -o tokit.exe main.go
```

## Usage

### 1. Wallet Management

**Create a new wallet:**
```bash
./tokit wallet create
```
*Generates a new BIP39 mnemonic and imports the first account.*

**Import an existing wallet:**
```bash
./tokit wallet import
```
*Supports importing via Mnemonic phrase or Private Key (hex).*

**List accounts:**
```bash
./tokit wallet list
```

### 2. Balance Check

**Check ETH balance:**
```bash
./tokit balance ethereum [address]
```
*If address is omitted, checks the first local account.*

**Check ERC20 Token balance:**
```bash
./tokit balance ethereum --token 0xdac17f958d2ee523a2206206994597c13d831ec7
```

### 3. Transfer Funds

**Send ETH:**
```bash
./tokit transfer ethereum 0xRecipientAddress 0.1
```

**Send ERC20 Tokens:**
```bash
./tokit transfer ethereum 0xRecipientAddress 100 --token 0xdac17f958d2ee523a2206206994597c13d831ec7
```

## Configuration

The wallet uses a configuration file located at `~/.tokit/config.yaml`.

```yaml
default: ethereum
networks:
  ethereum:
    rpc: https://eth.llamarpc.com
    chain_id: 1
    symbol: ETH
    explorer: https://etherscan.io
  arbitrum:
    rpc: https://arb1.arbitrum.io/rpc
    chain_id: 42161
    symbol: ETH
    explorer: https://arbiscan.io
```

## Security

*   **Private Keys**: Stored in `~/.tokit/keystore` as encrypted JSON files.
*   **Passwords**: Never stored, only requested interactively for signing.
*   **Mnemonics**: Only displayed once during creation.

## License

MIT
