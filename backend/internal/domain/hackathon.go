package domain

import "time"

type Hackathon struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Title       string    `json:"title"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	Goal        string    `json:"goal"`
	Status      string    `json:"status"`
}

type HackathonRepository interface {
	GetAll() ([]Hackathon, error)
	GetByID(id string) (Hackathon, error)
}
