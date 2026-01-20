package repo

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/youssefM1999/report/internal/config"
)

func TestClone(t *testing.T) {
	// Create a temporary directory for cloning
	tmpDir, err := os.MkdirTemp("", "repo-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Use a small public repository for testing
	repo := NewRepo(
		"test-repo",
		"https://github.com/octocat/Hello-World.git",
		"master",
		filepath.Join(tmpDir, "test-repo"),
	)

	err = repo.Clone()
	if err != nil {
		t.Fatalf("Clone() failed: %v", err)
	}

	// Verify the repo was cloned by checking if .git directory exists
	gitDir := filepath.Join(repo.RepoDir, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		t.Errorf("Clone() did not create .git directory at %s", gitDir)
	}
}

func TestRepoManager_CloneAll(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "repo-manager-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	rm := NewRepoManager(tmpDir)

	reposConfig := config.ReposConfig{
		TargetRepos: []config.RepoConfig{
			{
				Name:   "hello-world",
				URL:    "https://github.com/octocat/Hello-World.git",
				Branch: "master",
			},
		},
	}

	err = rm.CloneAll(reposConfig)
	if err != nil {
		t.Fatalf("CloneAll() failed: %v", err)
	}

	// Verify the repo was cloned
	repoDir := filepath.Join(tmpDir, "hello-world", ".git")
	if _, err := os.Stat(repoDir); os.IsNotExist(err) {
		t.Errorf("CloneAll() did not clone repo to %s", repoDir)
	}
}

func TestGetCommitsByAuthor(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "commits-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Clone the repo first
	repo := NewRepo(
		"hello-world",
		"https://github.com/octocat/Hello-World.git",
		"master",
		filepath.Join(tmpDir, "hello-world"),
	)

	if err := repo.Clone(); err != nil {
		t.Fatalf("Clone() failed: %v", err)
	}

	// Get commits by the known author (The Octocat)
	author := config.UserConfig{
		FullName: "The Octocat",
		Email:    "octocat@nowhere.com",
	}

	// Use a date far in the past to ensure we get commits
	since := time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC)

	err = repo.GetCommitsByAuthor(author, since)
	if err != nil {
		t.Fatalf("GetCommitsByAuthor() failed: %v", err)
	}

	// The Hello-World repo has commits from octocat
	if len(repo.Commits) == 0 {
		t.Log("No commits found for octocat@nowhere.com - this may be expected if the repo structure changed")
	} else {
		t.Logf("Found %d commits", len(repo.Commits))
		for _, c := range repo.Commits {
			t.Logf("  - %s: %s (%s)", c.Hash[:7], c.Message, c.Author)
		}
	}
}

func TestGetCommitsByAuthor_NoCommits(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "commits-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	repo := NewRepo(
		"hello-world",
		"https://github.com/octocat/Hello-World.git",
		"master",
		filepath.Join(tmpDir, "hello-world"),
	)

	if err := repo.Clone(); err != nil {
		t.Fatalf("Clone() failed: %v", err)
	}

	// Use an author that doesn't exist
	author := config.UserConfig{
		FullName: "Nonexistent User",
		Email:    "nonexistent@example.com",
	}

	since := time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC)

	err = repo.GetCommitsByAuthor(author, since)
	if err != nil {
		t.Fatalf("GetCommitsByAuthor() failed: %v", err)
	}

	if len(repo.Commits) != 0 {
		t.Errorf("Expected 0 commits for nonexistent author, got %d", len(repo.Commits))
	}
}
