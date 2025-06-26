from app.core.exception import AppError
from app.core.logging import get_logger
from app.core.trace import get_tracer
from openai import AsyncOpenAI

tracer = get_tracer("repository.greetings")


class GreetingRepoOpenAI:
    __log = get_logger()

    def __init__(self, client: AsyncOpenAI):
        self.client = client
        self.model = "gpt-4.1-mini"
        self.temperature = 0.7

    async def greetings(self):
        self.__log.info(
            "Fetching greetings in 5 languages from OpenAI",
            extra={"layer": "repository"},
        )
        messages = [
            {
                "role": "system",
                "content": "You are a helpful assistant that provides greetings in different languages. Always respond in valid JSON format.",
            },
            {
                "role": "user",
                "content": 'Give me greetings in exactly 5 different languages. Return only a JSON object with this exact structure: {"greetings": [{"language": "name of language", "greeting": "greeting in that language"}]}. Include no explanations or additional text.',
            },
        ]
        with tracer.start_as_current_span("repository.fetch_greetings") as span:
            span.set_attribute("llm.req.model", self.model)
            span.set_attribute("llm.req.temperature", self.temperature)
            span.set_attribute(
                "llm.req.messages",
                [m["content"] for m in messages],
            )
            try:
                response = await self.client.chat.completions.create(
                    model=self.model,
                    messages=messages,
                    temperature=self.temperature,
                    max_tokens=150,
                    response_format={"type": "json_object"},
                )

                content = response.choices[0].message.content
                span.set_attribute("llm.resp.model", response.model)

                usage = response.usage
                span.set_attribute("llm.usage.prompt_tokens",
                                   usage.prompt_tokens)
                span.set_attribute(
                    "llm.usage.completion_tokens", usage.completion_tokens
                )
                span.set_attribute("llm.usage.total_tokens",
                                   usage.total_tokens)

                self.__log.info(
                    "Successfully received greetings from OpenAI",
                    extra={"layer": "repository"},
                )

            except Exception as e:
                span.record_exception(e)
                raise AppError(
                    f"OpenAI unavailable: {str(e)}",
                    status_code=503,
                    code="OPENAI_ERROR",
                )

            try:
                import json

                greetings_data = json.loads(content)

                self.__log.info(
                    f"Successfully parsed greetings in {
                        len(greetings_data.get('greetings', []))} languages",
                    extra={"layer": "repository"},
                )

                return {"response": greetings_data}
            except json.JSONDecodeError as e:
                span.record_exception(e)
                raise AppError(
                    "Invalid JSON from OpenAI",
                    status_code=502,
                    code="PARSE_ERROR",
                )
