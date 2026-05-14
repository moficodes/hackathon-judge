# src/main.py
import logging

# Configure logging at the very beginning to ensure all logs are captured correctly
logging.basicConfig(
    level=logging.INFO,
    format="%(levelname)s:%(name)s:%(message)s",
    force=True
)

from fastapi import FastAPI
import uvicorn
import os
from dotenv import load_dotenv
from contextlib import asynccontextmanager

from src.adapters.outbound.sandbox_direct import SandboxDirectAdapter
from src.adapters.outbound.pubsub_publisher import PubSubPublisherAdapter, MockPubSubPublisherAdapter
from src.adapters.inbound.pubsub_subscriber import BackgroundSubscriber

# Load environment variables from .env file
load_dotenv()

import sys

def get_required_env(var_name: str) -> str:
    value = os.getenv(var_name)
    if not value:
        print(f"Error: {var_name} environment variable is required", file=sys.stderr)
        sys.exit(1)
    return value

# Configuration from Environment Variables
PROJECT_ID = get_required_env("GOOGLE_CLOUD_PROJECT")
TASKS_SUBSCRIPTION = get_required_env("TASKS_SUBSCRIPTION")
RESULTS_TOPIC = get_required_env("RESULTS_TOPIC")
USE_MOCK = os.getenv("USE_MOCK_PUBSUB", "false").lower() == "true"

# Dependency Injection setup
agent_service = SandboxDirectAdapter()

if USE_MOCK:
    publisher = MockPubSubPublisherAdapter()
    subscriber = BackgroundSubscriber(
        agent_service=agent_service,
        publisher=publisher
    )
else:
    publisher = PubSubPublisherAdapter(project_id=PROJECT_ID, topic_id=RESULTS_TOPIC)
    subscriber = BackgroundSubscriber(
        agent_service=agent_service, 
        publisher=publisher,
        project_id=PROJECT_ID,
        subscription_id=TASKS_SUBSCRIPTION
    )

@asynccontextmanager
async def lifespan(app: FastAPI):
    # Startup
    print("Starting up Hackathon Judge Agent...")
    await subscriber.start()
    print("Hackathon Judge Agent started.")
    yield
    # Shutdown
    print("Shutting down Hackathon Judge Agent...")
    await subscriber.stop()
    print("Hackathon Judge Agent stopped.")

app = FastAPI(title="Hackathon Judge Agent", lifespan=lifespan)

@app.get("/")
async def root():
    return {"message": "Welcome to the Hackathon Judge Agent"}

@app.get("/health")
async def health():
    return {"status": "healthy"}

if __name__ == "__main__":
    uvicorn.run("src.main:app", host="0.0.0.0", port=8000, reload=True)
