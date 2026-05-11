FROM ghcr.io/astral-sh/uv:python3.12-alpine AS builder
ENV UV_COMPILE_BYTECODE=1 UV_LINK_MODE=copy
WORKDIR /app
COPY pyproject.toml uv.lock ./
RUN uv sync --frozen --no-install-project --no-dev
COPY . .
RUN uv sync --frozen --no-dev

FROM python:3.12-alpine
WORKDIR /app
ENV PATH="/app/.venv/bin:$PATH"
COPY --from=builder /app /app
EXPOSE 8000
ENTRYPOINT ["python", "src/main.py"]
