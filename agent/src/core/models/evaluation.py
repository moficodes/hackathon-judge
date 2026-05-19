from pydantic import BaseModel, Field
from typing import List

class EvaluationScore(BaseModel):
    name: str = Field(description="The name of the scoring criteria category.")
    score: float = Field(description="The awarded score.")
    reasoning: str = Field(description="The reasoning behind this score.")

class EvaluationOutput(BaseModel):
    scores: List[EvaluationScore]
    total_score: float = Field(description="The total score summing up all category scores.")
    overall_comments: str = Field(description="Overall comments on the project.")
    confidence_score: float = Field(description="Score between 0.0 and 1.0 indicating confidence in the evaluation.")
