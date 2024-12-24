# Decentralized Wallet

This is a command-line decentralized wallet that supports multiple EVM-compatible blockchains, including Ethereum, Arbitrum, Optimism, and Base, as well as Bitcoin and EOS.

## Features

*   **Multi-Chain Support:** Supports Ethereum, Arbitrum, Optimism, Base, Bitcoin, and EOS.
*   **Address Generation:** Generates new addresses for supported blockchains.
*   ** फंड Transfer:** Allows transferring funds on EVM-compatible chains.
*   **Transaction History:** Retrieves transaction history for specified addresses on EVM-compatible chains.

## Usage

The wallet is controlled through command-line arguments.

### Environment Setup

Before using the wallet, you need to set the following environment variables. You can copy the `.env.template` file, rename it to `.env`, and replace the placeholders with your actual values:

*   `ETHEREUM_RPC_URL`: RPC URL for the Ethereum network.
*   `ARBITRUM_RPC_URL`: RPC URL for the Arbitrum network.
*   `OPTIMISM_RPC_URL`: RPC URL for the Optimism network.
*   `BASE_RPC_URL`: RPC URL for the Base network.

### Commands

*   `address <chain>`: Retrieves the address for the specified blockchain.
    *   Example: `tokit address ethereum`
    *   Supported chains: `ethereum`, `arbitrum`, `optimism`, `base`, `bitcoin`, `eos`

*   `transfer <chain> <to_address> <amount>`: Initiates a transfer on the specified blockchain.
    *   Example: `tokit transfer ethereum 0x1234abcd 1.0`
    *   Parameters:
        *   `<chain>`: The blockchain to perform the transfer on (`ethereum`, `arbitrum`, `optimism`, `base`).
        *   `<to_address>`: The recipient's address.
        *   `<amount>`: The transfer amount.

*   `history <chain> <address>`: Retrieves the transaction history for the specified blockchain and address.
    *   Example: `tokit history ethereum 0x1234abcd`
    *   Parameters:
        *   `<chain>`: The blockchain to query the transaction history from (`ethereum`, `arbitrum`, `optimism`, `base`).
        *   `<address>`: The address to query the transaction history for.

### Supported Chains

*   Ethereum
*   Arbitrum
*   Optimism
*   Base
*   Bitcoin (address generation only)
*   EOS (address generation only)

### Important Notes

*   This is a basic implementation and does not include secure key management or transaction signing functionalities.
*   RPC URLs are read from environment variables. Ensure these are correctly configured before running the wallet.
*   The wallet generates a new key pair for each address request. In a real-world scenario, key management would be handled more securely.
*   The wallet uses a mnemonic phrase to generate a master private key and then uses Hierarchical Deterministic (HD) paths to derive child private keys for each chain. This implementation uses a placeholder for key derivation and should be replaced with secure key generation and management practices for production use.
*   The `transfer` and `history` commands are only supported for EVM-compatible chains (Ethereum, Arbitrum, Optimism, Base).
