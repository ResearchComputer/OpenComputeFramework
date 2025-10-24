---
title: CLI
description: Command line options for the OCF node.
---

# Usage

The released binaries are named by architecture (e.g., `ocf-amd64`). The CLI root command is `ocfcore` in source. Examples below use the binary:

```bash
./ocf-amd64 --help | cat
```

# Commands

- **start**: Start the node and HTTP API.
- **wallet**: Wallet management commands for node owner identification.
  - **create**: Generate a new Solana account managed by OCF.
  - **list**: Display all accounts managed by OCF (default account marked with `*`).
  - **info**: Show the default account information.
- **version**: Print version information.
- **init**: Initialize config (placeholder).
- **update**: Update the node (placeholder).

# Flags (start)

- `--wallet.account` string: Wallet account for node identification (defaults to the first managed account).
- `--account.wallet` string: Path to the keypair file for the managed account (auto-populated when using managed wallets).
- `--bootstrap.addr` string: Bootstrap source (HTTP URL returning `{"bootstraps": ["/ip4/x/tcp/43905/p2p/<ID>"]}` or a single multiaddr). Default: `http://152.67.71.5:8092/v1/dnt/bootstraps`.
- `--seed` string: Seed for deterministic peer key (use `0` to persist/load key). Default: `0`.
- `--mode` string: `standalone`, `local`, or `node`/`full` (default `node`).
- `--tcpport` string: LibP2P TCP port. Default `43905`.
- `--udpport` string: LibP2P QUIC UDP port. Default `59820`.
- `--subprocess` string: Start a critical subprocess (kept alive by OCF).
- `--public-addr` string: Public IP address to advertise (enables bootstrap role).
- `--service.name` string: Local service name to register (e.g., `llm`).
- `--service.port` string: Local service port to register (e.g., `8080`).
- `--solana.rpc` string: Solana RPC endpoint used for SPL token verification. Default: `https://api.mainnet-beta.solana.com`.
- `--solana.mint` string: SPL token mint the node must hold (default `EsmcTrdLkFqV3mv4CjLF3AmCx132ixfFSYYRWD78cDzR`).
- `--solana.skip_verification` bool: Skip SPL token balance checks (testing only). Default `false`.
- `--cleanslate` bool: Remove local CRDT db on start. Default `true`.

# Global Flags

- `--config` string: Config file path (default is `$HOME/.config/ocf/cfg.yaml`).

# Examples

Start a standalone dispatcher:

```bash
./ocf-amd64 start --mode standalone
```

Start a node using a known bootstrap multiaddr:

```bash
./ocf-amd64 start --bootstrap.addr=/ip4/1.2.3.4/tcp/43905/p2p/<BOOTSTRAP_PEER_ID>
```

Advertise as a public bootstrap:

```bash
./ocf-amd64 start --public-addr=203.0.113.10
```

Register a local LLM service (worker):

```bash
./ocf-amd64 start \
  --service.name=llm \
  --service.port=8080
```

# Wallet Management

Create a new wallet for node identification:

```bash
./ocf-amd64 wallet create
```

List managed wallets:

```bash
./ocf-amd64 wallet list
```

Show the default wallet information:

```bash
./ocf-amd64 wallet info
```

After creating a wallet, start the node (the CLI automatically wires the default account):

```bash
./ocf-amd64 start
```

To point at a different managed account, pass its public key explicitly:

```bash
./ocf-amd64 start --wallet.account=<wallet_address>
```

