import logging
import sys

from app.core.env import get_env
from pythonjsonlogger import jsonlogger

# Singleton logger instance
_logger_instance = None


def get_logger(
    logger_name="ml_app",
    log_format="%(asctime)s %(name)s %(levelname)s %(message)s",
    stream=sys.stdout,
):
    global _logger_instance

    if _logger_instance is None:
        level = logging.INFO
        if get_env().DEBUG is True:
            level = logging.DEBUG

        logger = logging.getLogger(logger_name)
        logger.setLevel(level)

        if not logger.handlers:
            handler = logging.StreamHandler(stream)
            formatter = jsonlogger.JsonFormatter(log_format)
            handler.setFormatter(formatter)
            logger.addHandler(handler)

        _logger_instance = logger

    return _logger_instance


logger = get_logger()
