import logging
import os
import sys
from logging.handlers import RotatingFileHandler

from app.core.env import get_env
from pythonjsonlogger import jsonlogger

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

        # Set logging level based on DEBUG setting
        level = logging.DEBUG if get_env.DEBUG else logging.INFO

        logger = logging.getLogger(logger_name)
        logger.setLevel(level)

        if logger.handlers:
            logger.handlers.clear()

        console_handler = logging.StreamHandler(stream)
        console_handler.setLevel(level)
        console_formatter = jsonlogger.JsonFormatter(log_format)
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
            file_formatter = jsonlogger.JsonFormatter(log_format)
            file_handler.setFormatter(file_formatter)
            logger.addHandler(file_handler)

        if env.DEBUG:
            logger.debug("Debug logging enabled")

        _logger_instance = logger

    return _logger_instance


logger = get_logger()
