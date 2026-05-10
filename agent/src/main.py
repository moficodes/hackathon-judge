# src/main.py
from fastapi import FastAPI
import uvicorn
import os
from dotenv import load_dotenv
from contextlib import asynccontextmanager
import logging

from src.adapters.outbound.adk_agent import ADKAgentAdapter
from src.adapters.outbound.pubsub_publisher import PubSubPublisherAdapter, MockPubSubPublisherAdapter
from src.adapters.inbound.pubsub_subscriber import BackgroundSubscriber

# Load environment variables from .env file
load_dotenv()

logging.basicConfig(level=logging.INFO)

# Configuration from Environment Variables
PROJECT_ID = os.getenv("GOOGLE_CLOUD_PROJECT", "mofilabs")
TASKS_SUBSCRIPTION = os.getenv("TASKS_SUBSCRIPTION", "agent-judging-tasks-sub")
RESULTS_TOPIC = os.getenv("RESULTS_TOPIC", "judging-results")
USE_MOCK = os.getenv("USE_MOCK_PUBSUB", "false").lower() == "true"

# Dependency Injection setup
agent_service = ADKAgentAdapter()

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
    await subscriber.start()
    yield
    # Shutdown
    await subscriber.stop()

app = FastAPI(title="Hackathon Judge Agent", lifespan=lifespan)

@app.get("/")
async def root():
    return {"message": "Welcome to the Hackathon Judge Agent"}

@app.get("/health")
async def health():
    return {"status": "healthy"}

if __name__ == "__main__":
    uvicorn.run("src.main:app", host="0.0.0.0", port=8000, reload=True)
