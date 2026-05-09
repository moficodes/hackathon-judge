# src/core/ports/message_publisher.py
from abc import ABC, abstractmethod
from src.core.models.message import AgentResponse

class MessagePublisher(ABC):
    @abstractmethod
    async def publish(self, response: AgentResponse) -> None:
        pass
