from typing import Annotated

from app.repository.openai.dependency import DepGreetingRepoOpenAI
from app.services.greeting import GreetingService
from fastapi import Depends


def get_greeting_service(repo: DepGreetingRepoOpenAI) -> GreetingService:
    # FastAPI will resolve the DepGreetingRepoOpenAI dependency first
    greeting_service = GreetingService(repo=repo)
    return greeting_service


DepGreetingService = Annotated[GreetingService, Depends(get_greeting_service)]
