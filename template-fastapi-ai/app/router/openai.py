from app.core.response import ErrorResponse, SuccessResponse
from app.core.trace import get_tracer
from app.services.dependency import DepGreetingService
from fastapi import APIRouter

router = APIRouter(
    prefix="/openai",
    tags=["openai"],
)

tracer = get_tracer("api.greetings")


@router.get(
    "/greetings",
    responses={
        200: {"model": SuccessResponse[dict]},
        400: {"model": ErrorResponse},
        502: {"model": ErrorResponse},
        503: {"model": ErrorResponse},
    },
)
async def greetings(service: DepGreetingService):
    with tracer.start_as_current_span("route.greetings") as span:
        span.set_attribute("http.method", "GET")
        span.set_attribute("http.route", "/greetings")
        greetings = await service.greetings()
        return greetings
