package repo

import (
	"strconv"
	"strings"
	"time"
)

type Commit struct {
	Hash    string
	Message string
	Author  string
	Date    time.Time
	Content string
}

func parseToCommits(output []byte) ([]*Commit, error) {
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) == 0 || lines[0] == "" {
		return nil, nil
	}

	commits := make([]*Commit, 0, len(lines))
	for _, line := range lines {
		parts := strings.Split(line, "|||")
		if len(parts) != 4 {
			continue
		}

		timestamp, err := strconv.ParseInt(parts[2], 10, 64)
		if err != nil {
			return nil, err
		}

		commits = append(commits, &Commit{
			Hash:    parts[0],
			Author:  parts[1],
			Date:    time.Unix(timestamp, 0),
			Message: parts[3],
		})
	}

	return commits, nil
}
