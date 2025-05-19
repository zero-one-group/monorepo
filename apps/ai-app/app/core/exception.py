class AppError(Exception):
    """
    Generic application error.
    :param message: human‐readable message
    :param status_code: HTTP status to return
    :param code: machine‐readable error code
    """

    def __init__(
        self,
        message: str = "Operation failed",
        status_code: int = 400,
        code: str = "BAD_REQUEST",
    ):
        super().__init__(message)
        self.status_code = status_code
        self.code = code
