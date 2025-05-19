from app.core.logging import get_logger
from app.core.response import ErrorResponse, SuccessResponse, success_response
from app.repository.openai.dependency import DepGreetingRepoOpenAI


class GreetingService:
    __log = get_logger()

    def __init__(self, repo: DepGreetingRepoOpenAI):
        self.__repo = repo

    async def greetings(self) -> SuccessResponse[dict] | ErrorResponse:
        self.__log.info("Service layer log", extra={"layer": "service"})
        greetings = await self.__repo.greetings()
        return success_response(data=greetings)
