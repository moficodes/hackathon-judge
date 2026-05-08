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
		ID:         "e1",
		ProjectID:  "p1",
		TotalScore: 80.0,
	}

	err := svc.AddEvaluation(eval1)
	assert.NoError(t, err)

	projects, err := svc.ListProjectsByHackathon("1")
	assert.NoError(t, err)
	assert.Len(t, projects, 1)
	assert.Equal(t, 80.0, projects[0].Score)

	eval2 := domain.Evaluation{
		ID:         "e2",
		ProjectID:  "p1",
		TotalScore: 90.0,
	}

	err = svc.AddEvaluation(eval2)
	assert.NoError(t, err)

	projects, err = svc.ListProjectsByHackathon("1")
	assert.NoError(t, err)
	assert.Len(t, projects, 1)
	assert.Equal(t, 85.0, projects[0].Score)
}
