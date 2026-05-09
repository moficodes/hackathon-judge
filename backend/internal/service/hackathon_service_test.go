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
	svc := service.NewHackathonService(repo, repo, repo)

	// memoryRepo initializes a project with ID "p1" under hackathon "1"
	eval1 := domain.Evaluation{
		ID:        "e1",
		ProjectID: "p1",
		Criteria: []domain.CriteriaScore{
			{Name: "Technology", Score: 4}, // Weight is 2 from mock data -> 8
			{Name: "Design", Score: 5},     // Weight is 1 from mock data -> 5
		},
	}

	err := svc.AddEvaluation(eval1)
	assert.NoError(t, err)

	projects, err := svc.ListProjectsByHackathon("1")
	assert.NoError(t, err)
	assert.Len(t, projects, 1)
	assert.Equal(t, 13.0, projects[0].Score) // (4*2) + (5*1) = 13

	eval2 := domain.Evaluation{
		ID:        "e2",
		ProjectID: "p1",
		Criteria: []domain.CriteriaScore{
			{Name: "Technology", Score: 2}, // Weight 2 -> 4
			{Name: "Design", Score: 3},     // Weight 1 -> 3
		},
	}

	err = svc.AddEvaluation(eval2)
	assert.NoError(t, err)

	projects, err = svc.ListProjectsByHackathon("1")
	assert.NoError(t, err)
	assert.Len(t, projects, 1)
	assert.Equal(t, 10.0, projects[0].Score) // (13 + 7) / 2 = 10
}
