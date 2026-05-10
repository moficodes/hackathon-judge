package service

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/moficodes/hackathon-judge/backend/internal/domain"
)

type HackathonService interface {
	ListHackathons() ([]domain.Hackathon, error)
	ListProjectsByHackathon(id string) ([]domain.Project, error)
	AddEvaluation(eval domain.Evaluation) error
	ListEvaluationsByProject(projectID string) ([]domain.Evaluation, error)
	TriggerJudging(projectID string) (string, error)
}

type hackathonService struct {
	repo        domain.HackathonRepository
	projectRepo domain.ProjectRepository
	evalRepo    domain.EvaluationRepository
	publisher   domain.TaskPublisher
}

func NewHackathonService(repo domain.HackathonRepository, projectRepo domain.ProjectRepository, evalRepo domain.EvaluationRepository, publisher domain.TaskPublisher) HackathonService {
	return &hackathonService{repo: repo, projectRepo: projectRepo, evalRepo: evalRepo, publisher: publisher}
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

func (s *hackathonService) TriggerJudging(projectID string) (string, error) {
	// Fetch project
	project, err := s.projectRepo.GetProjectByID(projectID)
	if err != nil {
		return "", fmt.Errorf("failed to get project: %w", err)
	}

	// Fetch hackathon for criteria
	hackathon, err := s.repo.GetByID(project.HackathonID)
	if err != nil {
		return "", fmt.Errorf("failed to get hackathon: %w", err)
	}

	// Create judging criteria mapping
	var scoringCriteria []domain.ScoringCriteria
	for _, c := range hackathon.Criteria {
		scoringCriteria = append(scoringCriteria, domain.ScoringCriteria{
			Name:     c.Name,
			Weight:   c.Weight,
			MaxScore: 10.0, // Assuming 10 for now
		})
	}

	taskID := "tsk_" + uuid.New().String()

	task := domain.JudgingTask{
		TaskID:          taskID,
		ProjectName:     project.Name,
		GithubURL:       project.GitHubURL,
		SubmissionText:  project.Document,
		JudgingRubric:   hackathon.Goal + "\n" + hackathon.Description,
		ScoringCriteria: scoringCriteria,
	}

	if s.publisher != nil {
		if err := s.publisher.PublishTask(task); err != nil {
			return "", fmt.Errorf("failed to publish task: %w", err)
		}
	} else {
		// Mock handling when publisher is nil (e.g. for simple tests)
		fmt.Printf("Mock published task %s for project %s\n", taskID, project.Name)
	}

	return taskID, nil
}
