package report

import (
	"fmt"
	"strings"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/youssefM1999/report/internal/config"
	"github.com/youssefM1999/report/internal/repo"
)

const (
	dateFormat = "Jan 2, 2006"
)

type RepoCommits struct {
	RepoName string
	Commits  []repo.Commit
}

type Report struct {
	Author    config.UserConfig
	StartDate time.Time
	EndDate   time.Time
	Repos     []RepoCommits
}

func NewReport(author config.UserConfig, startDate, endDate time.Time) *Report {
	return &Report{
		Author:    author,
		StartDate: startDate,
		EndDate:   endDate,
		Repos:     []RepoCommits{},
	}
}

func (r *Report) AddRepoCommits(repoName string, commits []repo.Commit) {
	r.Repos = append(r.Repos, RepoCommits{
		RepoName: repoName,
		Commits:  commits,
	})
}

func (r *Report) ToMarkdown() string {
	var sb strings.Builder

	sb.WriteString("# Work Report\n\n")
	sb.WriteString(fmt.Sprintf("**Author:** %s (%s)\n\n", r.Author.FullName, r.Author.Email))
	sb.WriteString(fmt.Sprintf("**Period:** %s - %s\n\n",
		r.StartDate.Format(dateFormat),
		r.EndDate.Format(dateFormat)))
	sb.WriteString("---\n\n")

	totalCommits := 0
	for _, rc := range r.Repos {
		totalCommits += len(rc.Commits)
	}
	sb.WriteString(fmt.Sprintf("**Total Commits:** %d\n\n", totalCommits))

	for _, rc := range r.Repos {
		if len(rc.Commits) == 0 {
			continue
		}

		sb.WriteString(fmt.Sprintf("## %s\n\n", rc.RepoName))
		sb.WriteString(fmt.Sprintf("*%d commits*\n\n", len(rc.Commits)))

		for _, c := range rc.Commits {
			sb.WriteString(fmt.Sprintf("- **%s** - %s (`%s`)\n",
				c.Date.Format("Jan 2"),
				c.Message,
				c.Hash[:7],
			))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

func (r *Report) ToHTML() string {
	md := r.ToMarkdown()
	html := markdown.ToHTML([]byte(md), nil, nil)
	return string(html)
}
