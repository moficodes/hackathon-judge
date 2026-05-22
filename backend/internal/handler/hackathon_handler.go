package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/moficodes/hackathon-judge/backend/internal/service"
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
		api.GET("/hackathons/:id", h.GetHackathon)
		api.GET("/hackathons/:id/projects", h.GetProjects)
		api.GET("/projects/:id", h.GetProject)
		api.GET("/projects/:id/evaluations", h.GetEvaluations)
		api.POST("/projects/:id/judge", h.TriggerJudgingAgent)
		api.POST("/projects/:id/judge/bq", h.TriggerJudgingBQ)
	}
}

func (h *HackathonHandler) GetHackathons(c *gin.Context) {
	res, err := h.svc.ListHackathons()
	if err != nil {
		log.Printf("[ERROR] GetHackathons failed: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *HackathonHandler) GetHackathon(c *gin.Context) {
	id := c.Param("id")
	res, err := h.svc.GetHackathon(id)
	if err != nil {
		log.Printf("[ERROR] GetHackathon failed for %s: %v\n", id, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Hackathon not found"})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *HackathonHandler) GetProjects(c *gin.Context) {
	id := c.Param("id")
	res, err := h.svc.ListProjectsByHackathon(id)
	if err != nil {
		log.Printf("[ERROR] GetProjects failed for hackathon %s: %v\n", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *HackathonHandler) GetProject(c *gin.Context) {
	id := c.Param("id")
	res, err := h.svc.GetProject(id)
	if err != nil {
		log.Printf("[ERROR] GetProject failed for %s: %v\n", id, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *HackathonHandler) GetEvaluations(c *gin.Context) {
	id := c.Param("id")
	res, err := h.svc.ListEvaluationsByProject(id)
	if err != nil {
		// Use a simple string check to determine if the project wasn't found
		if err.Error() == "failed to get project: project not found" {
			log.Printf("[WARNING] GetEvaluations project not found: %s\n", id)
			c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
			return
		}
		log.Printf("[ERROR] GetEvaluations failed for project %s: %v\n", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *HackathonHandler) TriggerJudgingAgent(c *gin.Context) {
	id := c.Param("id")
	taskID, err := h.svc.TriggerJudgingAgent(id)
	if err != nil {
		log.Printf("[ERROR] TriggerJudgingAgent failed for project %s: %v\n", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{
		"message": "Agent judging task created",
		"task_id": taskID,
	})
}

func (h *HackathonHandler) TriggerJudgingBQ(c *gin.Context) {
	id := c.Param("id")
	taskID, err := h.svc.TriggerJudging(id, true)
	if err != nil {
		log.Printf("[ERROR] TriggerJudgingBQ failed for project %s: %v\n", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{
		"message": "BQ AI judging task created",
		"task_id": taskID,
	})
}
