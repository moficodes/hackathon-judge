package service

import "github.com/moficodes/hackathon-judge/backend/internal/domain"

type HackathonService interface {
	ListHackathons() ([]domain.Hackathon, error)
	ListProjectsByHackathon(id string) ([]domain.Project, error)
	AddEvaluation(eval domain.Evaluation) error
	ListEvaluationsByProject(projectID string) ([]domain.Evaluation, error)
}

type hackathonService struct {
	repo        domain.HackathonRepository
	projectRepo domain.ProjectRepository
	evalRepo    domain.EvaluationRepository
}

func NewHackathonService(repo domain.HackathonRepository, projectRepo domain.ProjectRepository, evalRepo domain.EvaluationRepository) HackathonService {
	return &hackathonService{repo: repo, projectRepo: projectRepo, evalRepo: evalRepo}
}

func (s *hackathonService) ListHackathons() ([]domain.Hackathon, error) {
	return s.repo.GetAll()
}

func (s *hackathonService) ListProjectsByHackathon(id string) ([]domain.Project, error) {
	return s.projectRepo.GetByHackathonID(id)
}

func (s *hackathonService) ListEvaluationsByProject(projectID string) ([]domain.Evaluation, error) {
	return s.evalRepo.GetByProjectID(projectID)
}

func (s *hackathonService) AddEvaluation(eval domain.Evaluation) error {
	if err := s.evalRepo.Save(eval); err != nil {
		return err
	}

	evals, err := s.evalRepo.GetByProjectID(eval.ProjectID)
	if err != nil {
		return err
	}

	var total float64
	for _, e := range evals {
		total += e.TotalScore
	}
	avg := total / float64(len(evals))

	return s.projectRepo.UpdateScore(eval.ProjectID, avg)
}
