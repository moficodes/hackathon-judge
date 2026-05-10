package pubsub

import (
	"context"
	"encoding/json"
	"log"

	"cloud.google.com/go/pubsub"
	"github.com/moficodes/hackathon-judge/backend/internal/domain"
)

type GoogleResultSubscriber struct {
	client *pubsub.Client
	sub    *pubsub.Subscription
}

func NewGoogleResultSubscriber(projectID, subID string) (*GoogleResultSubscriber, error) {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}

	sub := client.Subscription(subID)
	return &GoogleResultSubscriber{
		client: client,
		sub:    sub,
	}, nil
}

// Start listens for messages and passes them to the provided handler. It blocks until ctx is canceled or an error occurs.
func (s *GoogleResultSubscriber) Start(ctx context.Context, handle func(domain.JudgingResult) error) error {
	log.Printf("Listening for messages on subscription: %s", s.sub.String())
	return s.sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		var result domain.JudgingResult
		if err := json.Unmarshal(msg.Data, &result); err != nil {
			log.Printf("Failed to unmarshal judging result: %v", err)
			msg.Nack()
			return
		}

		if err := handle(result); err != nil {
			log.Printf("Failed to handle judging result: %v", err)
			msg.Nack()
			return
		}

		msg.Ack()
	})
}
