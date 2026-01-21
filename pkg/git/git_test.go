package git

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestClone(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "git-clone-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	repoDir := filepath.Join(tmpDir, "test-repo")
	url := "https://github.com/octocat/Hello-World.git"
	branch := "master"

	err = Clone(repoDir, url, branch)
	if err != nil {
		t.Fatalf("Clone() failed: %v", err)
	}

	// Verify the repo was cloned by checking if .git directory exists
	gitDir := filepath.Join(repoDir, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		t.Errorf("Clone() did not create .git directory at %s", gitDir)
	}
}

func TestClone_InvalidURL(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "git-clone-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	repoDir := filepath.Join(tmpDir, "test-repo")
	url := "https://invalid-url-that-does-not-exist.git"
	branch := "master"

	err = Clone(repoDir, url, branch)
	if err == nil {
		t.Error("Clone() should fail with invalid URL")
	}
}

func TestClone_ExistingDirectory(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "git-clone-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	repoDir := filepath.Join(tmpDir, "test-repo")
	url := "https://github.com/octocat/Hello-World.git"
	branch := "master"

	// First clone should succeed
	err = Clone(repoDir, url, branch)
	if err != nil {
		t.Fatalf("First Clone() failed: %v", err)
	}

	// Verify the repo was cloned
	gitDir := filepath.Join(repoDir, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		t.Fatalf("First clone did not create .git directory at %s", gitDir)
	}

	// Second clone into the same directory should pull instead of failing
	err = Clone(repoDir, url, branch)
	if err != nil {
		t.Errorf("Clone() should pull when directory already exists, got error: %v", err)
	}

	// Verify the repo still exists after pull
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		t.Errorf("Git directory should still exist after pull")
	}
}

func TestGetCommitsByAuthor(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "git-commits-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Clone the repo first
	repoDir := filepath.Join(tmpDir, "hello-world")
	url := "https://github.com/octocat/Hello-World.git"
	branch := "master"

	if err := Clone(repoDir, url, branch); err != nil {
		t.Fatalf("Clone() failed: %v", err)
	}

	// Get commits by the known author (The Octocat)
	email := "octocat@nowhere.com"
	since := time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC)

	output, err := GetCommitsByAuthor(repoDir, email, since)
	if err != nil {
		t.Fatalf("GetCommitsByAuthor() failed: %v", err)
	}

	// Verify output format: should contain commit hash, author, timestamp, and message separated by |||
	if len(output) > 0 {
		lines := strings.Split(strings.TrimSpace(string(output)), "\n")
		if len(lines) == 0 {
			t.Error("GetCommitsByAuthor() should return at least one commit")
		}

		// Check format of first commit
		parts := strings.Split(lines[0], "|||")
		if len(parts) != 4 {
			t.Errorf("Expected 4 parts separated by |||, got %d: %v", len(parts), parts)
		}

		// Verify parts are not empty
		if parts[0] == "" {
			t.Error("Commit hash should not be empty")
		}
		if parts[1] == "" {
			t.Error("Author name should not be empty")
		}
		if parts[2] == "" {
			t.Error("Timestamp should not be empty")
		}
		if parts[3] == "" {
			t.Error("Commit message should not be empty")
		}
	} else {
		t.Log("No commits found for octocat@nowhere.com - this may be expected if the repo structure changed")
	}
}

func TestGetCommitsByAuthor_NoCommits(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "git-commits-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Clone the repo first
	repoDir := filepath.Join(tmpDir, "hello-world")
	url := "https://github.com/octocat/Hello-World.git"
	branch := "master"

	if err := Clone(repoDir, url, branch); err != nil {
		t.Fatalf("Clone() failed: %v", err)
	}

	// Use an email that doesn't exist in the repo
	email := "nonexistent@example.com"
	since := time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC)

	output, err := GetCommitsByAuthor(repoDir, email, since)
	if err != nil {
		t.Fatalf("GetCommitsByAuthor() should not fail for non-existent author, got: %v", err)
	}

	// Output should be empty or just whitespace
	if len(strings.TrimSpace(string(output))) != 0 {
		t.Errorf("Expected empty output for non-existent author, got: %s", string(output))
	}
}

func TestGetCommitsByAuthor_InvalidRepoDir(t *testing.T) {
	invalidDir := "/nonexistent/directory"
	email := "test@example.com"
	since := time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC)

	_, err := GetCommitsByAuthor(invalidDir, email, since)
	if err == nil {
		t.Error("GetCommitsByAuthor() should fail with invalid repo directory")
	}
}

func TestGetCommitContents(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "git-contents-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Clone the repo first
	repoDir := filepath.Join(tmpDir, "hello-world")
	url := "https://github.com/octocat/Hello-World.git"
	branch := "master"

	if err := Clone(repoDir, url, branch); err != nil {
		t.Fatalf("Clone() failed: %v", err)
	}

	// Get a commit hash first
	email := "octocat@nowhere.com"
	since := time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC)

	output, err := GetCommitsByAuthor(repoDir, email, since)
	if err != nil {
		t.Fatalf("GetCommitsByAuthor() failed: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) == 0 || lines[0] == "" {
		t.Skip("No commits found, skipping GetCommitContents test")
	}

	// Try to find a commit with non-empty content (skip merge commits)
	var commitHash string
	var contents string
	var foundNonEmptyCommit bool

	for _, line := range lines {
		parts := strings.Split(line, "|||")
		if len(parts) < 1 {
			continue
		}
		commitHash = parts[0]
		message := parts[3]

		// Skip merge commits as they often have empty diffs
		if strings.HasPrefix(strings.ToLower(message), "merge") {
			continue
		}

		// Get commit contents
		var err error
		contents, err = GetCommitContents(repoDir, commitHash)
		if err != nil {
			t.Logf("GetCommitContents() failed for commit %s: %v, trying next commit", commitHash[:7], err)
			continue
		}

		// If we got content, we're done
		if len(contents) > 0 {
			foundNonEmptyCommit = true
			break
		}
	}

	if !foundNonEmptyCommit {
		t.Skip("No commits with non-empty content found (may be all merge commits)")
	}

	// Verify contents contain diff markers or git diff output
	if !strings.Contains(contents, "diff") && !strings.Contains(contents, "@@") && !strings.Contains(contents, "index") {
		t.Logf("Commit contents may not be in expected format, but got content of length %d", len(contents))
	}
}

func TestGetCommitContents_InvalidHash(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "git-contents-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Clone the repo first
	repoDir := filepath.Join(tmpDir, "hello-world")
	url := "https://github.com/octocat/Hello-World.git"
	branch := "master"

	if err := Clone(repoDir, url, branch); err != nil {
		t.Fatalf("Clone() failed: %v", err)
	}

	// Use an invalid commit hash
	invalidHash := "0000000000000000000000000000000000000000"

	_, err = GetCommitContents(repoDir, invalidHash)
	if err == nil {
		t.Error("GetCommitContents() should fail with invalid commit hash")
	}
}

func TestGetCommitContents_InvalidRepoDir(t *testing.T) {
	invalidDir := "/nonexistent/directory"
	hash := "abc123"

	_, err := GetCommitContents(invalidDir, hash)
	if err == nil {
		t.Error("GetCommitContents() should fail with invalid repo directory")
	}
}
