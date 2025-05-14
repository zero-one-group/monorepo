import logging
import os
import sys
import uuid
from logging.handlers import RotatingFileHandler
from typing import Annotated

from app.core.env import get_env
from fastapi import Depends, Request
from pythonjsonlogger.jsonlogger import JsonFormatter
from starlette.middleware.base import BaseHTTPMiddleware

# Singleton logger instance
_logger_instance = None


def get_logger(
    logger_name="ml_app",
    log_format="%(asctime)s %(name)s %(levelname)s %(message)s",
    stream=sys.stdout,
    log_to_file=True,
    log_dir="logs",
    log_file="app.log",
    max_file_size=10 * 1024 * 1024,  # 10 MB
    backup_count=5,
):
    global _logger_instance

    if _logger_instance is None:
        env = get_env()

        is_development = env.APP_ENVIRONMENT == "development"
        # Set logging level based on DEBUG setting
        level = logging.DEBUG if is_development else logging.INFO

        logger = logging.getLogger(logger_name)
        logger.setLevel(level)

        if logger.handlers:
            logger.handlers.clear()

        console_handler = logging.StreamHandler(stream)
        console_handler.setLevel(level)

        if is_development == False:
            console_formatter = JsonFormatter(log_format)
            console_handler.setFormatter(console_formatter)

        logger.addHandler(console_handler)

        # Add file handler if log_to_file is True
        if log_to_file:
            os.makedirs(log_dir, exist_ok=True)

            log_file_path = os.path.join(log_dir, log_file)
            file_handler = RotatingFileHandler(
                log_file_path, maxBytes=max_file_size, backupCount=backup_count
            )
            file_handler.setLevel(level)
            file_formatter = JsonFormatter(log_format)
            file_handler.setFormatter(file_formatter)
            logger.addHandler(file_handler)

        if is_development:
            logger.debug("Debug logging enabled")

        _logger_instance = logger

    return _logger_instance


logger = get_logger()


class RequestIdMiddleware(BaseHTTPMiddleware):
    async def dispatch(self, request: Request, call_next):
        request_id = str(uuid.uuid4())

        # Add request_id to the request state
        request.state.request_id = request_id
        response = await call_next(request)
        return response


def get_request_id(request: Request) -> str:
    """Get the request ID from the current request state"""
    return getattr(request.state, "request_id", str(uuid.uuid4()))


def get_logger_dependency(request: Request):
    """Dependency to get the logger instance with request_id context."""
    logger = get_logger()
    request_id = get_request_id(request)
    # Return a logger adapter with the request_id in the extra dict
    return logging.LoggerAdapter(logger, {"request_id": request_id})


# Type alias
DepLogger = Annotated[logging.LoggerAdapter, Depends(get_logger_dependency)]
