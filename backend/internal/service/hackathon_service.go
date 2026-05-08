package service

import "github.com/moficodes/hackathon-judge/backend/internal/domain"

type HackathonService interface {
	ListHackathons() ([]domain.Hackathon, error)
	ListProjectsByHackathon(id string) ([]domain.Project, error)
}

type hackathonService struct {
	repo        domain.HackathonRepository
	projectRepo domain.ProjectRepository
}

func NewHackathonService(repo domain.HackathonRepository, projectRepo domain.ProjectRepository) HackathonService {
	return &hackathonService{repo: repo, projectRepo: projectRepo}
}

func (s *hackathonService) ListHackathons() ([]domain.Hackathon, error) {
	return s.repo.GetAll()
}

func (s *hackathonService) ListProjectsByHackathon(id string) ([]domain.Project, error) {
	return s.projectRepo.GetByHackathonID(id)
}
