package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/moficodes/hackathon-judge/backend/internal/service"
	"net/http"
)

type HackathonHandler struct {
	svc service.HackathonService
}

func NewHackathonHandler(svc service.HackathonService) *HackathonHandler {
	return &HackathonHandler{svc: svc}
}

func (h *HackathonHandler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.GET("/hackathons", h.GetHackathons)
		api.GET("/hackathons/:id/projects", h.GetProjects)
	}
}

func (h *HackathonHandler) GetHackathons(c *gin.Context) {
	res, err := h.svc.ListHackathons()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *HackathonHandler) GetProjects(c *gin.Context) {
	id := c.Param("id")
	res, err := h.svc.ListProjectsByHackathon(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}
