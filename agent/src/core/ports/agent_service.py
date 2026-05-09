# src/core/ports/agent_service.py
from abc import ABC, abstractmethod
from src.core.models.message import AgentRequest, AgentResponse

class AgentService(ABC):
    @abstractmethod
    async def process_message(self, request: AgentRequest) -> AgentResponse:
        pass
