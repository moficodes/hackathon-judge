package repository

import (
	"errors"
	"sync"

	"github.com/moficodes/hackathon-judge/backend/internal/domain"
)

type memoryRepo struct {
	mu          sync.RWMutex
	hackathons  []domain.Hackathon
	projects    []domain.Project
	evaluations []domain.Evaluation
}

func NewMemoryRepo() *memoryRepo {
	return &memoryRepo{
		hackathons: []domain.Hackathon{
			{ID: "1", Name: "Hack1", Title: "Summer Hack", Status: "Active"},
		},
		projects: []domain.Project{
			{ID: "p1", Name: "Proj1", HackathonID: "1", Score: 0},
		},
		evaluations: []domain.Evaluation{},
	}
}

func (r *memoryRepo) GetAll() ([]domain.Hackathon, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return append([]domain.Hackathon(nil), r.hackathons...), nil
}

func (r *memoryRepo) GetByID(id string) (domain.Hackathon, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, h := range r.hackathons {
		if h.ID == id {
			return h, nil
		}
	}
	return domain.Hackathon{}, errors.New("hackathon not found")
}

func (r *memoryRepo) GetByHackathonID(hackathonID string) ([]domain.Project, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []domain.Project
	for _, p := range r.projects {
		if p.HackathonID == hackathonID {
			result = append(result, p)
		}
	}
	return result, nil
}

func (r *memoryRepo) Save(eval domain.Evaluation) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.evaluations = append(r.evaluations, eval)
	return nil
}

func (r *memoryRepo) GetByProjectID(projectID string) ([]domain.Evaluation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []domain.Evaluation
	for _, e := range r.evaluations {
		if e.ProjectID == projectID {
			result = append(result, e)
		}
	}
	return result, nil
}

func (r *memoryRepo) UpdateScore(projectID string, score float64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, p := range r.projects {
		if p.ID == projectID {
			r.projects[i].Score = score
			return nil
		}
	}
	return errors.New("project not found")
}
