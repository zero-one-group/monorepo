from app.core.logging import get_logger
from app.core.response import ErrorResponse, SuccessResponse, success_response
from app.core.trace import get_tracer
from app.repository.openai.dependency import DepGreetingRepoOpenAI

tracer = get_tracer("service.greetings")


class GreetingService:
    __log = get_logger()

    def __init__(self, repo: DepGreetingRepoOpenAI):
        self.__repo = repo

    async def greetings(self) -> SuccessResponse[dict] | ErrorResponse:
        with tracer.start_as_current_span("service.greetings") as span:
            self.__log.info("Service layer log", extra={"layer": "service"})
            greetings = await self.__repo.greetings()
            return success_response(data=greetings)
