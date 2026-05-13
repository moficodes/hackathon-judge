package domain

import "time"

// Evaluation Models (Existing)

type CriteriaScore struct {
	Name      string  `json:"name"`
	Score     float64 `json:"score"`
	Reasoning string  `json:"reasoning"`
	MaxScore  float64 `json:"max_score"`
	Weight    float64 `json:"weight"`
}

type Evaluation struct {
	ID         string          `json:"id"`
	ProjectID  string          `json:"project_id"`
	JudgeID    string          `json:"judge_id"`
	Status     string          `json:"status"` // "UNKNOWN", "SUCCESS", "RUNNING", "FAILED"
	Criteria   []CriteriaScore `json:"criteria"`
	TotalScore float64         `json:"total_score"`
	Comment    string          `json:"comment"`
	CreatedAt  time.Time       `json:"created_at"`
}

type EvaluationRepository interface {
	Save(eval Evaluation) error
	Update(eval Evaluation) error
	GetByProjectID(projectID string) ([]Evaluation, error)
	GetEvaluationByID(id string) (Evaluation, error)
}

type JudgingTask struct {
	TaskID          string      `json:"task_id"`
	ProjectName     string      `json:"project_name"`
	GithubURL       string      `json:"github_url"`
	SubmissionText  string      `json:"submission_text"`
	JudgingRubric   string      `json:"judging_rubric"`
	ScoringCriteria []Criterion `json:"scoring_criteria"`
}

type JudgingResult struct {
	TaskID          string          `json:"task_id"`
	Status          string          `json:"status"` // "success" or "error"
	ErrorMessage    *string         `json:"error_message,omitempty"`
	Scores          []CriteriaScore `json:"scores"`
	TotalScore      float64         `json:"total_score"`
	OverallComments string          `json:"overall_comments"`
	ConfidenceScore float64         `json:"confidence_score"`
}

type TaskPublisher interface {
	PublishTask(task JudgingTask) error
}

