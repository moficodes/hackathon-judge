package service

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/moficodes/hackathon-judge/backend/internal/domain"
)

type HackathonService interface {
	ListHackathons() ([]domain.Hackathon, error)
	GetHackathon(id string) (domain.Hackathon, error)
	ListProjectsByHackathon(id string) ([]domain.Project, error)
	GetProject(id string) (domain.Project, error)
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

func (s *hackathonService) GetHackathon(id string) (domain.Hackathon, error) {
	return s.repo.GetByID(id)
}

func (s *hackathonService) GetProject(id string) (domain.Project, error) {
	return s.projectRepo.GetProjectByID(id)
}

func (s *hackathonService) ListProjectsByHackathon(id string) ([]domain.Project, error) {
	return s.projectRepo.GetByHackathonID(id)
}

func (s *hackathonService) ListEvaluationsByProject(projectID string) ([]domain.Evaluation, error) {
	// Verify project exists
	if _, err := s.projectRepo.GetProjectByID(projectID); err != nil {
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	evals, err := s.evalRepo.GetByProjectID(projectID)
	if err != nil {
		return nil, err
	}

	// Ensure we return an empty array instead of null for JSON serialization
	if evals == nil {
		return []domain.Evaluation{}, nil
	}
	return evals, nil
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
		// Check standard criteria
		for _, hc := range hackathon.Criteria {
			if hc.Name == cs.Name {
				weight = hc.Weight
				break
			}
		}
		// If not found, check bonus criteria
		if weight == 0.0 {
			for _, bc := range hackathon.BonusCriteria {
				if bc.Name == cs.Name {
					weight = bc.Weight
					break
				}
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
	var scoringCriteria []domain.Criterion
	scoringCriteria = append(scoringCriteria, hackathon.Criteria...)
	scoringCriteria = append(scoringCriteria, hackathon.BonusCriteria...)

	taskID := "tsk_" + uuid.New().String()

	// Save RUNNING evaluation
	eval := domain.Evaluation{
		ID:        taskID,
		ProjectID: projectID,
		JudgeID:   "system-agent",
		Status:    "RUNNING",
		CreatedAt: time.Now(),
	}
	if err := s.evalRepo.Save(eval); err != nil {
		return "", fmt.Errorf("failed to save running evaluation: %w", err)
	}

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
