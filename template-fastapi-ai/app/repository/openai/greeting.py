from app.core.logging import get_logger


class GreetingRepoOpenAI:
    __log = get_logger()

    def greetings(self):
        self.__log.info(
            "Hello world from root repo openai", extra={"layer": "repository"}
        )
        return {"response": "Hello World"}
