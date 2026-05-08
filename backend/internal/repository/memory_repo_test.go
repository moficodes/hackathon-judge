package repository_test

import (
	"github.com/moficodes/hackathon-judge/backend/internal/domain"
	"github.com/moficodes/hackathon-judge/backend/internal/repository"
	"sync"
	"testing"
)

func TestMemoryRepo_Concurrency(t *testing.T) {
	repo := repository.NewMemoryRepo()

	var wg sync.WaitGroup
	// Concurrently read and write to verify thread safety
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			repo.GetByHackathonID("1")
			repo.UpdateScore("p1", float64(i))
			repo.Save(domain.Evaluation{ProjectID: "p1", TotalScore: float64(i)})
			repo.GetAll()
			repo.GetByID("1")
			repo.GetByProjectID("p1")
		}(i)
	}
	wg.Wait()
}

func TestMemoryRepo_SaveDoesNotUpdateScore(t *testing.T) {
	repo := repository.NewMemoryRepo()

	// Ensure saving evaluation doesn't modify project score automatically anymore
	eval := domain.Evaluation{ProjectID: "p1", TotalScore: 100.0}
	repo.Save(eval)

	projects, _ := repo.GetByHackathonID("1")
	var p1 *domain.Project
	for _, p := range projects {
		if p.ID == "p1" {
			p1 = &p
			break
		}
	}

	if p1 == nil {
		t.Fatalf("project p1 not found")
	}

	// Original score was 0.0, shouldn't change to 100.0 based solely on save.
	if p1.Score != 0.0 {
		t.Errorf("Expected score to remain 0.0, got %f", p1.Score)
	}
}

func TestMemoryRepo_UpdateScore(t *testing.T) {
	repo := repository.NewMemoryRepo()

	err := repo.UpdateScore("p1", 99.5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	projects, _ := repo.GetByHackathonID("1")
	var p1 *domain.Project
	for _, p := range projects {
		if p.ID == "p1" {
			p1 = &p
			break
		}
	}

	if p1.Score != 99.5 {
		t.Errorf("Expected score to be 99.5, got %f", p1.Score)
	}
}
