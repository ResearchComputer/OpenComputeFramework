#!/bin/bash

# Load environment variables from .env file if it exists
if [ -f "../.env" ]; then
    export $(cat ../.env | grep -v '^#' | xargs)
fi

# Start Docker container with environment variables
docker run \
    -e ENVIRONMENT="${ENVIRONMENT:-development}" \
    -e PG_URI="${PG_URI}" \
    -e REDIS_URL="${REDIS_URL}" \
    -e OCF_HEAD_URL="${OCF_HEAD_URL}" \
    -e JWT_SECRET_KEY="${JWT_SECRET_KEY}" \
    -p 10590:10590 \
    researchcomputer/ocf-proxy