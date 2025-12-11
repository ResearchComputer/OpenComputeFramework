# Configuration

OpenComputeFramework uses `viper` for configuration management. It looks for a configuration file at `$HOME/.config/ocf/cfg.yaml` by default. You can also specify a config file using the `--config` flag.

## Configuration File (`cfg.yaml`)

Here is an example configuration file with default values:

```yaml
# Path to the data directory (default: "")
path: ""

# Main HTTP server port (default: "8092")
port: "8092"

# Node name (default: "relay")
name: "relay"

# Seed for generating node identity (default: "0")
seed: "0"

# TCP port for libp2p (default: "43905")
tcp_port: "43905"

# UDP port for libp2p (default: "59820")
udp_port: "59820"

# P2P configuration
p2p:
  port: "8093"

# Vacuum configuration
vacuum:
  interval: 10

# Queue configuration
queue:
  port: "8094"

# Account configuration
account:
  wallet: ""

# Solana integration configuration
solana:
  rpc: "https://api.mainnet-beta.solana.com"
  mint: "EsmcTrdLkFqV3mv4CjLF3AmCx132ixfFSYYRWD78cDzR"
  skip_verification: false
```

## Environment Variables

Configuration options can also be overridden using environment variables. The mapping is typically uppercase with `_` separators, but check the code for specific bindings if needed.

## CLI Flags

Many configuration options can be passed as CLI flags when starting the application. See the [Usage](usage.md) section for more details.

## CRDT Configuration

The following CRDT-related settings can be configured (defaults shown):

*   `crdt.tombstone_retention`: `24h`
*   `crdt.tombstone_compaction_interval`: `1h`
*   `crdt.tombstone_compaction_batch`: `512`
