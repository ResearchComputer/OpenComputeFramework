from fastapi import APIRouter, Depends, HTTPException, status
from sqlalchemy.orm import Session
from .database import get_db, User, APIKey, APIUsageLog
from .auth import verify_address_signature, create_access_token, generate_api_key, hash_api_key, get_current_user
from .models import AuthRequest, AuthResponse, UserResponse, APIKeyCreate, APIKeyResponse, APIKeyInfo, ErrorResponse

router = APIRouter()

@router.post("/api/auth/connect", response_model=AuthResponse)
async def connect_wallet(request: AuthRequest, db: Session = Depends(get_db)):
    """Connect wallet and create JWT token"""
    # Verify signature
    if not verify_address_signature(request.address, request.signature, request.chain_id):
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Invalid signature",
        )

    # Find or create user
    user = db.query(User).filter(User.address == request.address.lower()).first()
    if not user:
        user = User(address=request.address.lower())
        db.add(user)
        db.commit()
        db.refresh(user)

    # Create JWT token
    access_token = create_access_token(data={"sub": user.id})

    return {
        "token": access_token,
        "user": UserResponse(id=user.id, address=user.address, created_at=user.created_at)
    }

@router.get("/api/user/profile", response_model=UserResponse)
async def get_user_profile(current_user: User = Depends(get_current_user), db: Session = Depends(get_db)):
    """Get current user profile"""
    return UserResponse(id=current_user.id, address=current_user.address, created_at=current_user.created_at)

@router.post("/api/api-keys", response_model=APIKeyResponse)
async def create_api_key(
    key_data: APIKeyCreate,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    """Create new API key for user"""
    # Generate API key
    api_key_value = generate_api_key()
    api_key_hash = hash_api_key(api_key_value)

    # Store API key
    api_key = APIKey(
        user_id=current_user.id,
        name=key_data.name,
        key_hash=api_key_hash
    )
    db.add(api_key)
    db.commit()
    db.refresh(api_key)

    return APIKeyResponse(
        id=api_key.id,
        name=api_key.name,
        key=api_key_value,  # Return the actual key only on creation
        created_at=api_key.created_at,
        last_used=api_key.last_used,
        is_active=api_key.is_active
    )

@router.get("/api/api-keys", response_model=dict)
async def list_api_keys(
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    """List user's API keys"""
    api_keys = db.query(APIKey).filter(APIKey.user_id == current_user.id).all()

    key_list = []
    for key in api_keys:
        key_info = APIKeyInfo(
            id=key.id,
            name=key.name,
            created_at=key.created_at,
            last_used=key.last_used,
            is_active=key.is_active
        )
        key_list.append(key_info)

    return {"apiKeys": key_list}

@router.delete("/api/api-keys/{key_id}")
async def delete_api_key(
    key_id: str,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    """Delete (deactivate) API key"""
    api_key = db.query(APIKey).filter(
        APIKey.id == key_id,
        APIKey.user_id == current_user.id
    ).first()

    if not api_key:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="API key not found"
        )

    api_key.is_active = False
    db.commit()

    return {"success": True}

@router.patch("/api/api-keys/{key_id}/usage")
async def update_api_key_usage(
    key_id: str,
    db: Session = Depends(get_db)
):
    """Update API key last used timestamp (called when API key is used)"""
    api_key = db.query(APIKey).filter(APIKey.id == key_id).first()

    if api_key:
        from datetime import datetime
        api_key.last_used = datetime.utcnow()
        db.commit()

    return {"success": True}