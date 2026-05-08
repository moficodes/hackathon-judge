package domain

import "time"

type Project struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Title       string    `json:"title"`
	URL         string    `json:"url"`
	GitHubURL   string    `json:"github_url"`
	TeamName    string    `json:"team_name"`
	Document    string    `json:"document"`
	Date        time.Time `json:"date"`
	HackathonID string    `json:"hackathon_id"`
	Score       float64   `json:"score"`
}

type ProjectRepository interface {
	GetByHackathonID(hackathonID string) ([]Project, error)
	UpdateScore(projectID string, score float64) error
}
