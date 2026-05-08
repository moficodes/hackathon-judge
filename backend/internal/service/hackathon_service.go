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
	// Fetch Project to get HackathonID
	project, err := s.projectRepo.GetProjectByID(eval.ProjectID)
	if err != nil {
		return err
	}

	// Fetch Hackathon to get weights
	hackathon, err := s.repo.GetByID(project.HackathonID)
	if err != nil {
		return err
	}

	// Calculate TotalScore
	var evalTotal float64
	for _, cs := range eval.Criteria {
		weight := 0.0
		for _, hc := range hackathon.Criteria {
			if hc.Name == cs.Name {
				weight = hc.Weight
				break
			}
		}
		evalTotal += cs.Score * weight
	}
	eval.TotalScore = evalTotal

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
	average := total / float64(len(evals))

	return s.projectRepo.UpdateScore(eval.ProjectID, average)
}
