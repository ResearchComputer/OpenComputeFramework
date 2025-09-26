import os
import json
import aiohttp
import asyncio
from typing import Optional, Dict, Any
from fastapi.responses import JSONResponse
from fastapi.middleware.cors import CORSMiddleware
from fastapi import FastAPI, Request, HTTPException
from .utils import get_all_models
from .database import init_db
from .routers import router as auth_router
from .models_api import router as models_router
from .middleware import RateLimitMiddleware, ErrorHandlerMiddleware, SecurityHeadersMiddleware

app = FastAPI(title="OpenAI Compatible Proxy Service")

# Add middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)
app.add_middleware(ErrorHandlerMiddleware)
app.add_middleware(SecurityHeadersMiddleware)
app.add_middleware(RateLimitMiddleware)
OCF_HEAD_URL = os.getenv("OCF_HEAD_URL", "http://140.238.223.116:8092")

async def proxy_request(
    session: aiohttp.ClientSession,
    method: str,
    url: str,
    headers: Dict[str, str],
    content: Optional[bytes] = None,
    params: Optional[Dict[str, Any]] = None
) -> JSONResponse:
    try:
        async with session.request(
            method=method,
            url=url,
            headers=headers,
            data=content,
            params=params
        ) as response:
            response_data = await response.json()
            return JSONResponse(
                status_code=response.status,
                content=response_data,
                headers=dict(response.headers)
            )
    except asyncio.TimeoutError:
        raise HTTPException(status_code=504, detail="Request timeout")
    except aiohttp.ClientError as e:
        raise HTTPException(status_code=502, detail=f"Request failed: {str(e)}")
    except json.JSONDecodeError:
        return JSONResponse(
            status_code=response.status,
            content={"error": "Invalid JSON response"},
            headers=dict(response.headers)
        )

@app.api_route("/v1/chat/completions", methods=["POST"])
@app.api_route("/v1/completions", methods=["POST"])
@app.api_route("/v1/embeddings", methods=["POST"])
async def openai_proxy(request: Request):
    method = request.method
    path = request.url.path
    target_url = f"{OCF_HEAD_URL}/v1/service/llm{path}"
    headers = dict(request.headers)
    headers.pop("host", None)
    content = await request.body()

    query_params = dict(request.query_params)

    async with aiohttp.ClientSession() as session:
        return await proxy_request(
            session=session,
            method=method,
            url=target_url,
            headers=headers,
            content=content if content else None,
            params=query_params if query_params else None
        )

@app.get("/health")
async def health_check():
    return {"status": "healthy", "service": "ocf-entry"}


@app.get("/v1/models")
async def list_models():
    try:
        models = await get_all_models(f"{OCF_HEAD_URL}/v1/dnt/table")
        return {
            "object": "list",
            "data": models
        }
    except Exception as e:
        raise HTTPException(status_code=502, detail=f"Failed to fetch models: {str(e)}")

# Include routers
app.include_router(auth_router)
app.include_router(models_router)

@app.get("/")
async def root():
    return {
        "message": "OpenAI Compatible Proxy Service",
        "version": "1.0.0",
        "endpoints": [
            "POST /v1/chat/completions",
            "POST /v1/completions",
            "POST /v1/embeddings",
            "GET /v1/models",
            "POST /api/auth/connect",
            "GET /api/user/profile",
            "POST /api/api-keys",
            "GET /api/api-keys",
            "DELETE /api/api-keys/{key_id}",
            "GET /api/models",
            "GET /api/models/public",
            "POST /api/v1/chat/completions",
            "POST /api/v1/completions"
        ]
    }

@app.on_event("startup")
async def startup_event():
    """Initialize database on startup"""
    init_db()

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)