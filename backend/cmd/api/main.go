package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/moficodes/hackathon-judge/backend/internal/domain"
	"github.com/moficodes/hackathon-judge/backend/internal/handler"
	"github.com/moficodes/hackathon-judge/backend/internal/infrastructure/pubsub"
	"github.com/moficodes/hackathon-judge/backend/internal/repository"
	"github.com/moficodes/hackathon-judge/backend/internal/service"
	"github.com/moficodes/hackathon-judge/backend/pkg/logger"
)

const defaultPort = ":8080"

func main() {
	logger.Init()
	
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	r := gin.Default()

	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		projectID = "mofilabs"
	}
	topicID := os.Getenv("TASKS_TOPIC")
	if topicID == "" {
		topicID = "judging-tasks"
	}
	subID := os.Getenv("RESULTS_SUB")
	if subID == "" {
		subID = "backend-judging-results-sub"
	}

	publisher, err := pubsub.NewGoogleTaskPublisher(projectID, topicID)
	if err != nil {
		log.Printf("Warning: failed to initialize PubSub publisher: %v. Using mock.", err)
		publisher = nil
	}

	subscriber, err := pubsub.NewGoogleResultSubscriber(projectID, subID)
	if err != nil {
		log.Printf("Warning: failed to initialize PubSub subscriber: %v. Result listening disabled.", err)
	} else {
		go func() {
			log.Println("Starting background result subscriber...")
			err := subscriber.Start(context.Background(), func(res domain.JudgingResult) error {
				log.Printf("--- RECEIVED JUDGING RESULT ---")
				log.Printf("Task ID: %s", res.TaskID)
				log.Printf("Status: %s", res.Status)
				if res.Status == "error" && res.ErrorMessage != nil {
					log.Printf("Error: %s", *res.ErrorMessage)
				}
				log.Printf("Total Score: %.2f", res.TotalScore)
				log.Printf("Comments: %s", res.OverallComments)
				log.Printf("-------------------------------")
				return nil
			})
			if err != nil {
				log.Printf("Subscriber stopped with error: %v", err)
			}
		}()
	}

	repo := repository.NewMemoryRepo()
	svc := service.NewHackathonService(repo, repo, repo, publisher)
	h := handler.NewHackathonHandler(svc)

	h.RegisterRoutes(r)

	log.Printf("Server starting on %s\n", defaultPort)
	if err := r.Run(defaultPort); err != nil {
		log.Fatal(err)
	}
}
