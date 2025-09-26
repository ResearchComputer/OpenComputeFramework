from fastapi import APIRouter, Depends, HTTPException, status, Request
from sqlalchemy.orm import Session
from .database import get_db, User, APIKey, APIUsageLog
from .auth import get_current_user, get_api_key_from_token
from .models import ModelResponse, ChatCompletionRequest, ErrorResponse
from .utils import get_all_models
import os
import aiohttp
import json

router = APIRouter()

OCF_HEAD_URL = os.getenv("OCF_HEAD_URL", "http://140.238.223.116:8092")

@router.get("/api/models")
async def get_available_models(
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    """Get available models (optional authentication for user-specific models/rate limits)"""
    try:
        # Get models from OCF
        ocf_models = await get_all_models(f"{OCF_HEAD_URL}/v1/dnt/table", with_details=True)

        # Transform to frontend format
        models = []
        for model in ocf_models:
            model_info = ModelResponse(
                id=model['id'],
                name=model['id'],
                description=f"AI model available on {model.get('device', 'unknown device')}",
                author=model.get('owner', 'unknown'),
                blockchain="ethereum",
                price="0.1",
                tags=["text-generation"],
                huggingFaceId=None
            )
            models.append(model_info)

        return models
    except Exception as e:
        raise HTTPException(status_code=502, detail=f"Failed to fetch models: {str(e)}")

@router.get("/api/models/public")
async def get_public_models():
    """Get available models without authentication"""
    try:
        # Get models from OCF
        ocf_models = await get_all_models(f"{OCF_HEAD_URL}/v1/dnt/table", with_details=True)

        # Transform to frontend format
        models = []
        for model in ocf_models:
            model_info = ModelResponse(
                id=model['id'],
                name=model['id'],
                description=f"AI model available on {model.get('device', 'unknown device')}",
                author=model.get('owner', 'unknown'),
                blockchain="ethereum",
                price="0.1",
                tags=["text-generation"],
                huggingFaceId=None
            )
            models.append(model_info)

        return models
    except Exception as e:
        raise HTTPException(status_code=502, detail=f"Failed to fetch models: {str(e)}")

async def proxy_request_with_logging(
    session: aiohttp.ClientSession,
    method: str,
    url: str,
    headers: dict,
    content: bytes = None,
    params: dict = None,
    api_key_id: str = None,
    model_id: str = None,
    db: Session = None
):
    """Proxy request with API usage logging"""
    try:
        async with session.request(
            method=method,
            url=url,
            headers=headers,
            data=content,
            params=params
        ) as response:
            response_data = await response.json()

            # Log API usage if we have API key info
            if api_key_id and model_id and db:
                # Extract token usage from response if available
                input_tokens = response_data.get('usage', {}).get('prompt_tokens', None)
                output_tokens = response_data.get('usage', {}).get('completion_tokens', None)

                usage_log = APIUsageLog(
                    api_key_id=api_key_id,
                    model_id=model_id,
                    input_tokens=input_tokens,
                    output_tokens=output_tokens,
                    cost=0.001  # Simple cost calculation
                )
                db.add(usage_log)
                db.commit()

            return response_data, response.status
    except Exception as e:
        raise HTTPException(status_code=502, detail=f"Request failed: {str(e)}")

@router.post("/api/v1/chat/completions")
async def chat_completions(
    request: Request,
    chat_request: ChatCompletionRequest,
    db: Session = Depends(get_db)
):
    """OpenAI-compatible chat completions endpoint with API key authentication"""
    # Get API key from Authorization header
    authorization = request.headers.get("Authorization")
    if not authorization:
        raise HTTPException(status_code=401, detail="Authorization header required")

    try:
        api_key_record = await get_api_key_from_token(authorization, db)
    except HTTPException as e:
        raise e

    # Update API key usage
    from datetime import datetime
    api_key_record.last_used = datetime.utcnow()
    db.commit()

    # Prepare proxy request
    target_url = f"{OCF_HEAD_URL}/v1/service/llm/v1/chat/completions"
    headers = dict(request.headers)
    headers.pop("host", None)
    headers.pop("authorization", None)  # Remove original auth header

    request_data = chat_request.dict()

    async with aiohttp.ClientSession() as session:
        try:
            response_data, status_code = await proxy_request_with_logging(
                session=session,
                method="POST",
                url=target_url,
                headers=headers,
                content=json.dumps(request_data).encode(),
                api_key_id=api_key_record.id,
                model_id=chat_request.model,
                db=db
            )
            return response_data
        except HTTPException as e:
            raise e
        except Exception as e:
            raise HTTPException(status_code=500, detail=f"Internal server error: {str(e)}")

@router.post("/api/v1/completions")
async def completions(
    request: Request,
    db: Session = Depends(get_db)
):
    """OpenAI-compatible completions endpoint with API key authentication"""
    # Get API key from Authorization header
    authorization = request.headers.get("Authorization")
    if not authorization:
        raise HTTPException(status_code=401, detail="Authorization header required")

    try:
        api_key_record = await get_api_key_from_token(authorization, db)
    except HTTPException as e:
        raise e

    # Update API key usage
    from datetime import datetime
    api_key_record.last_used = datetime.utcnow()
    db.commit()

    # Prepare proxy request
    target_url = f"{OCF_HEAD_URL}/v1/service/llm/v1/completions"
    headers = dict(request.headers)
    headers.pop("host", None)
    headers.pop("authorization", None)

    content = await request.body()

    async with aiohttp.ClientSession() as session:
        try:
            response_data, status_code = await proxy_request_with_logging(
                session=session,
                method="POST",
                url=target_url,
                headers=headers,
                content=content,
                api_key_id=api_key_record.id,
                model_id="unknown",  # Extract from request body if needed
                db=db
            )
            return response_data
        except HTTPException as e:
            raise e
        except Exception as e:
            raise HTTPException(status_code=500, detail=f"Internal server error: {str(e)}")