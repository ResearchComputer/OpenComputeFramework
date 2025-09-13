# OpenAI Compatible Proxy Service

A FastAPI-based proxy service that asynchronously routes requests to OpenAI-compatible services.

## Features

- Async request handling with FastAPI
- Proxy for OpenAI-compatible endpoints:
  - `/v1/chat/completions`
  - `/v1/completions`
  - `/v1/embeddings`
  - `/v1/models`
- Timeout handling
- Error handling and status code preservation
- Header forwarding
- Docker support

## Environment Variables

- `TARGET_SERVICE_URL`: Target OpenAI-compatible service URL (default: `http://localhost:8000/v1`)
- `TIMEOUT`: Request timeout in seconds (default: `30.0`)

## Running

### Local Development

```bash
pip install -r requirements.txt
python main.py
```

### Docker

```bash
docker build -t openai-proxy .
docker run -p 8000:8000 -e TARGET_SERVICE_URL=http://target-service:8000/v1 openai-proxy
```

### Docker Compose

```yaml
version: '3.8'
services:
  proxy:
    build: .
    ports:
      - "8000:8000"
    environment:
      - TARGET_SERVICE_URL=http://target-service:8000/v1
      - TIMEOUT=30.0
```