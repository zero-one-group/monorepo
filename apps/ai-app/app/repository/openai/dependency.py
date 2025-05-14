from typing import Annotated

from app.core.env import get_env
from app.repository.openai.greeting import GreetingRepoOpenAI
from fastapi import Depends
from openai import AsyncOpenAI


def get_openai_client() -> AsyncOpenAI:
    """Get OpenAI client with API key configured."""
    env = get_env()
    return AsyncOpenAI(api_key=env.OPENAI_API_KEY)


def get_greeting_repository_openai(
    openai: AsyncOpenAI = Depends(get_openai_client),
) -> GreetingRepoOpenAI:
    greeting_repo = GreetingRepoOpenAI(client=openai)
    return greeting_repo


DepGreetingRepoOpenAI = Annotated[
    GreetingRepoOpenAI, Depends(get_greeting_repository_openai)
]
