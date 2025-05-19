from app.services.dependency import DepGreetingService
from fastapi import APIRouter

router = APIRouter(
    prefix="/openai",
    tags=["openai"],
)


@router.get("/greetings")
async def greetings(service: DepGreetingService):
    greetings = await service.greetings()
    return greetings
