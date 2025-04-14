from app.core.database import DepDB, check_db_connection, engine
from app.core.env import get_env
from app.core.logging import DepLogger, RequestIdMiddleware, logger
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


@app.get("/")
async def root(logger: DepLogger):
    logger.info("Incoming request at root path!", extra={"path": "/"})
    return {"message": "Welcome to the Machine Learning API"}


@app.get("/health-check")
async def health_check(logger: DepLogger, db: DepDB):
    """
    Checks the health of the application and related connection.
    """
    logger.info("Performing health check...")
    if check_db_connection(db):
        logger.info("Health check successful: Database connection verified.")
        return {"status": "ok", "database": "connected"}
    else:
        logger.error("Health check failed: Database connection error.")
        raise HTTPException(
            status_code=503,
            detail="Database connection error",
        )


@app.on_event("shutdown")
def shutdown():
    logger.info("Application shutting down, preparing for graceful shutdown")
    logger.info("Disposing database connections")
    engine.dispose()
    logger.info("Database engine disposed")
