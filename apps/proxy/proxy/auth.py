from datetime import datetime, timedelta
from typing import Optional
from jose import JWTError, jwt
from passlib.context import CryptContext
from fastapi import HTTPException, status, Depends
from fastapi.security import HTTPBearer, HTTPAuthorizationCredentials
from sqlalchemy.orm import Session
import uuid
import secrets
from .database import get_db, User, APIKey

SECRET_KEY = "your-secret-key-here-change-in-production"
ALGORITHM = "HS256"
ACCESS_TOKEN_EXPIRE_MINUTES = 24 * 60  # 24 hours

pwd_context = CryptContext(schemes=["bcrypt"], deprecated="auto")
security = HTTPBearer()

def verify_address_signature(address: str, signature: str, chain_id: int) -> bool:
    """Verify Ethereum signature for wallet authentication"""
    try:
        from web3 import Web3
        # This is a simplified signature verification
        # In production, you should implement proper EIP-191/EIP-712 signature verification
        message = f"Welcome to OCF! Sign this message to authenticate your wallet.\n\nChain ID: {chain_id}"

        # Hash the message
        message_hash = Web3.solidityKeccak(
            ['string'],
            [message]
        )

        # Recover the address from signature
        recovered_address = Web3.eth.account.recoverHash(message_hash, signature=signature)

        return recovered_address.lower() == address.lower()
    except Exception as e:
        return False

def create_access_token(data: dict, expires_delta: Optional[timedelta] = None):
    to_encode = data.copy()
    if expires_delta:
        expire = datetime.utcnow() + expires_delta
    else:
        expire = datetime.utcnow() + timedelta(minutes=15)
    to_encode.update({"exp": expire})
    encoded_jwt = jwt.encode(to_encode, SECRET_KEY, algorithm=ALGORITHM)
    return encoded_jwt

def verify_token(token: str) -> dict:
    try:
        payload = jwt.decode(token, SECRET_KEY, algorithms=[ALGORITHM])
        user_id: str = payload.get("sub")
        if user_id is None:
            raise HTTPException(
                status_code=status.HTTP_401_UNAUTHORIZED,
                detail="Invalid authentication credentials",
                headers={"WWW-Authenticate": "Bearer"},
            )
        return {"user_id": user_id}
    except JWTError:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Invalid authentication credentials",
            headers={"WWW-Authenticate": "Bearer"},
        )

async def get_current_user(
    credentials: HTTPAuthorizationCredentials = Depends(security),
    db: Session = Depends(get_db)
) -> User:
    token = credentials.credentials
    payload = verify_token(token)
    user_id = payload["user_id"]

    user = db.query(User).filter(User.id == user_id).first()
    if user is None:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="User not found",
            headers={"WWW-Authenticate": "Bearer"},
        )
    return user

def generate_api_key():
    """Generate a secure API key"""
    prefix = "sk_rc_"
    key_part = secrets.token_urlsafe(32)
    return f"{prefix}{key_part}"

def hash_api_key(api_key: str) -> str:
    """Hash API key for storage"""
    return pwd_context.hash(api_key)

def verify_api_key(api_key: str, hashed_key: str) -> bool:
    """Verify API key against hash"""
    return pwd_context.verify(api_key, hashed_key)

async def get_api_key_from_token(authorization: str, db: Session) -> APIKey:
    """Extract and verify API key from Authorization header"""
    if not authorization.startswith("Bearer "):
        raise HTTPException(status_code=401, detail="Invalid authorization format")

    api_key = authorization[7:]  # Remove "Bearer " prefix

    # Find API key by hash
    api_key_record = db.query(APIKey).filter(APIKey.key_hash == hash_api_key(api_key)).first()

    if not api_key_record or not api_key_record.is_active:
        raise HTTPException(status_code=401, detail="Invalid or inactive API key")

    return api_key_record