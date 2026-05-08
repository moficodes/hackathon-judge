package domain

import "time"

type CriteriaScore struct {
	Name   string  `json:"name"`
	Score  float64 `json:"score"`
	Weight float64 `json:"weight"`
}

type Evaluation struct {
	ID         string          `json:"id"`
	ProjectID  string          `json:"project_id"`
	JudgeID    string          `json:"judge_id"`
	Criteria   []CriteriaScore `json:"criteria"`
	TotalScore float64         `json:"total_score"`
	Comment    string          `json:"comment"`
	CreatedAt  time.Time       `json:"created_at"`
}

type EvaluationRepository interface {
	Save(eval Evaluation) error
	GetByProjectID(projectID string) ([]Evaluation, error)
}
