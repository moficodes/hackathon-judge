# src/adapters/outbound/pubsub_publisher.py
from src.core.ports.message_publisher import MessagePublisher
from src.core.models.message import AgentResponse
import logging

logger = logging.getLogger(__name__)

class MockPubSubPublisherAdapter(MessagePublisher):
    def __init__(self):
        self.published_messages = []

    async def publish(self, response: AgentResponse) -> None:
        logger.info(f"Mock publishing: {response.response_text}")
        self.published_messages.append(response)

# In a real implementation, you would have:
# class PubSubPublisherAdapter(MessagePublisher):
#    def __init__(self, project_id: str, topic_id: str): ...
#    async def publish(self, response: AgentResponse) -> None: ...
