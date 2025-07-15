from datetime import datetime

from app.schemas.types import IndonesianPhoneNumber
from pydantic import BaseModel, ConfigDict, EmailStr


class UserBase(BaseModel):
    full_name: str
    username: str
    phone_number: IndonesianPhoneNumber
    email: EmailStr


class UserCreate(UserBase):
    password_hash: str


class UserResponse(UserBase):
    id: int
    created_at: datetime
    updated_at: datetime | None = None

    model_config = ConfigDict(from_attributes=True)
