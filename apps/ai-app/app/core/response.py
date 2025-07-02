from typing import Generic, Optional, TypeVar

from fastapi import status
from fastapi.responses import JSONResponse
from pydantic import BaseModel, ConfigDict

T = TypeVar("T")


class ResponseModel(BaseModel, Generic[T]):
    success: bool
    message: str
    data: Optional[T] = None
    metadata: Optional[dict] = None
    error_code: Optional[str] = None


class SuccessResponse(ResponseModel[T], Generic[T]):
    success: bool = True
    error_code: None = None

    model_config = ConfigDict(
        json_schema_extra={
            "example": {
                "success": True,
                "message": "Operation completed successfully",
                "data": {},
                "error_code": None,
            }
        }
    )


class ErrorResponse(ResponseModel[None]):
    success: bool = False
    data: None = None

    model_config = ConfigDict(
        json_schema_extra={
            "example": {
                "success": False,
                "message": "An error occurred",
                "data": None,
                "error_code": "INTERNAL_ERROR",
            }
        }
    )


def success_response(
    data: T,
    message: str = "Operation successful",
    status_code: int = status.HTTP_200_OK,
) -> JSONResponse:
    """
    Wraps a SuccessResponse in a JSONResponse so FastAPI
    returns the proper HTTP status code.
    """
    payload = SuccessResponse(message=message, data=data).model_dump(exclude_none=True)
    return JSONResponse(status_code=status_code, content=payload)
