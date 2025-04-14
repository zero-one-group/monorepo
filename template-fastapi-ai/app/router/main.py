from app.core.database import DepDB, check_db_connection
from app.core.logging import DepLogger
from fastapi import APIRouter, HTTPException

router = APIRouter()


@router.get("/")
async def root(logger: DepLogger):
    logger.info("Incoming request at root path!", extra={"path": "/"})
    return {"message": "Welcome to the Machine Learning API"}


@router.get("/health-check")
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
