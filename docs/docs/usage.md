# Usage

## CLI Commands

The `ocfcore` CLI provides several commands to interact with the system.

### Global Flags

*   `--config string`: config file (default is `$HOME/.config/ocf/cfg.yaml`)

### `init`

Initialize the node configuration.

```bash
ocfcore init
```

### `start`

Start the OCF node and listen for incoming connections.

```bash
ocfcore start [flags]
```

#### Flags

*   `--wallet.account string`: Wallet account.
*   `--account.wallet string`: Path to wallet key file.
*   `--bootstrap.addr string`: Bootstrap address (default "http://152.67.71.5:8092/v1/dnt/bootstraps").
*   `--bootstrap.source stringSlice`: Bootstrap source (HTTP URL, dnsaddr://host, or multiaddr). Repeatable.
*   `--bootstrap.static stringSlice`: Static bootstrap multiaddr (repeatable).
*   `--seed string`: Seed for identity generation (default "0").
*   `--mode string`: Mode (standalone, local, full) (default "node").
*   `--tcpport string`: TCP Port (default "43905").
*   `--udpport string`: UDP Port (default "59820").
*   `--subprocess string`: Subprocess to start.
*   `--public-addr string`: Public address if you have one (enables being a bootstrap node).
*   `--service.name string`: Service name to register.
*   `--service.port string`: Service port.
*   `--solana.rpc string`: Solana RPC endpoint.
*   `--solana.mint string`: SPL token mint to verify ownership.
*   `--solana.skip_verification`: Skip Solana token ownership verification (use for testing only).
*   `--cleanslate`: Clean slate, removing the database before starting (default `true`).

### `version`

Print the version number of ocfcore.

```bash
ocfcore version
```

### `update`

Update the application (if supported).

```bash
ocfcore update
```

## Running a Node

To start a standard node:

```bash
./ocfcore start
```

To start a node with a specific seed (identity):

```bash
./ocfcore start --seed myseed
```

To register a local service:

```bash
./ocfcore start --service.name myservice --service.port 8080
```
