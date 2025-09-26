# OCF Proxy Service

A FastAPI-based OpenAI-compatible proxy service for the Open Compute Framework with Web3 wallet authentication and API key management.

## Features

- OpenAI-compatible API endpoints (/v1/*)
- Web3 wallet authentication with JWT tokens
- API key management with usage tracking
- Database integration (PostgreSQL)
- Rate limiting and security middleware
- CORS support
- Async request handling
- Timeout management
- Docker containerization

## Architecture

The proxy service provides:

### Authentication & User Management
- Wallet signature verification
- JWT token generation and validation
- User profile management

### API Key Management
- Secure API key generation (sk_rc_* format)
- API key creation, listing, and deactivation
- Usage tracking and logging
- Rate limiting per API key

### Models API
- Fetch available models from OCF
- OpenAI-compatible inference endpoints
- Usage logging and cost tracking

## Quick Start

### Local Development

1. Install dependencies:
```bash
pip install -r requirements.txt
```

2. Set environment variables:
```bash
cp .env.example .env
# Edit .env with your configuration
export $(cat .env | xargs)
```

3. Run the service:
```bash
python -m proxy.main
```

### Docker

1. Build the image:
```bash
./scripts/build_docker.sh
```

2. Run the container:
```bash
docker run -p 8000:8000 \
  -e PG_URI=postgresql://user:pass@host:5432/db \
  -e JWT_SECRET_KEY=your-secret-key \
  researchcomputer/ocf-proxy
```

## API Endpoints

### Authentication Endpoints

- `POST /api/auth/connect` - Connect wallet and get JWT token
- `GET /api/user/profile` - Get user profile (requires JWT)

### API Key Management

- `POST /api/api-keys` - Create new API key
- `GET /api/api-keys` - List user's API keys
- `DELETE /api/api-keys/{key_id}` - Deactivate API key
- `PATCH /api/api-keys/{key_id}/usage` - Update usage timestamp

### Models API

- `GET /api/models` - Get models (requires authentication)
- `GET /api/models/public` - Get models (public access)
- `POST /api/v1/chat/completions` - Chat completions (requires API key)
- `POST /api/v1/completions` - Text completions (requires API key)

### Legacy OpenAI Compatible Endpoints

- `POST /v1/chat/completions` - Chat completions (legacy)
- `POST /v1/completions` - Text completions (legacy)
- `POST /v1/embeddings` - Embeddings (legacy)
- `GET /v1/models` - List models (legacy)

### System Endpoints

- `GET /health` - Health check
- `GET /` - API information

## Environment Variables

- `PG_URI` - PostgreSQL database connection string (required)
- `OCF_HEAD_URL` - Target OCF service URL (default: http://140.238.223.116:8092)
- `JWT_SECRET_KEY` - JWT secret key for authentication (required)
- `REDIS_URL` - Redis URL for rate limiting (optional)

## Database Schema

### Users Table
- `id` - UUID primary key
- `address` - Wallet address (unique, indexed)
- `created_at` - Creation timestamp
- `updated_at` - Last update timestamp

### API Keys Table
- `id` - UUID primary key
- `user_id` - Foreign key to users
- `name` - API key name
- `key` - API key (unique, indexed)
- `key_hash` - Hashed API key for security
- `created_at` - Creation timestamp
- `last_used` - Last usage timestamp
- `is_active` - Active status

### API Usage Logs Table
- `id` - UUID primary key
- `api_key_id` - Foreign key to API keys
- `model_id` - Model used
- `input_tokens` - Input token count
- `output_tokens` - Output token count
- `timestamp` - Request timestamp
- `cost` - Request cost

## Security Features

- JWT token authentication with expiration
- API key hashing with bcrypt
- Rate limiting (100 requests/minute per IP)
- Request logging and usage tracking
- CORS and security headers
- Error handling with proper HTTP status codes

## Testing

Run the test suite:
```bash
python -m pytest tests/
```

Or run the simple implementation test:
```bash
python test_implementation.py
```

## Frontend Integration

The proxy service is designed to work with the Web3 AI Model Console frontend. Key integration points:

1. **Wallet Connection**: Use `/api/auth/connect` with signature verification
2. **API Keys**: Generate keys via `/api/api-keys` for inference requests
3. **Model Discovery**: Fetch models via `/api/models` or `/api/models/public`
4. **Inference**: Use OpenAI-compatible endpoints with Bearer token authentication

## Error Response Format

```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Human readable error message"
  }
}
```

Common error codes:
- `UNAUTHORIZED` - Invalid authentication
- `FORBIDDEN` - Insufficient permissions
- `NOT_FOUND` - Resource not found
- `RATE_LIMITED` - Too many requests
- `INVALID_INPUT` - Malformed request data
- `INTERNAL_ERROR` - Server-side error