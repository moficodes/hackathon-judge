package main

import (
	"log"
	"github.com/gin-gonic/gin"
	"github.com/moficodes/hackathon-judge/backend/internal/handler"
	"github.com/moficodes/hackathon-judge/backend/internal/repository"
	"github.com/moficodes/hackathon-judge/backend/internal/service"
)

func main() {
	r := gin.Default()

	repo := repository.NewMemoryRepo()
	svc := service.NewHackathonService(repo, repo)
	h := handler.NewHackathonHandler(svc)

	h.RegisterRoutes(r)

	log.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
