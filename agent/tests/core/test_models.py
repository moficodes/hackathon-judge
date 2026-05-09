# tests/core/test_models.py
from src.core.models.message import AgentRequest, AgentResponse

def test_agent_request_model():
    req = AgentRequest(message_id="123", text="Hello")
    assert req.message_id == "123"
    assert req.text == "Hello"

def test_agent_response_model():
    res = AgentResponse(original_message_id="123", response_text="Hi there")
    assert res.original_message_id == "123"
    assert res.response_text == "Hi there"
