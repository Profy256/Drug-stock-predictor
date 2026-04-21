from fastapi import APIRouter, Depends, HTTPException, status
from sqlalchemy.orm import Session
from app.db import get_db
from app.models import (
    LoginRequest,
    LoginResponse,
    RegisterRequest,
    UserResponse,
    UserMe,
)
from app.models.database import User, Pharmacy, UserRole as DBUserRole
from app.core.security import hash_password, create_access_token, verify_password
from app.middleware.auth import get_current_user, get_client_ip
from app.services.audit import AuditService
from app.models import TokenData
from fastapi import Request
import logging

logger = logging.getLogger(__name__)

router = APIRouter(prefix="/api/v1/auth", tags=["auth"])


@router.post("/login", response_model=LoginResponse)
async def login(
    request: Request,
    credentials: LoginRequest,
    db: Session = Depends(get_db),
):
    """Login user and return JWT token"""
    audit_service = AuditService(db)
    client_ip = get_client_ip(request)
    
    # Find user by email
    user = db.query(User).filter(User.email == credentials.email).first()
    
    if not user:
        logger.warning(f"Login failed: user not found - {credentials.email}")
        audit_service.log_action(
            None, None, "LOGIN_FAILED", "USER", None, client_ip,
            {"reason": "user_not_found"}
        )
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Invalid credentials"
        )
    
    # Verify password
    if not verify_password(credentials.password, user.password_hash):
        logger.warning(f"Login failed: invalid password - {credentials.email}")
        audit_service.log_action(
            str(user.id), str(user.pharmacy_id), "LOGIN_FAILED", "USER",
            str(user.id), client_ip, {"reason": "invalid_password"}
        )
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Invalid credentials"
        )
    
    # Create token
    token = create_access_token(str(user.id), str(user.pharmacy_id), user.role.value)
    
    # Log successful login
    audit_service.log_action(
        str(user.id), str(user.pharmacy_id), "LOGIN_SUCCESS", "USER",
        str(user.id), client_ip
    )
    
    return LoginResponse(
        token=token,
        user_id=str(user.id),
        pharmacy_id=str(user.pharmacy_id),
        name=user.name,
        email=user.email,
        role=user.role,
    )


@router.post("/register", response_model=UserResponse)
async def register(
    request: Request,
    registration: RegisterRequest,
    db: Session = Depends(get_db),
):
    """Register a new user"""
    audit_service = AuditService(db)
    client_ip = get_client_ip(request)
    
    # Check if user already exists
    existing_user = db.query(User).filter(User.email == registration.email).first()
    if existing_user:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail="Email already registered"
        )
    
    # Check if pharmacy exists
    pharmacy = db.query(Pharmacy).filter(Pharmacy.id == registration.pharmacy_id).first()
    if not pharmacy:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail="Pharmacy not found"
        )
    
    # Create new user
    hashed_password = hash_password(registration.password)
    new_user = User(
        pharmacy_id=registration.pharmacy_id,
        name=registration.name,
        email=registration.email,
        password_hash=hashed_password,
        role=DBUserRole.data_entrant,
    )
    
    try:
        db.add(new_user)
        db.commit()
        db.refresh(new_user)
        
        # Log registration
        audit_service.log_action(
            str(new_user.id), str(new_user.pharmacy_id), "USER_REGISTERED",
            "USER", str(new_user.id), client_ip
        )
        
        return UserResponse(
            id=str(new_user.id),
            pharmacy_id=str(new_user.pharmacy_id),
            name=new_user.name,
            email=new_user.email,
            role=new_user.role,
            is_active=new_user.is_active,
            created_at=new_user.created_at,
            updated_at=new_user.updated_at,
        )
    except Exception as e:
        db.rollback()
        logger.error(f"Registration failed: {str(e)}")
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail="Registration failed"
        )


@router.get("/me", response_model=UserMe)
async def get_me(
    current_user: TokenData = Depends(get_current_user),
    db: Session = Depends(get_db),
):
    """Get current user information"""
    user = db.query(User).filter(User.id == current_user.user_id).first()
    
    if not user:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="User not found"
        )
    
    return UserMe(
        id=str(user.id),
        pharmacy_id=str(user.pharmacy_id),
        name=user.name,
        email=user.email,
        role=user.role,
    )
