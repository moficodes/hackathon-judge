# src/main.py
from fastapi import FastAPI
import uvicorn
from contextlib import asynccontextmanager
import logging

from src.adapters.outbound.adk_agent import ADKAgentAdapter
from src.adapters.outbound.pubsub_publisher import MockPubSubPublisherAdapter
from src.adapters.inbound.pubsub_subscriber import BackgroundSubscriber

logging.basicConfig(level=logging.INFO)

# Dependency Injection setup
agent_service = ADKAgentAdapter()
publisher = MockPubSubPublisherAdapter()
subscriber = BackgroundSubscriber(agent_service=agent_service, publisher=publisher)

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
