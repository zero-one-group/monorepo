from app.core.env import get_env
from app.core.logging import logger
from opentelemetry import trace
from opentelemetry.exporter.otlp.proto.grpc.trace_exporter import OTLPSpanExporter
from opentelemetry.instrumentation.fastapi import FastAPIInstrumentor
from opentelemetry.sdk.resources import SERVICE_NAME, Resource
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.trace.export import BatchSpanProcessor
from prometheus_fastapi_instrumentator import Instrumentator


def init_tracer(
    service_name: str,
    otlp_endpoint: str = "localhost:4317",
) -> None:
    """
    Configure a TracerProvider with an OTLP gRPC exporter
    sending to otlp_endpoint (default localhost:4317).
    """
    resource = Resource.create({SERVICE_NAME: service_name})
    provider = TracerProvider(resource=resource)
    otlp_exporter = OTLPSpanExporter(
        endpoint=otlp_endpoint,
        insecure=True,
    )
    provider.add_span_processor(BatchSpanProcessor(otlp_exporter))
    trace.set_tracer_provider(provider)


def instrument_app(app) -> None:
    env = get_env()
    if env.APP_ENVIRONMENT != "production":
        return

    logger.info(
        "The environment is set to production; instrumentation is being configured."
    )

    # Tracing
    init_tracer(
        service_name=env.APP_NAME,
        otlp_endpoint=env.OTEL_EXPORTER_OTLP_ENDPOINT,
    )

    # @see: https://opentelemetry-python-contrib.readthedocs.io/en/latest/instrumentation/fastapi/fastapi.html
    FastAPIInstrumentor.instrument_app(app, exclude_spans=["receive", "send"])

    # Metrics
    # @see: https://github.com/trallnag/prometheus-fastapi-instrumentator?tab=readme-ov-file#advanced-usage
    Instrumentator(
        should_group_status_codes=True,  # e.g. 2xx, 4xx, 5xx
        should_ignore_untemplated=True,  # drop paths with path parameters
        should_respect_env_var=False,  # always on
    ).instrument(app).expose(
        app,
        include_in_schema=False,  # hide from your OpenAPI docs
        endpoint="/metrics",  # where Prometheus will scrape
    )


def get_tracer(name: str = __name__):
    return trace.get_tracer(name)
