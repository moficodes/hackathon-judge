# Hackathon Judge Agent Context

This `GEMINI.md` file contains the instructional context and architectural guidance for the `agent` project. It serves as the primary source of truth for the codebase's structure, conventions, and operational procedures.

## Project Overview

The Hackathon Judge Agent is a production-ready Python application designed to:
1. Listen to a Google Cloud Pub/Sub subscription in the background.
2. Asynchronously process incoming messages using a Google ADK (Agent Development Kit) agent.
3. Publish the agent's responses back to a Pub/Sub topic.

The application leverages **FastAPI** to manage the lifecycle of the background tasks (using `lifespan` events) and to expose basic HTTP endpoints (like `/health` and `/`).

### Key Technologies
- **Python:** `3.12+`
- **Package Manager:** `uv`
- **Web Framework:** `fastapi`, `uvicorn`
- **Messaging:** `google-cloud-pubsub`
- **Agent SDK:** `google-adk`
- **Testing:** `pytest`, `pytest-asyncio`, `httpx`

## Architecture: Hexagonal (Ports and Adapters)

The project strictly follows a **Hexagonal Architecture** to decouple the core business logic from external delivery mechanisms and infrastructure dependencies.

### Directory Structure

```text
agent/
├── src/
│   ├── core/                  # Domain/Core logic (Independent of frameworks)
│   │   ├── ports/             # Interfaces (Abstract Base Classes)
│   │   │   ├── message_publisher.py
│   │   │   └── agent_service.py
│   │   └── models/            # Core data models (Pydantic)
│   │       └── message.py
│   ├── adapters/              # External integrations (Depends on core ports)
│   │   ├── inbound/           # Driving adapters (Triggers the application)
│   │   │   └── pubsub_subscriber.py
│   │   └── outbound/          # Driven adapters (Infrastructure/Services)
│   │       ├── pubsub_publisher.py
│   │       └── adk_agent.py
│   └── main.py                # Composition root & FastAPI Application
├── tests/
│   ├── api/                   # Integration tests for FastAPI endpoints
│   ├── adapters/              # Tests for inbound/outbound adapters
│   └── core/                  # Tests for domain models and business logic
```

### Key Components

- **Core Layer:** Defines the `AgentRequest` and `AgentResponse` models using Pydantic. It also defines the contracts (ports) that adapters must implement: `AgentService` and `MessagePublisher`.
- **Inbound Adapters:** The `BackgroundSubscriber` listens for Pub/Sub messages, translates them into `AgentRequest` models, and invokes the `AgentService`.
- **Outbound Adapters:**
  - `ADKAgentAdapter`: Implements the `AgentService` port and interacts with the `google-adk` SDK.
  - `MockPubSubPublisherAdapter` (and eventually the real publisher): Implements the `MessagePublisher` port to send responses.
- **Composition Root:** `src/main.py` is responsible for initializing the concrete adapters, injecting them into the subscriber, and attaching the subscriber to the FastAPI application's `lifespan`.

## Development Conventions

1. **Dependency Injection:** Dependencies should be passed into classes rather than instantiated internally. This is crucial for maintaining the Hexagonal Architecture and allowing easy substitution of mock adapters during testing.
2. **Asynchronous I/O:** The application is highly asynchronous. Use `async`/`await` for all I/O bound operations (Pub/Sub calls, ADK agent invocations, FastAPI endpoints) to prevent blocking the event loop.
3. **Type Hinting:** Comprehensive type hinting is required for all function signatures and class properties.
4. **Testing (TDD):** The project follows a Test-Driven Development approach.
   - Tests are located in the `tests/` directory, mirroring the structure of `src/`.
   - Use `pytest` for all testing.
   - Use `pytest.mark.asyncio` for testing async functions.
   - Run tests before committing any changes.

## Building and Running

### Setup
Ensure `uv` is installed, then sync the environment:
```bash
uv sync
```

### Running the Application
Run the FastAPI application (which automatically starts the background subscriber):
```bash
# Standard run
uv run python src/main.py

# Live-reloading for development
uv run uvicorn src.main:app --reload
```

### Running Tests
Execute the full test suite:
```bash
uv run pytest
```
*(Note: `pyproject.toml` is configured with `pythonpath = ["."]` so imports work seamlessly).*
