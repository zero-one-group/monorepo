import asyncio
from collections.abc import AsyncGenerator
from contextlib import asynccontextmanager
from typing import Annotated, Optional

from app.core.env import get_env
from app.core.logging import logger
from fastapi import Depends
from sqlalchemy import text
from sqlalchemy.exc import OperationalError, SQLAlchemyError
from sqlalchemy.ext.asyncio import AsyncSession, async_sessionmaker, create_async_engine
from sqlalchemy.orm import DeclarativeBase


class Database:
    _instance: Optional["Database"] = None
    _lock = asyncio.Lock()
    _initialized = False
    engine = None
    async_session_maker = None

    def __new__(cls):
        if cls._instance is None:
            cls._instance = super(Database, cls).__new__(cls)
        return cls._instance

    @classmethod
    async def initialize(cls):
        async with cls._lock:
            if cls._initialized:
                return

            env = get_env()
            db_url = env.DATABASE_URL
            if db_url.startswith("postgresql://"):
                db_url = db_url.replace(
                    "postgresql://", "postgresql+asyncpg://")

            is_development = env.APP_ENVIRONMENT == "development"
            cls.engine = create_async_engine(
                db_url, echo=is_development, pool_pre_ping=True
            )

            cls.async_session_maker = async_sessionmaker(
                autocommit=False,
                autoflush=False,
                bind=cls.engine,
                expire_on_commit=False,
            )

            cls._initialized = True
            logger.debug("Database initialized successfully")

    @classmethod
    @asynccontextmanager
    async def get_session(cls) -> AsyncGenerator[AsyncSession, None]:
        if not cls._initialized:
            await cls.initialize()

        session = cls.async_session_maker()
        try:
            yield session
        finally:
            await session.close()

    @classmethod
    async def check_connection(cls) -> bool:
        if not cls._initialized:
            await cls.initialize()

        async with cls.get_session() as session:
            try:
                await session.execute(text("SELECT 1"))
                logger.debug("Database connection check successful.")
                return True
            except OperationalError as e:
                logger.error(
                    f"Database connection operational error: {e}",
                    exc_info=False,
                )
                return False
            except SQLAlchemyError as e:
                logger.error(
                    f"Database query execution error during health check: {e}",
                    exc_info=True,
                )
                return False
            except Exception as e:
                logger.error(
                    f"Unexpected error during database connection check: {e}",
                    exc_info=True,
                )
                return False

    @classmethod
    async def dispose(cls):
        if cls._initialized and cls.engine:
            await cls.engine.dispose()
            logger.info("Database engine disposed")
            cls._initialized = False


# NOTE: Base class that will be inherited across model
class Base(DeclarativeBase):
    pass


# Initialize database and get session
async def get_db() -> AsyncGenerator[AsyncSession, None]:
    async with Database.get_session() as session:
        yield session


# NOTE: Type alias for database session dependency
DepDB = Annotated[AsyncSession, Depends(get_db)]


# Expose connection check function for convenience
async def check_db_connection() -> bool:
    return await Database.check_connection()
