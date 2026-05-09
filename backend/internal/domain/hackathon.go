package domain

import "time"

type Criterion struct {
	Name        string  `json:"name"`
	Weight      float64 `json:"weight"`
	Description string  `json:"description"`
}

type Hackathon struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Title       string      `json:"title"`
	Date        time.Time   `json:"date"`
	Description string      `json:"description"`
	Goal        string      `json:"goal"`
	Status      string      `json:"status"`
	Criteria    []Criterion `json:"criteria"`
}

type HackathonRepository interface {
	GetAll() ([]Hackathon, error)
	GetByID(id string) (Hackathon, error)
}
