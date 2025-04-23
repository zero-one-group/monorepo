from collections.abc import Generator
from typing import Annotated

from app.core.env import get_env
from app.core.logging import logger
from fastapi import Depends
from sqlalchemy import create_engine, text
from sqlalchemy.exc import OperationalError, SQLAlchemyError
from sqlalchemy.orm import DeclarativeBase, Session, sessionmaker

env = get_env()

engine = create_engine(env.DATABASE_URL, echo=env.DEBUG)

SessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)


# NOTE: Base class that will be inherited across model
class Base(DeclarativeBase):
    pass


def get_db() -> Generator[Session, None, None]:
    db = SessionLocal()
    try:
        yield db
    finally:
        # closing and returning to connection pool
        # ref: https://stackoverflow.com/a/8705750
        db.close()


# NOTE: Type alias for database session dependency
DepDB = Annotated[Session, Depends(get_db)]


def check_db_connection(db: Session) -> bool:
    try:
        db.execute(text("SELECT 1"))
        logger.debug("Database connection check successful.")
        return True
    except OperationalError as e:
        logger.error(
            f"Database connection operational error: {
                e}",
            exc_info=False,
        )
        return False
    except SQLAlchemyError as e:
        logger.error(
            f"Database query execution error during health check: {e}", exc_info=True
        )
        return False
    except Exception as e:
        logger.error(
            f"Unexpected error during database connection check: {e}", exc_info=True
        )
        return False
