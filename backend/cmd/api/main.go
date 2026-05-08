package main

import (
	"log"
	"github.com/gin-gonic/gin"
	"github.com/moficodes/hackathon-judge/backend/internal/handler"
	"github.com/moficodes/hackathon-judge/backend/internal/repository"
	"github.com/moficodes/hackathon-judge/backend/internal/service"
	"github.com/moficodes/hackathon-judge/backend/pkg/logger"
)

const defaultPort = ":8080"

func main() {
	logger.Init()
	r := gin.Default()

	repo := repository.NewMemoryRepo()
	svc := service.NewHackathonService(repo, repo)
	h := handler.NewHackathonHandler(svc)

	h.RegisterRoutes(r)

	log.Printf("Server starting on %s\n", defaultPort)
	if err := r.Run(defaultPort); err != nil {
		log.Fatal(err)
	}
}
