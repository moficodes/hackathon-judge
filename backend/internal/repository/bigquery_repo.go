package repository

import (
	"context"
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
	datasetID string
}

func NewBigQueryRepo(client *bigquery.Client, projectID string, datasetID string) *BigQueryRepo {
	return &BigQueryRepo{
		client:    client,
		projectID: projectID,
		datasetID: datasetID,
	}
}

type bqCriterion struct {
	ID          string  `bigquery:"id"`
	Name        string  `bigquery:"name"`
	Description string  `bigquery:"description"`
	Weight      float64 `bigquery:"weight"`
	Score       float64 `bigquery:"score"`
	MaxScore    float64 `bigquery:"max_score"`
}

type bqHackathon struct {
	ID            string        `bigquery:"id"`
	Title         string        `bigquery:"title"`
	Date          time.Time     `bigquery:"date"`
	Description   string        `bigquery:"description"`
	Goal          string        `bigquery:"goal"`
	Status        string        `bigquery:"status"`
	Criteria      []bqCriterion `bigquery:"criteria"`
	BonusCriteria []bqCriterion `bigquery:"bonus_criteria"`
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

	h.Criteria = make([]domain.Criterion, len(bqH.Criteria))
	for i, c := range bqH.Criteria {
		h.Criteria[i] = domain.Criterion{
			Name:        c.Name,
			Weight:      c.Weight,
			Description: c.Description,
			MaxScore:    c.MaxScore,
		}
	}

	h.BonusCriteria = make([]domain.Criterion, len(bqH.BonusCriteria))
	for i, c := range bqH.BonusCriteria {
		h.BonusCriteria[i] = domain.Criterion{
			Name:        c.Name,
			Weight:      c.Weight,
			Description: c.Description,
			MaxScore:    c.MaxScore,
		}
	}
	return h, nil
}

func (r *BigQueryRepo) GetAll() ([]domain.Hackathon, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	query := r.client.Query(fmt.Sprintf("SELECT * FROM `%s.%s.hackathons`", r.projectID, r.datasetID))
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
	query := r.client.Query(fmt.Sprintf("SELECT * FROM `%s.%s.hackathons` WHERE id = @id LIMIT 1", r.projectID, r.datasetID))
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

type bqCriteriaScore struct {
	ID          string  `bigquery:"id"`
	Name        string  `bigquery:"name"`
	Description string  `bigquery:"description"`
	Weight      float64 `bigquery:"weight"`
	Score       float64 `bigquery:"score"`
	MaxScore    float64 `bigquery:"max_score"`
}

type bqNestedEvaluation struct {
	ID         string            `bigquery:"id"`
	JudgeID    string            `bigquery:"judge_id"`
	Status     string            `bigquery:"status"`
	Criteria   []bqCriteriaScore `bigquery:"criteria"`
	TotalScore float64           `bigquery:"total_score"`
	Comment    string            `bigquery:"comment"`
	CreatedAt  time.Time         `bigquery:"created_at"`
}

type bqProject struct {
	ID          string               `bigquery:"id"`
	Name        string               `bigquery:"name"`
	Title       string               `bigquery:"title"`
	URL         string               `bigquery:"url"`
	GitHubURL   string               `bigquery:"github_url"`
	TeamName    string               `bigquery:"team_name"`
	Document    bigquery.NullString  `bigquery:"document"`
	Date        civil.Date           `bigquery:"date"`
	HackathonID string               `bigquery:"hackathon_id"`
	Score       float64              `bigquery:"score"`
	Evaluations []bqNestedEvaluation `bigquery:"evaluations"`
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

func mapBQNestedEvaluation(bqEval bqNestedEvaluation, projectID string) domain.Evaluation {
	criteria := make([]domain.CriteriaScore, len(bqEval.Criteria))
	for i, c := range bqEval.Criteria {
		criteria[i] = domain.CriteriaScore{
			Name:      c.Name,
			Score:     c.Score,
			Reasoning: c.Description,
			MaxScore:  c.MaxScore,
			Weight:    c.Weight,
		}
	}

	return domain.Evaluation{
		ID:         bqEval.ID,
		ProjectID:  projectID,
		JudgeID:    bqEval.JudgeID,
		Status:     bqEval.Status,
		Criteria:   criteria,
		TotalScore: bqEval.TotalScore,
		Comment:    bqEval.Comment,
		CreatedAt:  bqEval.CreatedAt,
	}
}

func (r *BigQueryRepo) GetByHackathonID(hackathonID string) ([]domain.Project, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	query := r.client.Query(fmt.Sprintf("SELECT * FROM `%s.%s.projects` WHERE hackathon_id = @hackathon_id", r.projectID, r.datasetID))
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
	query := r.client.Query(fmt.Sprintf("SELECT * FROM `%s.%s.projects` WHERE id = @id LIMIT 1", r.projectID, r.datasetID))
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
	query := r.client.Query(fmt.Sprintf("UPDATE `%s.%s.projects` SET score = @score WHERE id = @id", r.projectID, r.datasetID))
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

type bqEvaluationRow struct {
	ProjectID  string            `bigquery:"project_id"`
	ID         string            `bigquery:"id"`
	JudgeID    string            `bigquery:"judge_id"`
	Status     string            `bigquery:"status"`
	Criteria   []bqCriteriaScore `bigquery:"criteria"`
	TotalScore float64           `bigquery:"total_score"`
	Comment    string            `bigquery:"comment"`
	CreatedAt  time.Time         `bigquery:"created_at"`
}

type bqProjectEvaluations struct {
	Evaluations []bqNestedEvaluation `bigquery:"evaluations"`
}

func (r *BigQueryRepo) Save(eval domain.Evaluation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if eval.CreatedAt.IsZero() {
		eval.CreatedAt = time.Now()
	}

	bqEval := bqNestedEvaluation{
		ID:         eval.ID,
		JudgeID:    eval.JudgeID,
		Status:     eval.Status,
		TotalScore: eval.TotalScore,
		Comment:    eval.Comment,
		CreatedAt:  eval.CreatedAt,
	}

	bqEval.Criteria = make([]bqCriteriaScore, len(eval.Criteria))
	for i, c := range eval.Criteria {
		bqEval.Criteria[i] = bqCriteriaScore{
			Name:        c.Name,
			Score:       c.Score,
			Description: c.Reasoning,
			Weight:      c.Weight,
			MaxScore:    c.MaxScore,
		}
	}

	query := r.client.Query(fmt.Sprintf("UPDATE `%s.%s.projects` SET evaluations = ARRAY_CONCAT(evaluations, [@new_eval]) WHERE id = @project_id", r.projectID, r.datasetID))
	query.Parameters = []bigquery.QueryParameter{
		{Name: "new_eval", Value: bqEval},
		{Name: "project_id", Value: eval.ProjectID},
	}

	job, err := query.Run(ctx)
	if err != nil {
		return fmt.Errorf("failed to run insert (append) query: %w", err)
	}
	status, err := job.Wait(ctx)
	if err != nil {
		return fmt.Errorf("failed to wait for insert (append) query: %w", err)
	}
	if status.Err() != nil {
		return fmt.Errorf("insert (append) query failed: %w", status.Err())
	}
	return nil
}

func (r *BigQueryRepo) Update(eval domain.Evaluation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	bqEval := bqNestedEvaluation{
		ID:         eval.ID,
		JudgeID:    eval.JudgeID,
		Status:     eval.Status,
		TotalScore: eval.TotalScore,
		Comment:    eval.Comment,
		CreatedAt:  eval.CreatedAt,
	}

	bqEval.Criteria = make([]bqCriteriaScore, len(eval.Criteria))
	for i, c := range eval.Criteria {
		bqEval.Criteria[i] = bqCriteriaScore{
			Name:        c.Name,
			Score:       c.Score,
			Description: c.Reasoning,
			Weight:      c.Weight,
			MaxScore:    c.MaxScore,
		}
	}

	query := r.client.Query(fmt.Sprintf(`
		UPDATE %s.%s.projects 
		SET evaluations = ARRAY(
			SELECT AS STRUCT 
				e.id,
				IF(e.id = @eval_id, @updated_eval.judge_id, e.judge_id) as judge_id,
				IF(e.id = @eval_id, @updated_eval.status, e.status) as status,
				IF(e.id = @eval_id, @updated_eval.criteria, e.criteria) as criteria,
				IF(e.id = @eval_id, @updated_eval.total_score, e.total_score) as total_score,
				IF(e.id = @eval_id, @updated_eval.comment, e.comment) as comment,
				IF(e.id = @eval_id, @updated_eval.created_at, e.created_at) as created_at
			FROM UNNEST(evaluations) e
		)
		WHERE id = @project_id`, r.projectID, r.datasetID))

	query.Parameters = []bigquery.QueryParameter{
		{Name: "eval_id", Value: eval.ID},
		{Name: "updated_eval", Value: bqEval},
		{Name: "project_id", Value: eval.ProjectID},
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

	query := r.client.Query(fmt.Sprintf("SELECT p.id as project_id, e.id, e.judge_id, e.status, e.criteria, e.total_score, e.comment, e.created_at FROM `%s.%s.projects` p, UNNEST(p.evaluations) e WHERE e.id = @eval_id LIMIT 1", r.projectID, r.datasetID))
	query.Parameters = []bigquery.QueryParameter{
		{Name: "eval_id", Value: id},
	}

	it, err := query.Read(ctx)
	if err != nil {
		return domain.Evaluation{}, fmt.Errorf("failed to read evaluation: %w", err)
	}

	var row bqEvaluationRow
	err = it.Next(&row)
	if err == iterator.Done {
		return domain.Evaluation{}, fmt.Errorf("evaluation not found")
	}
	if err != nil {
		return domain.Evaluation{}, fmt.Errorf("failed to iterate evaluation: %w", err)
	}

	criteria := make([]domain.CriteriaScore, len(row.Criteria))
	for i, c := range row.Criteria {
		criteria[i] = domain.CriteriaScore{
			Name:      c.Name,
			Score:     c.Score,
			Reasoning: c.Description,
			MaxScore:  c.MaxScore,
			Weight:    c.Weight,
		}
	}

	return domain.Evaluation{
		ID:         row.ID,
		ProjectID:  row.ProjectID,
		JudgeID:    row.JudgeID,
		Status:     row.Status,
		TotalScore: row.TotalScore,
		Comment:    row.Comment,
		CreatedAt:  row.CreatedAt,
		Criteria:   criteria,
	}, nil
}

func (r *BigQueryRepo) GetByProjectID(projectID string) ([]domain.Evaluation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := r.client.Query(fmt.Sprintf("SELECT evaluations FROM `%s.%s.projects` WHERE id = @project_id LIMIT 1", r.projectID, r.datasetID))
	query.Parameters = []bigquery.QueryParameter{
		{Name: "project_id", Value: projectID},
	}

	it, err := query.Read(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to query evaluations: %w", err)
	}

	var row bqProjectEvaluations
	err = it.Next(&row)
	if err == iterator.Done {
		return []domain.Evaluation{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to iterate evaluations: %w", err)
	}

	evaluations := make([]domain.Evaluation, len(row.Evaluations))
	for i, e := range row.Evaluations {
		evaluations[i] = mapBQNestedEvaluation(e, projectID)
	}

	// Sort by CreatedAt DESC (to mimic old DB order)
	for i := 0; i < len(evaluations)/2; i++ {
		j := len(evaluations) - i - 1
		evaluations[i], evaluations[j] = evaluations[j], evaluations[i]
	}

	return evaluations, nil
}
