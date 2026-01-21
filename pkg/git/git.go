package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func Clone(repoDir, url, branch string) error {
	// Check if the directory exists and is a git repository
	// if it is, pull the repository
	if IsGitRepository(repoDir) {
		return Pull(repoDir, branch)
	}

	// if it is not, clone the repository
	cmd := exec.Command("git", "clone", "--branch", branch, url, repoDir)
	return cmd.Run()
}

func GetCommitsByAuthor(repoDir, email string, since time.Time) ([]byte, error) {
	cmd := exec.Command("git", "-C", repoDir, "log",
		"--format=%H|||%an|||%at|||%s",
		"--since", since.Format(time.RFC3339),
		"--author", email,
	)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return output, nil
}

func GetCommitContents(repoDir, hash string) (string, error) {
	cmd := exec.Command("git", "-C", repoDir, "diff", hash+"^!", "--")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

func IsGitRepository(repoDir string) bool {
	gitDir := filepath.Join(repoDir, ".git")
	_, err := os.Stat(gitDir)
	return err == nil
}

func Pull(repoDir, branch string) error {
	cmd := exec.Command("git", "-C", repoDir, "pull", "origin", branch)
	return cmd.Run()
}
