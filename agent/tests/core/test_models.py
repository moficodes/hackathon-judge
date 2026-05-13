# tests/core/test_models.py
from src.core.models.message import AgentRequest, AgentResponse, ScoringCriteria, CategoryScore

def test_agent_request_model():
    req = AgentRequest(
        task_id="tsk_123",
        project_name="Test Project",
        github_url="https://github.com/moficodes/test",
        submission_text="Here is my code",
        judging_rubric="Be nice",
        scoring_criteria=[ScoringCriteria(name="Innovation", description="How new is this?", weight=0.5)]
    )
    assert req.task_id == "tsk_123"
    assert req.scoring_criteria[0].max_score == 10.0
    assert req.scoring_criteria[0].description == "How new is this?"

def test_agent_response_model():
    res = AgentResponse(
        task_id="tsk_123",
        status="success",
        scores=[CategoryScore(name="Innovation", score=8, reasoning="Good")],
        total_score=4.0
    )
    assert res.task_id == "tsk_123"
    assert res.status == "success"
    assert len(res.scores) == 1
