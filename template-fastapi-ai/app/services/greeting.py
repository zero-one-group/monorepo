from app.core.logging import get_logger
from app.repository.openai.dependency import DepGreetingRepoOpenAI


class GreetingService:
    __log = get_logger()

    def __init__(self, repo: DepGreetingRepoOpenAI):
        self.__repo = repo

    def greetings(self) -> dict:
        self.__log.info("Service layer log", extra={"layer": "service"})
        return self.__repo.greetings()
