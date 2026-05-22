package service_test

import (
	"testing"
	"time"

	"github.com/moficodes/hackathon-judge/backend/internal/domain"
	"github.com/moficodes/hackathon-judge/backend/internal/repository"
	"github.com/moficodes/hackathon-judge/backend/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestGetHackathon(t *testing.T) {
	repo := repository.NewMemoryRepo()
	svc := service.NewHackathonService(repo, repo, repo, nil)

	h, err := svc.GetHackathon("1")
	assert.NoError(t, err)
	assert.Equal(t, "1", h.ID)
	assert.Equal(t, "Summer Hack", h.Title)

	_, err = svc.GetHackathon("invalid")
	assert.Error(t, err)
}

func TestGetProject(t *testing.T) {
	repo := repository.NewMemoryRepo()
	svc := service.NewHackathonService(repo, repo, repo, nil)

	p, err := svc.GetProject("p1")
	assert.NoError(t, err)
	assert.Equal(t, "p1", p.ID)
	assert.Equal(t, "Proj1", p.Name)

	_, err = svc.GetProject("invalid")
	assert.Error(t, err)
}

func TestAddEvaluationUpdatesProjectScore(t *testing.T) {
	repo := repository.NewMemoryRepo()
	svc := service.NewHackathonService(repo, repo, repo, nil)

	// memoryRepo initializes a project with ID "p1" under hackathon "1"
	eval1 := domain.Evaluation{
		ID:        "e1",
		ProjectID: "p1",
		Criteria: []domain.CriteriaScore{
			{Name: "Innovation & Originality", Score: 4}, // Weight is 0.2 from mock data -> 0.8
			{Name: "Technical Execution", Score: 5},      // Weight is 0.25 from mock data -> 1.25
		},
	}

	err := svc.AddEvaluation(eval1)
	assert.NoError(t, err)

	projects, err := svc.ListProjectsByHackathon("1")
	assert.NoError(t, err)
	assert.Len(t, projects, 1)
	assert.Equal(t, 2.05, projects[0].Score) // (4*0.2) + (5*0.25) = 0.8 + 1.25 = 2.05

	eval2 := domain.Evaluation{
		ID:        "e2",
		ProjectID: "p1",
		Criteria: []domain.CriteriaScore{
			{Name: "Innovation & Originality", Score: 2}, // Weight 0.2 -> 0.4
			{Name: "Technical Execution", Score: 3},      // Weight 0.25 -> 0.75
			{Name: "Clean Code", Score: 2},               // Bonus: Weight 2 -> 4.0
		},
	}

	err = svc.AddEvaluation(eval2)
	assert.NoError(t, err)

	projects, err = svc.ListProjectsByHackathon("1")
	assert.NoError(t, err)
	assert.Len(t, projects, 1)
	assert.InDelta(t, 3.6, projects[0].Score, 0.0001) // (2.05 + 5.15) / 2 = 7.2 / 2 = 3.6
}

func TestTriggerJudgingDirectBQ(t *testing.T) {
	repo := repository.NewMemoryRepo()
	svc := service.NewHackathonService(repo, repo, repo, nil)

	taskID, err := svc.TriggerJudging("p1", true)
	assert.NoError(t, err)
	assert.NotEmpty(t, taskID)

	// Initially it should be RUNNING
	eval, err := repo.GetEvaluationByID(taskID)
	assert.NoError(t, err)
	assert.Equal(t, "RUNNING", eval.Status)

	// Wait a brief moment for the background goroutine to finish scoring
	time.Sleep(50 * time.Millisecond)

	// Now it should be SUCCESS
	eval, err = repo.GetEvaluationByID(taskID)
	assert.NoError(t, err)
	assert.Equal(t, "SUCCESS", eval.Status)
	assert.Equal(t, "BQ AI Function", eval.JudgeID)
	assert.NotEmpty(t, eval.Criteria)

	// Sum of weights:
	// Innovation & Originality: 4.5 * 0.2 = 0.9
	// Theme Alignment: 4.0 * 0.25 = 1.0
	// Technical Execution: 4.0 * 0.25 = 1.0
	// UX/UI: 5.0 * 0.2 = 1.0
	// Pitch: 4.5 * 0.1 = 0.45
	// Total score expected = 0.9 + 1.0 + 1.0 + 1.0 + 0.45 = 4.35
	assert.InDelta(t, 4.35, eval.TotalScore, 0.0001)
	assert.Equal(t, "Automated BigQuery AI evaluation completed successfully", eval.Comment)

	// Verify project score got updated
	p, err := repo.GetProjectByID("p1")
	assert.NoError(t, err)
	assert.InDelta(t, 4.35, p.Score, 0.0001)
}
