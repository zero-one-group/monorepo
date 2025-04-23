from functools import lru_cache

from pydantic_settings import BaseSettings, SettingsConfigDict


class Env(BaseSettings):
    """Application settings.

    Loads environment variables with pydantic-settings.
    """

    ML_PREFIX_API: str
    DEBUG: bool = False
    DATABASE_URL: str
    OPENAI_API_KEY: str

    model_config = SettingsConfigDict(
        env_file=".env",
        env_file_encoding="utf-8",
        case_sensitive=True,
    )


@lru_cache
def get_env() -> Env:
    """Get cached application settings.

    Uses lru_cache decorator to cache the settings instance.

    Returns:
        Settings: Application settings instance.
    """
    return Env()
