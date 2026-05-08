package service

import "github.com/moficodes/hackathon-judge/backend/internal/domain"

type HackathonService struct {
	repo        domain.HackathonRepository
	projectRepo domain.ProjectRepository
}

func NewHackathonService(repo domain.HackathonRepository, projectRepo domain.ProjectRepository) *HackathonService {
	return &HackathonService{repo: repo, projectRepo: projectRepo}
}

func (s *HackathonService) ListHackathons() ([]domain.Hackathon, error) {
	return s.repo.GetAll()
}

func (s *HackathonService) ListProjectsByHackathon(id string) ([]domain.Project, error) {
	return s.projectRepo.GetByHackathonID(id)
}
