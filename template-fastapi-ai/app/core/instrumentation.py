from opentelemetry import trace
from opentelemetry.exporter.otlp.proto.grpc.trace_exporter import OTLPSpanExporter
from opentelemetry.instrumentation.fastapi import FastAPIInstrumentor
from opentelemetry.sdk.resources import SERVICE_NAME, Resource
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.trace.export import BatchSpanProcessor


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
    # @see: https://opentelemetry-python-contrib.readthedocs.io/en/latest/instrumentation/fastapi/fastapi.html
    FastAPIInstrumentor.instrument_app(app, exclude_spans=["receive", "send"])


def get_tracer(name: str = __name__):
    return trace.get_tracer(name)
