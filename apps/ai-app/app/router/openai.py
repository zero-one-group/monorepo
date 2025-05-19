from app.core.response import ErrorResponse, SuccessResponse
from app.services.dependency import DepGreetingService
from fastapi import APIRouter

router = APIRouter(
    prefix="/openai",
    tags=["openai"],
)


@router.get(
    "/greetings",
    responses={
        200: {"model": SuccessResponse[dict]},
        400: {"model": ErrorResponse},
        502: {"model": ErrorResponse},
        503: {"model": ErrorResponse},
    },
)
async def greetings(service: DepGreetingService):
    greetings = await service.greetings()
    return greetings
