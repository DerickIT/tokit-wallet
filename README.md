# Tokit Wallet

A secure, modular CLI wallet for EVM-compatible blockchains.

## Features

- Secure key management with encrypted keystore files (`go-ethereum/accounts/keystore`).
- BIP39 mnemonic generation/import and BIP44 HD derivation (`m/44'/60'/0'/0/0`).
- Multi-chain support (Ethereum, Arbitrum, Optimism, Base), configurable via `~/.tokit/config.yaml`.
- Native token and ERC20 balance query.
- EIP-1559 transfers with dynamic fee estimation.
- Interactive password prompts for import/signing.

## Installation

```bash
go build -o tokit.exe main.go
```

## Usage

### Wallet management

Create a wallet:

```bash
./tokit wallet create
```

Import a wallet:

```bash
./tokit wallet import
```

List wallets:

```bash
./tokit wallet list
```

Set default account:

```bash
./tokit wallet use 0xYourAddress
```

### Balance

Check native balance (address optional):

```bash
./tokit balance ethereum [address]
```

If address is omitted, Tokit uses `--from`, then configured default account, then the first local account.

Check ERC20 balance:

```bash
./tokit balance ethereum --token 0xdac17f958d2ee523a2206206994597c13d831ec7
```

### Transfer

Send native token:

```bash
./tokit transfer ethereum 0xRecipientAddress 0.1
```

Send ERC20 token:

```bash
./tokit transfer ethereum 0xRecipientAddress 100 --token 0xdac17f958d2ee523a2206206994597c13d831ec7
```

## Configuration

Tokit uses `~/.tokit/config.yaml`:

```yaml
default: ethereum
default_account: ""
networks:
  ethereum:
    rpc: https://eth.llamarpc.com
    chain_id: 1
    symbol: ETH
    explorer: https://etherscan.io
```

## SSH proxy troubleshooting (GitHub push)

If `git push` shows an error like `exec : The term 'exec' is not recognized` on Windows:

- Check whether environment variable `SHELL` is set to PowerShell.
- Remove it for the current session and retry:

```powershell
Remove-Item Env:SHELL -ErrorAction SilentlyContinue
git push --porcelain origin
```

This issue usually appears when SSH `ProxyCommand` is interpreted by the wrong shell.

## Security notes

- Private keys are stored as encrypted JSON files under `~/.tokit/keystore`.
- Passwords are only entered interactively.
- Mnemonic is shown once during creation and should be stored offline securely.

## License

MIT
