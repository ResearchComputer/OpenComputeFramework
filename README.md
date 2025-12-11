# OpenComputeFramework

OpenComputeFramework is a distributed service registry and P2P network built with Go. It utilizes libp2p and IPFS technologies to provide a decentralized infrastructure for computing services.

## Features

- **P2P Networking**: Leverages libp2p for robust peer-to-peer communication.
- **Distributed Service Registry**: Allows nodes to register and discover services.
- **CRDT-based State Management**: Uses Conflict-free Replicated Data Types for eventual consistency.
- **Solana Integration**: Optional integration with Solana for token-gated access (verification).

## Documentation

Full documentation is available in the `docs/` directory. You can also view it by running the documentation server (if configured) or browsing the markdown files directly.

- [Installation](docs/docs/installation.md)
- [Configuration](docs/docs/configuration.md)
- [Usage](docs/docs/usage.md)
- [API Reference](docs/docs/api.md)

## quick start

```bash
cd src
make build
./build/ocfcore start
```
