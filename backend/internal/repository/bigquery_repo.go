package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/moficodes/hackathon-judge/backend/internal/domain"
	"google.golang.org/api/iterator"
)

type BigQueryRepo struct {
	client    *bigquery.Client
	projectID string
}

func NewBigQueryRepo(client *bigquery.Client, projectID string) *BigQueryRepo {
	return &BigQueryRepo{
		client:    client,
		projectID: projectID,
	}
}

type bqHackathon struct {
	ID            string    `bigquery:"id"`
	Title         string    `bigquery:"title"`
	Date          time.Time `bigquery:"date"`
	Description   string    `bigquery:"description"`
	Goal          string    `bigquery:"goal"`
	Status        string    `bigquery:"status"`
	Criteria      string    `bigquery:"criteria"`
	BonusCriteria string    `bigquery:"bonus_criteria"`
}

func (r *BigQueryRepo) mapBQHackathon(bqH bqHackathon) (domain.Hackathon, error) {
	h := domain.Hackathon{
		ID:          bqH.ID,
		Title:       bqH.Title,
		Date:        bqH.Date,
		Description: bqH.Description,
		Goal:        bqH.Goal,
		Status:      bqH.Status,
	}

	if bqH.Criteria != "" {
		if err := json.Unmarshal([]byte(bqH.Criteria), &h.Criteria); err != nil {
			return h, fmt.Errorf("failed to unmarshal criteria: %w", err)
		}
	}
	if bqH.BonusCriteria != "" {
		if err := json.Unmarshal([]byte(bqH.BonusCriteria), &h.BonusCriteria); err != nil {
			return h, fmt.Errorf("failed to unmarshal bonus criteria: %w", err)
		}
	}
	return h, nil
}

func (r *BigQueryRepo) GetAll() ([]domain.Hackathon, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	query := r.client.Query(fmt.Sprintf("SELECT * FROM `%s.hackathons.hackathons`", r.projectID))
	it, err := query.Read(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to read hackathons: %w", err)
	}

	var hackathons []domain.Hackathon
	for {
		var row bqHackathon
		err := it.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to iterate hackathons: %w", err)
		}
		h, err := r.mapBQHackathon(row)
		if err != nil {
			return nil, fmt.Errorf("failed to map hackathon row: %w", err)
		}
		hackathons = append(hackathons, h)
	}
	return hackathons, nil
}

func (r *BigQueryRepo) GetByID(id string) (domain.Hackathon, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	query := r.client.Query(fmt.Sprintf("SELECT * FROM `%s.hackathons.hackathons` WHERE id = @id LIMIT 1", r.projectID))
	query.Parameters = []bigquery.QueryParameter{
		{Name: "id", Value: id},
	}
	it, err := query.Read(ctx)
	if err != nil {
		return domain.Hackathon{}, fmt.Errorf("failed to read hackathon: %w", err)
	}

	var row bqHackathon
	err = it.Next(&row)
	if err == iterator.Done {
		return domain.Hackathon{}, fmt.Errorf("hackathon not found")
	}
	if err != nil {
		return domain.Hackathon{}, fmt.Errorf("failed to iterate hackathon: %w", err)
	}

	return r.mapBQHackathon(row)
}
