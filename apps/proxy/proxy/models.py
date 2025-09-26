from pydantic import BaseModel
from typing import Optional, List
from datetime import datetime

class UserCreate(BaseModel):
    address: str
    signature: str
    chain_id: int

class UserResponse(BaseModel):
    id: str
    address: str
    created_at: datetime

    class Config:
        from_attributes = True

class APIKeyCreate(BaseModel):
    name: str

class APIKeyResponse(BaseModel):
    id: str
    name: str
    key: str
    created_at: datetime
    last_used: Optional[datetime]
    is_active: bool

    class Config:
        from_attributes = True

class APIKeyInfo(BaseModel):
    id: str
    name: str
    created_at: datetime
    last_used: Optional[datetime]
    is_active: bool

    class Config:
        from_attributes = True

class ErrorResponse(BaseModel):
    error: dict

class AuthRequest(BaseModel):
    address: str
    signature: str
    chain_id: int

class AuthResponse(BaseModel):
    token: str
    user: UserResponse

class ModelResponse(BaseModel):
    id: str
    name: str
    description: Optional[str]
    author: Optional[str]
    blockchain: Optional[str]
    price: Optional[str]
    tags: Optional[List[str]]
    huggingFaceId: Optional[str]

class ChatCompletionRequest(BaseModel):
    model: str
    messages: List[dict]
    temperature: Optional[float] = 0.7
    max_tokens: Optional[int] = 2048
    stream: Optional[bool] = False

class APIUsage(BaseModel):
    id: str
    api_key_id: str
    model_id: str
    input_tokens: Optional[int]
    output_tokens: Optional[int]
    timestamp: datetime
    cost: Optional[float]

    class Config:
        from_attributes = True