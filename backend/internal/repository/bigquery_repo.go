package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/civil"
	"github.com/moficodes/hackathon-judge/backend/internal/domain"
	"google.golang.org/api/iterator"
)

type BigQueryRepo struct {
	client    *bigquery.Client
	projectID string
}

func NewBigQueryRepo(client *bigquery.Client, projectID string) *BigQueryRepo {
	return &BigQueryRepo{
		client:    client,
		projectID: projectID,
	}
}

type bqHackathon struct {
	ID            string    `bigquery:"id"`
	Title         string    `bigquery:"title"`
	Date          time.Time `bigquery:"date"`
	Description   string    `bigquery:"description"`
	Goal          string    `bigquery:"goal"`
	Status        string    `bigquery:"status"`
	Criteria      string    `bigquery:"criteria"`
	BonusCriteria string    `bigquery:"bonus_criteria"`
}

func (r *BigQueryRepo) mapBQHackathon(bqH bqHackathon) (domain.Hackathon, error) {
	h := domain.Hackathon{
		ID:          bqH.ID,
		Title:       bqH.Title,
		Date:        bqH.Date,
		Description: bqH.Description,
		Goal:        bqH.Goal,
		Status:      bqH.Status,
	}

	if bqH.Criteria != "" {
		if err := json.Unmarshal([]byte(bqH.Criteria), &h.Criteria); err != nil {
			return h, fmt.Errorf("failed to unmarshal criteria: %w", err)
		}
	}
	if bqH.BonusCriteria != "" {
		if err := json.Unmarshal([]byte(bqH.BonusCriteria), &h.BonusCriteria); err != nil {
			return h, fmt.Errorf("failed to unmarshal bonus criteria: %w", err)
		}
	}
	return h, nil
}

func (r *BigQueryRepo) GetAll() ([]domain.Hackathon, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	query := r.client.Query(fmt.Sprintf("SELECT * FROM `%s.hackathons.hackathons`", r.projectID))
	it, err := query.Read(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to read hackathons: %w", err)
	}

	var hackathons []domain.Hackathon
	for {
		var row bqHackathon
		err := it.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to iterate hackathons: %w", err)
		}
		h, err := r.mapBQHackathon(row)
		if err != nil {
			return nil, fmt.Errorf("failed to map hackathon row: %w", err)
		}
		hackathons = append(hackathons, h)
	}
	return hackathons, nil
}

func (r *BigQueryRepo) GetByID(id string) (domain.Hackathon, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	query := r.client.Query(fmt.Sprintf("SELECT * FROM `%s.hackathons.hackathons` WHERE id = @id LIMIT 1", r.projectID))
	query.Parameters = []bigquery.QueryParameter{
		{Name: "id", Value: id},
	}
	it, err := query.Read(ctx)
	if err != nil {
		return domain.Hackathon{}, fmt.Errorf("failed to read hackathon: %w", err)
	}

	var row bqHackathon
	err = it.Next(&row)
	if err == iterator.Done {
		return domain.Hackathon{}, fmt.Errorf("hackathon not found")
	}
	if err != nil {
		return domain.Hackathon{}, fmt.Errorf("failed to iterate hackathon: %w", err)
	}

	return r.mapBQHackathon(row)
}

type bqProject struct {
	ID          string             `bigquery:"id"`
	Name        string             `bigquery:"name"`
	Title       string             `bigquery:"title"`
	URL         string             `bigquery:"url"`
	GitHubURL   string             `bigquery:"github_url"`
	TeamName    string             `bigquery:"team_name"`
	Document    bigquery.NullString `bigquery:"document"`
	Date        civil.Date         `bigquery:"date"`
	HackathonID string             `bigquery:"hackathon_id"`
	Score       float64            `bigquery:"score"`
}

func (r *BigQueryRepo) mapBQProject(bqP bqProject) domain.Project {
	doc := ""
	if bqP.Document.Valid {
		doc = bqP.Document.StringVal
	}

	return domain.Project{
		ID:          bqP.ID,
		Name:        bqP.Name,
		Title:       bqP.Title,
		URL:         bqP.URL,
		GitHubURL:   bqP.GitHubURL,
		TeamName:    bqP.TeamName,
		Document:    doc,
		Date:        bqP.Date.In(time.UTC),
		HackathonID: bqP.HackathonID,
		Score:       bqP.Score,
	}
}

func (r *BigQueryRepo) GetByHackathonID(hackathonID string) ([]domain.Project, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	query := r.client.Query(fmt.Sprintf("SELECT * FROM `%s.projects.projects` WHERE hackathon_id = @hackathon_id", r.projectID))
	query.Parameters = []bigquery.QueryParameter{
		{Name: "hackathon_id", Value: hackathonID},
	}
	it, err := query.Read(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to read projects: %w", err)
	}

	var projects []domain.Project
	for {
		var row bqProject
		err := it.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to iterate projects: %w", err)
		}
		projects = append(projects, r.mapBQProject(row))
	}
	return projects, nil
}

func (r *BigQueryRepo) GetProjectByID(id string) (domain.Project, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	query := r.client.Query(fmt.Sprintf("SELECT * FROM `%s.projects.projects` WHERE id = @id LIMIT 1", r.projectID))
	query.Parameters = []bigquery.QueryParameter{
		{Name: "id", Value: id},
	}
	it, err := query.Read(ctx)
	if err != nil {
		return domain.Project{}, fmt.Errorf("failed to read project: %w", err)
	}

	var row bqProject
	err = it.Next(&row)
	if err == iterator.Done {
		return domain.Project{}, fmt.Errorf("project not found")
	}
	if err != nil {
		return domain.Project{}, fmt.Errorf("failed to iterate project: %w", err)
	}

	return r.mapBQProject(row), nil
}

func (r *BigQueryRepo) UpdateScore(projectID string, score float64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	query := r.client.Query(fmt.Sprintf("UPDATE `%s.projects.projects` SET score = @score WHERE id = @id", r.projectID))
	query.Parameters = []bigquery.QueryParameter{
		{Name: "score", Value: score},
		{Name: "id", Value: projectID},
	}

	job, err := query.Run(ctx)
	if err != nil {
		return fmt.Errorf("failed to run update query: %w", err)
	}
	status, err := job.Wait(ctx)
	if err != nil {
		return fmt.Errorf("failed to wait for update query: %w", err)
	}
	if status.Err() != nil {
		return fmt.Errorf("update query failed: %w", status.Err())
	}
	return nil
}

type bqEvaluation struct {
	ID           string    `bigquery:"id"`
	ProjectID    string    `bigquery:"project_id"`
	JudgeID      string    `bigquery:"judge_id"`
	Status       string    `bigquery:"status"`
	TotalScore   float64   `bigquery:"total_score"`
	Comment      string    `bigquery:"comment"`
	CreatedAt    time.Time `bigquery:"created_at"`
	CriteriaJSON string    `bigquery:"criteria_json"`
}

func (r *BigQueryRepo) Save(eval domain.Evaluation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	criteriaBytes, err := json.Marshal(eval.Criteria)
	if err != nil {
		return fmt.Errorf("failed to marshal criteria: %w", err)
	}

	if eval.CreatedAt.IsZero() {
		eval.CreatedAt = time.Now()
	}

	query := r.client.Query(fmt.Sprintf("INSERT INTO `%s.evaluations.evaluations` (id, project_id, judge_id, status, total_score, comment, created_at, criteria_json) VALUES (@id, @project_id, @judge_id, @status, @total_score, @comment, @created_at, @criteria_json)", r.projectID))
	query.Parameters = []bigquery.QueryParameter{
		{Name: "id", Value: eval.ID},
		{Name: "project_id", Value: eval.ProjectID},
		{Name: "judge_id", Value: eval.JudgeID},
		{Name: "status", Value: eval.Status},
		{Name: "total_score", Value: eval.TotalScore},
		{Name: "comment", Value: eval.Comment},
		{Name: "created_at", Value: eval.CreatedAt},
		{Name: "criteria_json", Value: string(criteriaBytes)},
	}

	job, err := query.Run(ctx)
	if err != nil {
		return fmt.Errorf("failed to run insert query: %w", err)
	}
	status, err := job.Wait(ctx)
	if err != nil {
		return fmt.Errorf("failed to wait for insert query: %w", err)
	}
	if status.Err() != nil {
		return fmt.Errorf("insert query failed: %w", status.Err())
	}
	return nil
}

func (r *BigQueryRepo) Update(eval domain.Evaluation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	criteriaBytes, err := json.Marshal(eval.Criteria)
	if err != nil {
		return fmt.Errorf("failed to marshal criteria: %w", err)
	}

	query := r.client.Query(fmt.Sprintf(`
		UPDATE %s.evaluations.evaluations 
		SET status = @status, 
		    total_score = @total_score, 
		    comment = @comment, 
		    criteria_json = @criteria_json 
		WHERE id = @id`, r.projectID))

	query.Parameters = []bigquery.QueryParameter{
		{Name: "id", Value: eval.ID},
		{Name: "status", Value: eval.Status},
		{Name: "total_score", Value: eval.TotalScore},
		{Name: "comment", Value: eval.Comment},
		{Name: "criteria_json", Value: string(criteriaBytes)},
	}

	job, err := query.Run(ctx)
	if err != nil {
		return fmt.Errorf("failed to run update query: %w", err)
	}
	status, err := job.Wait(ctx)
	if err != nil {
		return fmt.Errorf("failed to wait for update query: %w", err)
	}
	if status.Err() != nil {
		return fmt.Errorf("update query failed: %w", status.Err())
	}
	return nil
}

func (r *BigQueryRepo) GetEvaluationByID(id string) (domain.Evaluation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := r.client.Query(fmt.Sprintf("SELECT id, project_id, judge_id, status, total_score, comment, created_at, criteria_json FROM `%s.evaluations.evaluations` WHERE id = @id LIMIT 1", r.projectID))
	query.Parameters = []bigquery.QueryParameter{
		{Name: "id", Value: id},
	}

	it, err := query.Read(ctx)
	if err != nil {
		return domain.Evaluation{}, fmt.Errorf("failed to read evaluation: %w", err)
	}

	var bqEval bqEvaluation
	err = it.Next(&bqEval)
	if err != nil {
		return domain.Evaluation{}, fmt.Errorf("evaluation not found: %w", err)
	}

	var criteria []domain.CriteriaScore
	if bqEval.CriteriaJSON != "" {
		if err := json.Unmarshal([]byte(bqEval.CriteriaJSON), &criteria); err != nil {
			return domain.Evaluation{}, fmt.Errorf("failed to unmarshal criteria: %w", err)
		}
	}

	return domain.Evaluation{
		ID:         bqEval.ID,
		ProjectID:  bqEval.ProjectID,
		JudgeID:    bqEval.JudgeID,
		Status:     bqEval.Status,
		TotalScore: bqEval.TotalScore,
		Comment:    bqEval.Comment,
		CreatedAt:  bqEval.CreatedAt,
		Criteria:   criteria,
	}, nil
}

func (r *BigQueryRepo) GetByProjectID(projectID string) ([]domain.Evaluation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := r.client.Query(fmt.Sprintf("SELECT id, project_id, judge_id, status, total_score, comment, created_at, criteria_json FROM `%s.evaluations.evaluations` WHERE project_id = @project_id ORDER BY created_at DESC", r.projectID))
	query.Parameters = []bigquery.QueryParameter{
		{Name: "project_id", Value: projectID},
	}

	it, err := query.Read(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to query evaluations: %w", err)
	}

	var evaluations []domain.Evaluation
	for {
		var bqEval bqEvaluation
		err := it.Next(&bqEval)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to iterate evaluations: %w", err)
		}

		var criteria []domain.CriteriaScore
		if bqEval.CriteriaJSON != "" {
			if err := json.Unmarshal([]byte(bqEval.CriteriaJSON), &criteria); err != nil {
				return nil, fmt.Errorf("failed to unmarshal criteria: %w", err)
			}
		}

		evaluations = append(evaluations, domain.Evaluation{
			ID:         bqEval.ID,
			ProjectID:  bqEval.ProjectID,
			JudgeID:    bqEval.JudgeID,
			Status:     bqEval.Status,
			TotalScore: bqEval.TotalScore,
			Comment:    bqEval.Comment,
			CreatedAt:  bqEval.CreatedAt,
			Criteria:   criteria,
		})
	}

	return evaluations, nil
}
