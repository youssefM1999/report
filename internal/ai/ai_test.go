package ai

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/youssefM1999/report/internal/config"
	"github.com/youssefM1999/report/internal/env"
	"github.com/youssefM1999/report/internal/repo"
)

func TestFullReportGenerationFlow(t *testing.T) {
	// Skip if no API key is set
	if err := env.LoadEnvFile(); err != nil {
		t.Fatalf("failed to load env file: %v", err)
	}
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		t.Skip("ANTHROPIC_API_KEY not set, skipping integration test")
	}

	// Step 1: Create temp directory and clone repo
	tmpDir, err := os.MkdirTemp("", "report-flow-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	t.Log("Step 1: Cloning repository...")
	r := repo.NewRepo(
		"hello-world",
		"https://github.com/octocat/Hello-World.git",
		"master",
		filepath.Join(tmpDir, "hello-world"),
	)

	if err := r.Clone(); err != nil {
		t.Fatalf("Clone() failed: %v", err)
	}
	t.Log("Repository cloned successfully")

	// Step 2: Get commits by author
	t.Log("Step 2: Getting commits by author...")
	author := config.UserConfig{
		FullName: "The Octocat",
		Email:    "octocat@nowhere.com",
	}
	since := time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC)

	if err := r.GetCommitsByAuthor(author, since); err != nil {
		t.Fatalf("GetCommitsByAuthor() failed: %v", err)
	}
	t.Logf("Found %d commits", len(r.Commits))

	if len(r.Commits) == 0 {
		t.Log("No commits found, skipping content and AI generation")
		return
	}

	// Step 3: Get commit contents
	t.Log("Step 3: Getting commit contents...")
	if err := r.GetCommitsContents(); err != nil {
		t.Fatalf("GetCommitsContents() failed: %v", err)
	}

	for _, c := range r.Commits {
		t.Logf("  Commit %s: %s (content length: %d)", c.Hash[:7], c.Message, len(c.Content))
	}

	// Step 4: Generate AI report
	t.Log("Step 4: Generating AI report...")
	ai := NewClaudeAI(apiKey)

	report, err := ai.GenerateRepoReport(r.Name, r.Commits)
	if err != nil {
		t.Fatalf("GenerateRepoReport() failed: %v", err)
	}

	t.Log("Generated report:")
	t.Log("---")
	t.Log(report)
	t.Log("---")

	// Verify report contains expected elements
	if len(report) == 0 {
		t.Error("Report is empty")
	}
}

func TestMultiRepoReportGenerationFlow(t *testing.T) {
	if err := env.LoadEnvFile(); err != nil {
		t.Fatalf("failed to load env file: %v", err)
	}
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		t.Skip("ANTHROPIC_API_KEY not set, skipping integration test")
	}

	tmpDir, err := os.MkdirTemp("", "multi-repo-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Setup RepoManager with multiple repos
	rm := repo.NewRepoManager(tmpDir)

	reposConfig := config.ReposConfig{
		TargetRepos: []config.RepoConfig{
			{
				Name:   "report",
				URL:    "https://github.com/youssefM1999/report.git",
				Branch: "main",
			},
			{
				Name:   "social",
				URL:    "https://github.com/youssefM1999/social.git",
				Branch: "main",
			},
		},
	}

	t.Log("Step 1: Cloning all repositories...")
	if err := rm.CloneAll(reposConfig); err != nil {
		t.Fatalf("CloneAll() failed: %v", err)
	}
	t.Logf("Cloned %d repositories", len(rm.Repos()))

	// Get commits for all repos
	t.Log("Step 2: Getting commits for all repositories...")
	author := config.UserConfig{
		FullName: "Youssef Mahmoud",
		Email:    "kamel.youssef1@gmail.com",
	}
	since := time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC)

	if err := rm.GetAllCommitsByAuthor(author, since); err != nil {
		t.Fatalf("GetAllCommitsByAuthor() failed: %v", err)
	}

	for _, r := range rm.Repos() {
		t.Logf("  %s: %d commits", r.Name, len(r.Commits))
	}

	// Generate full report
	t.Log("Step 3: Generating full report...")
	ai := NewClaudeAI(apiKey)

	report, err := ai.GenerateFullReport(rm.Repos())
	if err != nil {
		t.Fatalf("GenerateFullReport() failed: %v", err)
	}

	t.Log("Generated multi-repo report:")
	t.Log("---")
	t.Log(report)
	t.Log("---")

	if len(report) == 0 {
		t.Error("Report is empty")
	}
}

func TestGenerateRepoReport_NoCommits(t *testing.T) {
	ai := NewClaudeAI("fake-key") // won't make API call for empty commits

	report, err := ai.GenerateRepoReport("test-repo", []*repo.Commit{})
	if err != nil {
		t.Fatalf("GenerateRepoReport() failed: %v", err)
	}

	expected := "No commits in this period."
	if report != expected {
		t.Errorf("Expected %q, got %q", expected, report)
	}
}

func TestGenerateFullReport_NoCommits(t *testing.T) {
	ai := NewClaudeAI("fake-key")

	repos := []*repo.Repo{
		{Name: "repo1", Commits: []*repo.Commit{}},
		{Name: "repo2", Commits: []*repo.Commit{}},
	}

	report, err := ai.GenerateFullReport(repos)
	if err != nil {
		t.Fatalf("GenerateFullReport() failed: %v", err)
	}

	expected := "No commits across any repositories in this period."
	if report != expected {
		t.Errorf("Expected %q, got %q", expected, report)
	}
}

func TestFormatCommitsForPrompt(t *testing.T) {
	commits := []*repo.Commit{
		{
			Hash:    "abc1234567890",
			Author:  "Test Author",
			Date:    time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
			Message: "Add new feature",
			Content: "diff --git a/file.go\n+new line",
		},
		{
			Hash:    "def9876543210",
			Author:  "Test Author",
			Date:    time.Date(2024, 1, 16, 14, 0, 0, 0, time.UTC),
			Message: "Fix bug",
			Content: "",
		},
	}

	result := formatCommitsForPrompt(commits)

	// Verify it contains expected data
	if !contains(result, "abc1234567890") {
		t.Error("Result should contain first commit hash")
	}
	if !contains(result, "def9876543210") {
		t.Error("Result should contain second commit hash")
	}
	if !contains(result, "Add new feature") {
		t.Error("Result should contain first commit message")
	}
	if !contains(result, "Fix bug") {
		t.Error("Result should contain second commit message")
	}
	if !contains(result, "diff --git") {
		t.Error("Result should contain diff content")
	}

	t.Log("Formatted output:")
	t.Log(result)
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
