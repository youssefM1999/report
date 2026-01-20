package git

import (
	"os/exec"
	"time"
)

func Clone(repoDir, url, branch string) error {
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
