# src/adapters/inbound/pubsub_subscriber.py
import asyncio
import logging
from src.core.ports.agent_service import AgentService
from src.core.ports.message_publisher import MessagePublisher
from src.core.models.message import AgentRequest

logger = logging.getLogger(__name__)

class BackgroundSubscriber:
    def __init__(self, agent_service: AgentService, publisher: MessagePublisher):
        self.agent_service = agent_service
        self.publisher = publisher
        self._running = False
        self._task = None

    async def process_raw_message(self, message_id: str, text: str):
        logger.info(f"Processing message {message_id}")
        request = AgentRequest(message_id=message_id, text=text)
        response = await self.agent_service.process_message(request)
        await self.publisher.publish(response)

    async def start(self):
        self._running = True
        self._task = asyncio.create_task(self._mock_listen_loop())
        logger.info("Background subscriber started")

    async def stop(self):
        self._running = False
        if self._task:
            self._task.cancel()
            try:
                await self._task
            except asyncio.CancelledError:
                pass
        logger.info("Background subscriber stopped")

    async def _mock_listen_loop(self):
        # Scaffolded for Pub/Sub. In a real app, this would use google.cloud.pubsub_v1
        # async pull methods or wrap the synchronous subscriber client in a thread pool
        while self._running:
            await asyncio.sleep(1)
