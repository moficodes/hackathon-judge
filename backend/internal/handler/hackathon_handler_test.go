package handler

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/moficodes/hackathon-judge/backend/internal/domain"
	"github.com/moficodes/hackathon-judge/backend/internal/repository"
	"github.com/moficodes/hackathon-judge/backend/internal/service"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetHackathons(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	repo := repository.NewMemoryRepo()
	svc := service.NewHackathonService(repo, repo, repo)
	h := NewHackathonHandler(svc)
	h.RegisterRoutes(r)

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/api/hackathons", nil)
	assert.NoError(t, err)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var res []domain.Hackathon
	err = json.Unmarshal(w.Body.Bytes(), &res)
	assert.NoError(t, err)
	assert.Greater(t, len(res), 0)
}

func TestGetProjects(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	repo := repository.NewMemoryRepo()
	svc := service.NewHackathonService(repo, repo, repo)
	h := NewHackathonHandler(svc)
	h.RegisterRoutes(r)

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/api/hackathons/1/projects", nil)
	assert.NoError(t, err)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var res []domain.Project
	err = json.Unmarshal(w.Body.Bytes(), &res)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(res))
	assert.Equal(t, "p1", res[0].ID)
}

func TestGetEvaluations(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	repo := repository.NewMemoryRepo()
	svc := service.NewHackathonService(repo, repo, repo)
	h := NewHackathonHandler(svc)
	h.RegisterRoutes(r)

	// Pre-seed an evaluation via the service
	err := svc.AddEvaluation(domain.Evaluation{
		ID: "e1", ProjectID: "p1", JudgeID: "j1", TotalScore: 9.5,
	})
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/projects/p1/evaluations", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var res []domain.Evaluation
	err = json.Unmarshal(w.Body.Bytes(), &res)
	assert.NoError(t, err)
	assert.Greater(t, len(res), 0)
	assert.Equal(t, "e1", res[0].ID)
}
