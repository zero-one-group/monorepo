from app.core.env import get_env
from app.core.logging import DepLogger, logger
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


@app.get("/")
async def root(logger: DepLogger):
    logger.info("Incoming request at root path!", extra={"path": "/"})
    return {"message": "Welcome to the Machine Learning API"}
