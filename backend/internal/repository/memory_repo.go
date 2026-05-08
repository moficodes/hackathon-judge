package repository

import (
	"errors"
	"github.com/moficodes/hackathon-judge/backend/internal/domain"
)

type memoryRepo struct {
	hackathons []domain.Hackathon
	projects   []domain.Project
}

func NewMemoryRepo() *memoryRepo {
	return &memoryRepo{
		hackathons: []domain.Hackathon{
			{ID: "1", Name: "Hack1", Title: "Summer Hack", Status: "Active"},
		},
		projects: []domain.Project{
			{ID: "p1", Name: "Proj1", HackathonID: "1"},
		},
	}
}

func (r *memoryRepo) GetAll() ([]domain.Hackathon, error) {
	return r.hackathons, nil
}

func (r *memoryRepo) GetByID(id string) (domain.Hackathon, error) {
	for _, h := range r.hackathons {
		if h.ID == id {
			return h, nil
		}
	}
	return domain.Hackathon{}, errors.New("hackathon not found")
}

func (r *memoryRepo) GetByHackathonID(hackathonID string) ([]domain.Project, error) {
	var result []domain.Project
	for _, p := range r.projects {
		if p.HackathonID == hackathonID {
			result = append(result, p)
		}
	}
	return result, nil
}
