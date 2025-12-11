# Introduction

OpenComputeFramework (OCF) is a distributed service registry and peer-to-peer (P2P) network designed to facilitate decentralized computing. It enables nodes to discover each other, share state using CRDTs (Conflict-free Replicated Data Types), and proxy requests to registered services.

## Key Concepts

*   **Node**: An instance of the OCF software running on a machine.
*   **Peer**: Another node in the network.
*   **Service**: A computing service (e.g., an HTTP API) offered by a node.
*   **DNT (Distributed Node Table)**: A shared state containing information about all known nodes and their services.

## Architecture

OCF is built on top of [libp2p](https://libp2p.io/) and uses [IPFS](https://ipfs.tech/) technologies for data storage and retrieval. It supports different modes of operation, allowing for flexible deployment topologies.

## Getting Started

To get started with OpenComputeFramework, check out the following guides:

*   [Installation](installation.md): How to build and install the software.
*   [Configuration](configuration.md): How to configure your node.
*   [Usage](usage.md): How to run the CLI and start the server.
*   [API Reference](api.md): Documentation for the HTTP API.
