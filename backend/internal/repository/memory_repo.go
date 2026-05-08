package repository

import (
	"errors"
	"github.com/moficodes/hackathon-judge/backend/internal/domain"
)

type memoryRepo struct {
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

func (r *memoryRepo) Save(eval domain.Evaluation) error {
	r.evaluations = append(r.evaluations, eval)

	// Recalculate average project score
	var total float64
	var count float64
	for _, e := range r.evaluations {
		if e.ProjectID == eval.ProjectID {
			total += e.TotalScore
			count++
		}
	}

	if count > 0 {
		avg := total / count
		for i, p := range r.projects {
			if p.ID == eval.ProjectID {
				r.projects[i].Score = avg
				break
			}
		}
	}

	return nil
}

func (r *memoryRepo) GetByProjectID(projectID string) ([]domain.Evaluation, error) {
	var result []domain.Evaluation
	for _, e := range r.evaluations {
		if e.ProjectID == projectID {
			result = append(result, e)
		}
	}
	return result, nil
}
