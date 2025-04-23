from typing import Annotated

from app.repository.openai.greeting import GreetingRepoOpenAI
from fastapi import Depends


def get_greeting_repository_openai() -> GreetingRepoOpenAI:
    greeting_repo = GreetingRepoOpenAI()
    return greeting_repo


DepGreetingRepoOpenAI = Annotated[
    GreetingRepoOpenAI, Depends(get_greeting_repository_openai)
]
