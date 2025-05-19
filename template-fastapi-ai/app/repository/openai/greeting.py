from app.core.exception import AppError
from app.core.logging import get_logger
from openai import AsyncOpenAI


class GreetingRepoOpenAI:
    __log = get_logger()

    def __init__(self, client: AsyncOpenAI):
        self.client = client

    async def greetings(self):
        self.__log.info(
            "Fetching greetings in 5 languages from OpenAI",
            extra={"layer": "repository"},
        )
        try:
            response = await self.client.chat.completions.create(
                model="gpt-4.1-mini",
                messages=[
                    {
                        "role": "system",
                        "content": "You are a helpful assistant that provides greetings in different languages. Always respond in valid JSON format.",
                    },
                    {
                        "role": "user",
                        "content": 'Give me greetings in exactly 5 different languages. Return only a JSON object with this exact structure: {"greetings": [{"language": "name of language", "greeting": "greeting in that language"}]}. Include no explanations or additional text.',
                    },
                ],
                temperature=0.7,
                max_tokens=150,
                response_format={"type": "json_object"},
            )

            content = response.choices[0].message.content

            self.__log.info(
                "Successfully received greetings from OpenAI",
                extra={"layer": "repository"},
            )

        except Exception as e:
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
        except json.JSONDecodeError:
            raise AppError(
                "Invalid JSON from OpenAI",
                status_code=502,
                code="PARSE_ERROR",
            )
