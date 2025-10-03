from functools import lru_cache

from pydantic_settings import BaseSettings, SettingsConfigDict


class Env(BaseSettings):
    """Application settings.

    Loads environment variables with pydantic-settings.
    """

    ML_PREFIX_API: str
    APP_NAME: str = "{{ package_name | kebab_case }}"
    APP_ENVIRONMENT: str = "development"
    DATABASE_URL: str
    OPENAI_API_KEY: str
    OTEL_EXPORTER_OTLP_ENDPOINT: str

    model_config = SettingsConfigDict(
        env_file=".env", env_file_encoding="utf-8", case_sensitive=True, extra="ignore"
    )


@lru_cache
def get_env() -> Env:
    """Get cached application settings.

    Uses lru_cache decorator to cache the settings instance.

    Returns:
        Settings: Application settings instance.
    """
    return Env()
