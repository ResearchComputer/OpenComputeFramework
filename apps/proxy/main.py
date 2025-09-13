from fastapi import FastAPI, Request, HTTPException
from fastapi.responses import JSONResponse
import httpx
import os
from typing import Optional, Dict, Any
import asyncio
import json

app = FastAPI(title="OpenAI Compatible Proxy Service")

TARGET_SERVICE_URL = os.getenv("TARGET_SERVICE_URL", "http://localhost:8000/v1")
TIMEOUT = float(os.getenv("TIMEOUT", "30.0"))

async def proxy_request(
    client: httpx.AsyncClient,
    method: str,
    url: str,
    headers: Dict[str, str],
    content: Optional[bytes] = None,
    params: Optional[Dict[str, Any]] = None
) -> JSONResponse:
    try:
        response = await client.request(
            method=method,
            url=url,
            headers=headers,
            content=content,
            params=params,
            timeout=TIMEOUT
        )

        return JSONResponse(
            status_code=response.status_code,
            content=response.json(),
            headers=dict(response.headers)
        )
    except httpx.TimeoutException:
        raise HTTPException(status_code=504, detail="Request timeout")
    except httpx.RequestError as e:
        raise HTTPException(status_code=502, detail=f"Request failed: {str(e)}")
    except json.JSONDecodeError:
        return JSONResponse(
            status_code=response.status_code,
            content={"error": "Invalid JSON response"},
            headers=dict(response.headers)
        )

@app.api_route("/v1/chat/completions", methods=["POST"])
@app.api_route("/v1/completions", methods=["POST"])
@app.api_route("/v1/embeddings", methods=["POST"])
@app.api_route("/v1/models", methods=["GET"])
async def openai_proxy(request: Request):
    method = request.method
    path = request.url.path

    target_url = f"{TARGET_SERVICE_URL}{path}"

    headers = dict(request.headers)
    headers.pop("host", None)

    content = await request.body()

    query_params = dict(request.query_params)

    async with httpx.AsyncClient() as client:
        return await proxy_request(
            client=client,
            method=method,
            url=target_url,
            headers=headers,
            content=content if content else None,
            params=query_params if query_params else None
        )

@app.get("/health")
async def health_check():
    return {"status": "healthy", "service": "openai-proxy"}

@app.get("/")
async def root():
    return {
        "message": "OpenAI Compatible Proxy Service",
        "version": "1.0.0",
        "endpoints": [
            "POST /v1/chat/completions",
            "POST /v1/completions",
            "POST /v1/embeddings",
            "GET /v1/models"
        ]
    }

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)