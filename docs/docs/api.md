# API Reference

The OpenComputeFramework provides a RESTful API for interacting with the node. The API is available at `http://localhost:8092` by default.

## Health

### Check server health

`GET /v1/health`

Returns a simple health status check.

**Response:**

```json
{
  "status": "ok"
}
```

## DNT (Distributed Node Table)

### Get node table

`GET /v1/dnt/table`

Retrieve the current distributed node table.

### List connected peers

`GET /v1/dnt/peers`

Get a list of all currently connected peers.

### List peers with status

`GET /v1/dnt/peers_status`

Get a list of all peers with their connection status.

### List bootstraps

`GET /v1/dnt/bootstraps`

Get a list of all connected bootstrap nodes.

### Get resource statistics

`GET /v1/dnt/stats`

Retrieve resource usage and connection statistics.

### Update local node

`POST /v1/dnt/_node`

Update the local node's information in the node table.

### Delete local node

`DELETE /v1/dnt/_node`

Remove the local node from the node table.

## P2P Proxy

### Forward request to peer

`METHOD /v1/p2p/{peerId}/{path}`

Forward an HTTP request to a specific peer in the P2P network.

*   `peerId`: The ID of the target peer.
*   `path`: The path to forward the request to.

Supported methods: `GET`, `POST`, `PATCH`.

## Service Proxy

### Forward request to global service

`METHOD /v1/service/{service}/{path}`

Forward an HTTP request to a globally registered service. The system will route the request to a node offering this service.

*   `service`: The name of the service to forward to.
*   `path`: The path to forward the request to.

Supported methods: `GET`, `POST`, `PATCH`.

### Forward request to service (Local/Direct)

`METHOD /v1/_service/{service}/{path}`

Forward an HTTP request to a registered service (potentially with different routing rules or intended for local consumption/debugging).

*   `service`: The name of the service to forward to.
*   `path`: The path to forward the request to.

Supported methods: `GET`, `POST`, `PATCH`.
