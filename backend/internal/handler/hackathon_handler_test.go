package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/moficodes/hackathon-judge/backend/internal/domain"
	"github.com/moficodes/hackathon-judge/backend/internal/handler"
	"github.com/moficodes/hackathon-judge/backend/internal/repository"
	"github.com/moficodes/hackathon-judge/backend/internal/service"
	"github.com/stretchr/testify/assert"
)

func setupRouter() (*gin.Engine, service.HackathonService) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	repo := repository.NewMemoryRepo()
	svc := service.NewHackathonService(repo, repo, repo, nil)
	h := handler.NewHackathonHandler(svc)
	h.RegisterRoutes(r)
	return r, svc
}

func TestListHackathons(t *testing.T) {
	r, _ := setupRouter()

	req, _ := http.NewRequest(http.MethodGet, "/api/hackathons", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var res []domain.Hackathon
	err := json.Unmarshal(w.Body.Bytes(), &res)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(res), 1)
	assert.Equal(t, "Summer Hack", res[0].Title)
}

func TestListProjectsByHackathon(t *testing.T) {
	r, _ := setupRouter()

	req, _ := http.NewRequest(http.MethodGet, "/api/hackathons/1/projects", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var res []domain.Project
	err := json.Unmarshal(w.Body.Bytes(), &res)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(res), 1)
	assert.Equal(t, "p1", res[0].ID)
}

func TestListEvaluationsByProject(t *testing.T) {
	r, svc := setupRouter()

	// Pre-seed an evaluation via the service
	err := svc.AddEvaluation(domain.Evaluation{
		ID: "e1", ProjectID: "p1", JudgeID: "j1",
		Criteria: []domain.CriteriaScore{
			{Name: "Technology", Score: 4}, // 8
			{Name: "Design", Score: 5},     // 5 -> 13
		},
	})
	assert.NoError(t, err)

	req, _ := http.NewRequest(http.MethodGet, "/api/projects/p1/evaluations", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var res []domain.Evaluation
	err = json.Unmarshal(w.Body.Bytes(), &res)
	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, "e1", res[0].ID)
	assert.Equal(t, 13.0, res[0].TotalScore)
}
