from app.core.env import get_env
from app.core.logging import RequestIdMiddleware, logger
from app.router.main import router as system_router
from app.router.openai import router as openai_router
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

env = get_env()
app = FastAPI(
    title="Machine Learning API",
    description="Machine Learning API",
    version="0.1.0",
    docs_url="/docs",
    root_path=env.ML_PREFIX_API,
)
logger.debug(
    "Initializing Machine Learning app",
    extra={
        "docs_url": "/",
        "root_path": env.ML_PREFIX_API,
    },
)

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)
app.add_middleware(RequestIdMiddleware)

app.include_router(system_router)
app.include_router(openai_router)
