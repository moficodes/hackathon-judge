package repository

import (
	"testing"
	"time"
)

func TestMapBQHackathon(t *testing.T) {
	repo := &BigQueryRepo{} // We only need the method, not a real client

	date := time.Now()
	bqh := bqHackathon{
		ID:          "h1",
		Title:       "Hackathon 1",
		Date:        date,
		Description: "Desc",
		Goal:        "Goal",
		Status:      "ACTIVE",
		Criteria: []bqCriterion{
			{Name: "C1", Description: "D1", Weight: 0.5, MaxScore: 10},
		},
		BonusCriteria: []bqCriterion{},
	}

	h, err := repo.mapBQHackathon(bqh)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if h.ID != "h1" || h.Title != "Hackathon 1" || h.Description != "Desc" || h.Goal != "Goal" || h.Status != "ACTIVE" || !h.Date.Equal(date) {
		t.Errorf("basic fields mapped incorrectly. got: %+v", h)
	}

	if len(h.Criteria) != 1 || h.Criteria[0].Name != "C1" || h.Criteria[0].Weight != 0.5 || h.Criteria[0].MaxScore != 10 {
		t.Errorf("criteria mapped incorrectly. got: %+v", h.Criteria)
	}
}

func TestMapBQProject(t *testing.T) {
	repo := &BigQueryRepo{}

	// Note: bigquery.NullString is handled internally in mapBQProject, we need to pass a mocked version if needed or just use the struct
	bqp := bqProject{
		ID:          "p1",
		Name:        "Project 1",
		Title:       "Title 1",
		URL:         "url",
		GitHubURL:   "gh",
		TeamName:    "Team",
		HackathonID: "h1",
		Score:       95.5,
	}

	p := repo.mapBQProject(bqp)

	if p.ID != "p1" || p.Name != "Project 1" || p.Title != "Title 1" || p.URL != "url" || p.GitHubURL != "gh" || p.TeamName != "Team" || p.HackathonID != "h1" || p.Score != 95.5 {
		t.Errorf("fields mapped incorrectly. got: %+v", p)
	}
}
