# src/core/models/message.py
from pydantic import BaseModel

class AgentRequest(BaseModel):
    message_id: str
    text: str

class AgentResponse(BaseModel):
    original_message_id: str
    response_text: str
