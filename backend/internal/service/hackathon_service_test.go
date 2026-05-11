package service_test

import (
	"testing"

	"github.com/moficodes/hackathon-judge/backend/internal/domain"
	"github.com/moficodes/hackathon-judge/backend/internal/repository"
	"github.com/moficodes/hackathon-judge/backend/internal/service"
	"github.com/stretchr/testify/assert"
)

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
