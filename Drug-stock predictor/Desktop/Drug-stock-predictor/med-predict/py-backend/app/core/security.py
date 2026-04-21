from passlib.context import CryptContext
import jwt
from datetime import datetime, timedelta, timezone
from app.core.config import settings
from app.models import TokenData, UserRole

pwd_context = CryptContext(schemes=["bcrypt"], deprecated="auto")


def hash_password(password: str) -> str:
    """Hash a password"""
    return pwd_context.hash(password)


def verify_password(plain_password: str, hashed_password: str) -> bool:
    """Verify a password"""
    return pwd_context.verify(plain_password, hashed_password)


def create_access_token(user_id: str, pharmacy_id: str, role: str, expires_delta: timedelta = None) -> str:
    """Create a JWT access token"""
    if expires_delta:
        expire = datetime.now(timezone.utc) + expires_delta
    else:
        expire = datetime.now(timezone.utc) + timedelta(hours=settings.JWT_EXPIRATION_HOURS)
    
    to_encode = {
        "sub": user_id,
        "pharmacy_id": pharmacy_id,
        "role": role,
        "exp": expire,
    }
    
    encoded_jwt = jwt.encode(
        to_encode,
        settings.JWT_SECRET,
        algorithm=settings.JWT_ALGORITHM,
    )
    return encoded_jwt


def decode_token(token: str) -> TokenData:
    """Decode a JWT token"""
    try:
        payload = jwt.decode(
            token,
            settings.JWT_SECRET,
            algorithms=[settings.JWT_ALGORITHM],
        )
        user_id: str = payload.get("sub")
        pharmacy_id: str = payload.get("pharmacy_id")
        role: str = payload.get("role")
        
        if user_id is None:
            raise jwt.InvalidTokenError("Invalid token")
        
        return TokenData(user_id=user_id, pharmacy_id=pharmacy_id, role=role)
    except jwt.InvalidTokenError:
        raise
