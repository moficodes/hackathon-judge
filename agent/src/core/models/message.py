# src/core/models/message.py
from pydantic import BaseModel, Field
from typing import List, Optional

class ScoringCriteria(BaseModel):
    name: str
    weight: float
    max_score: float = 10.0

class AgentRequest(BaseModel):
    """The incoming task from the backend (Judging Task)"""
    task_id: str
    project_name: str
    github_url: str
    submission_text: str
    judging_rubric: str
    scoring_criteria: List[ScoringCriteria]

class CategoryScore(BaseModel):
    name: str
    score: float
    reasoning: str
    max_score: float = 10.0

class AgentResponse(BaseModel):
    """The result sent back to the backend (Judging Result)"""
    task_id: str
    status: str = Field(description="'success' or 'error'")
    error_message: Optional[str] = None
    scores: List[CategoryScore] = []
    total_score: float = 0.0
    overall_comments: str = ""
    confidence_score: float = 0.0
