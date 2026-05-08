package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/gin-gonic/gin"
	"github.com/moficodes/hackathon-judge/backend/internal/domain"
	"github.com/moficodes/hackathon-judge/backend/internal/repository"
	"github.com/moficodes/hackathon-judge/backend/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestGetHackathons(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	repo := repository.NewMemoryRepo()
	svc := service.NewHackathonService(repo, repo)
	h := NewHackathonHandler(svc)
	h.RegisterRoutes(r)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/hackathons", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var res []domain.Hackathon
	json.Unmarshal(w.Body.Bytes(), &res)
	assert.Greater(t, len(res), 0)
}
