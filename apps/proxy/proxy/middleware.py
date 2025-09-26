from fastapi import Request, HTTPException, status
from fastapi.responses import JSONResponse
from fastapi.middleware.base import BaseHTTPMiddleware
from slowapi import Limiter
from slowapi.util import get_remote_address
from slowapi.errors import RateLimitExceeded
import time
import json

# Rate limiting setup
limiter = Limiter(key_func=get_remote_address)

class RateLimitMiddleware(BaseHTTPMiddleware):
    def __init__(self, app, redis_url=None):
        super().__init__(app)
        self.redis_url = redis_url
        self.request_counts = {}  # Simple in-memory rate limiting for development

    async def dispatch(self, request: Request, call_next):
        client_ip = request.client.host
        current_time = time.time()

        # Simple rate limiting (100 requests per minute per IP)
        if client_ip not in self.request_counts:
            self.request_counts[client_ip] = []

        # Clean old requests (older than 1 minute)
        self.request_counts[client_ip] = [
            req_time for req_time in self.request_counts[client_ip]
            if current_time - req_time < 60
        ]

        # Check if rate limit exceeded
        if len(self.request_counts[client_ip]) >= 100:
            return JSONResponse(
                status_code=429,
                content={
                    "error": {
                        "code": "RATE_LIMITED",
                        "message": "Too many requests. Please try again later."
                    }
                }
            )

        # Add current request
        self.request_counts[client_ip].append(current_time)

        response = await call_next(request)
        return response

class ErrorHandlerMiddleware(BaseHTTPMiddleware):
    async def dispatch(self, request: Request, call_next):
        try:
            response = await call_next(request)
            return response
        except HTTPException as e:
            return JSONResponse(
                status_code=e.status_code,
                content={
                    "error": {
                        "code": self._get_error_code(e.status_code),
                        "message": e.detail
                    }
                }
            )
        except Exception as e:
            return JSONResponse(
                status_code=500,
                content={
                    "error": {
                        "code": "INTERNAL_ERROR",
                        "message": "Internal server error"
                    }
                }
            )

    def _get_error_code(self, status_code: int) -> str:
        error_codes = {
            400: "INVALID_INPUT",
            401: "UNAUTHORIZED",
            403: "FORBIDDEN",
            404: "NOT_FOUND",
            429: "RATE_LIMITED",
            500: "INTERNAL_ERROR"
        }
        return error_codes.get(status_code, "UNKNOWN_ERROR")

class SecurityHeadersMiddleware(BaseHTTPMiddleware):
    async def dispatch(self, request: Request, call_next):
        response = await call_next(request)

        # Add security headers
        response.headers["X-Content-Type-Options"] = "nosniff"
        response.headers["X-Frame-Options"] = "DENY"
        response.headers["X-XSS-Protection"] = "1; mode=block"
        response.headers["Strict-Transport-Security"] = "max-age=31536000; includeSubDomains"
        response.headers["Content-Security-Policy"] = "default-src 'self'"

        return response