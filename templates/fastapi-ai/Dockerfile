FROM python:3.13-slim-bookworm

# The installer requires curl (and certificates) to download the release archive
RUN apt-get update && apt-get install -y --no-install-recommends \
    curl \
    ca-certificates \
    ffmpeg \
    libpq-dev \
    build-essential

# Download the latest installer
RUN curl -LsSf https://astral.sh/uv/install.sh | sh

WORKDIR /opt

# Copy necessary files
COPY pyproject.toml README.md ./

# Ensure the installed binary is on the `PATH`
ENV PATH="/root/.local/bin/:$PATH"

# Set default environment variables
ARG SERVER_ENV=production
ARG HOST=0.0.0.0
ARG PORT={{ port_number }}
ENV SERVER_ENV=${SERVER_ENV}
ENV DEBUG=false

ENV HOST=${HOST}
ENV PORT=${PORT}

# Copy the rest of the application
COPY . .

# Set up the virtual environment and install dependencies
RUN uv venv && uv sync

# Expose the port
EXPOSE ${PORT}

# Run the application using a shell command
CMD ["/bin/sh", "-c", "/root/.local/bin/uv run fastapi run /opt/app/main.py --port=${PORT}"]
