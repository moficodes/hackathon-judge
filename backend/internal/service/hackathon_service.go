// Copyright 2026 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package service

import (
	"fmt"
	"sync"
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
	TriggerJudging(projectID string, useBQ bool) (string, error)
	TriggerJudgingAgent(projectID string) (string, error)
}

type hackathonService struct {
	repo        domain.HackathonRepository
	projectRepo domain.ProjectRepository
	evalRepo    domain.EvaluationRepository
	publisher   domain.TaskPublisher
	scoreMu     sync.Mutex
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

	s.scoreMu.Lock()
	defer s.scoreMu.Unlock()

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

func (s *hackathonService) TriggerJudging(projectID string, useBQ bool) (string, error) {
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

	judgeID := "sandbox"
	isBQ := false

	type bqScorer interface {
		JudgeProjectWithBQ(projectID string, evaluationID string) ([]domain.CriteriaScore, error)
	}

	var bqRepo bqScorer
	if useBQ {
		if r, ok := s.evalRepo.(bqScorer); ok {
			bqRepo = r
			judgeID = "BQ AI Function"
			isBQ = true
		}
	}

	// Save RUNNING evaluation
	eval := domain.Evaluation{
		ID:        taskID,
		ProjectID: projectID,
		JudgeID:   judgeID,
		Status:    "RUNNING",
		CreatedAt: time.Now(),
	}
	if err := s.evalRepo.Save(eval); err != nil {
		return "", fmt.Errorf("failed to save running evaluation: %w", err)
	}

	if isBQ {
		go func() {
			scores, err := bqRepo.JudgeProjectWithBQ(projectID, taskID)

			eval, getErr := s.evalRepo.GetEvaluationByID(taskID)
			if getErr != nil {
				fmt.Printf("Error: failed to fetch active evaluation task %s: %v\n", taskID, getErr)
				return
			}

			if err != nil {
				eval.Status = "FAILED"
				eval.Comment = fmt.Sprintf("BigQuery AI scoring failed: %v", err)
			} else {
				eval.Status = "SUCCESS"
				eval.Criteria = scores

				// Calculate average score based on criteria weights
				var evalTotal float64
				for _, cs := range scores {
					evalTotal += cs.Score * cs.Weight
				}
				eval.TotalScore = evalTotal
				eval.Comment = "Automated BigQuery AI evaluation completed successfully"
			}

			if updateErr := s.evalRepo.Update(eval); updateErr != nil {
				fmt.Printf("Error: failed to update evaluation %s: %v\n", taskID, updateErr)
				return
			}

			if eval.Status == "SUCCESS" {
				s.scoreMu.Lock()
				evals, err := s.evalRepo.GetByProjectID(projectID)
				if err == nil {
					var total float64
					var count int
					for _, e := range evals {
						if e.Status == "SUCCESS" {
							total += e.TotalScore
							count++
						}
					}
					if count > 0 {
						average := total / float64(count)
						s.projectRepo.UpdateScore(projectID, average)
					}
				}
				s.scoreMu.Unlock()
			}
		}()

		return taskID, nil
	}

	// Pub/Sub tasks publishing
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

func (s *hackathonService) TriggerJudgingAgent(projectID string) (string, error) {
	return s.TriggerJudging(projectID, false)
}
