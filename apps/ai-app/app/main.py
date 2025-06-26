from contextlib import asynccontextmanager

from app.core.database import Database
from app.core.env import get_env
from app.core.exception import AppError
from app.core.logging import RequestIdMiddleware, logger
from app.core.trace import init_tracer, instrument_app
from app.router.openai import router as openai_router
from app.router.root import router as root_router
from fastapi import FastAPI, Request
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import JSONResponse


@asynccontextmanager
async def lifespan(app: FastAPI):
    # see: https://fastapi.tiangolo.com/advanced/events/#async-context-manager
    # On startup hook
    logger.debug(
        "Initializing Machine Learning app",
        extra={
            "docs_url": "/",
            "root_path": env.ML_PREFIX_API,
        },
    )

    yield

    # On shutdown hook
    logger.info("Application shutting down, preparing for graceful shutdown")
    logger.info("Disposing database connections")
    await Database.dispose()


env = get_env()
app = FastAPI(
    title="Machine Learning API",
    description="Machine Learning API",
    version="0.1.0",
    docs_url="/docs",
    root_path=env.ML_PREFIX_API,
    lifespan=lifespan,
)

if env.APP_ENVIRONMENT == "production":
    logger.info(
        "The environment is set to production; instrumentation is being configured."
    )
    init_tracer(
        service_name=env.APP_NAME,
        otlp_endpoint=env.OTEL_EXPORTER_OTLP_ENDPOINT,
    )
    instrument_app(app)

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)
app.add_middleware(RequestIdMiddleware)

app.include_router(root_router)
app.include_router(openai_router)


@app.exception_handler(AppError)
async def app_error_handler(request: Request, exc: AppError):
    message = str(exc)
    logger.error(msg=message, exc_info=exc)

    return JSONResponse(
        status_code=exc.status_code,
        content={
            "success": False,
            "message": message,
            "error_code": exc.code,
            "data": exc.data,
        },
    )
