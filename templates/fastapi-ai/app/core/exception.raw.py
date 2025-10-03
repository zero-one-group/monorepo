from typing import Any, Optional


class AppError(Exception):
    """
    Generic application error.

    :param message: human-readable message
    :param status_code: HTTP status to return
    :param code: machine-readable error code
    :param data: optional extra payload to include in the response
    """

    def __init__(
        self,
        message: str = "Operation failed",
        status_code: int = 400,
        code: str = "BAD_REQUEST",
        data: Optional[Any] = None,
    ):
        super().__init__(message)
        self.status_code = status_code
        self.code = code
        self.data = data
