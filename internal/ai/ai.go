package ai

import (
	"context"
	"fmt"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/youssefM1999/report/internal/repo"
)

const systemPrompt = `You are a technical writer creating concise developer activity reports. Your task is to summarize git commits into clear, actionable bullet points.

Guidelines:
- For each commit, provide 2-3 bullet points that explain WHAT was done and WHY it matters
- Focus on the impact and purpose, not implementation details
- Use clear, professional language suitable for stakeholders and team leads
- Group related commits if they are part of the same feature or fix
- Always include the commit hash (short form, 7 characters) at the start of each commit section
- Use markdown formatting

Output format for each commit:
### <commit hash> - <brief title>
- <bullet point 1: what was changed>
- <bullet point 2: why it matters or what problem it solves>
- <bullet point 3: any notable technical decisions (optional, only if relevant)>

If there are no commits, respond with "No commits in this period."`

const multiRepoSystemPrompt = `You are a technical writer creating concise developer activity reports. Your task is to summarize git commits from multiple repositories into a clear, organized report.

Guidelines:
- Organize the report by repository
- For each commit, provide 2-3 bullet points that explain WHAT was done and WHY it matters
- Focus on the impact and purpose, not implementation details
- Use clear, professional language suitable for stakeholders and team leads
- Always include the commit hash (short form, 7 characters) at the start of each commit section
- Use markdown formatting

Output format:
## <Repository Name>

### <commit hash> - <brief title>
- <bullet point 1: what was changed>
- <bullet point 2: why it matters or what problem it solves>
- <bullet point 3: any notable technical decisions (optional, only if relevant)>

If a repository has no commits, include a note: "No commits in this period."
If all repositories have no commits, respond with "No commits across any repositories in this period."`

const userPromptTemplate = `Generate a report summary for the following commits from repository "%s":

%s

Provide a clear, concise summary with 2-3 bullet points per commit. Include the commit hash for each.`

const multiRepoUserPromptTemplate = `Generate a developer activity report for the following commits across multiple repositories:

%s

Provide a clear, organized summary with 2-3 bullet points per commit. Group by repository and include the commit hash for each.`

type AI interface {
	GenerateRepoReport(repoName string, commits []*repo.Commit) (string, error)
	GenerateFullReport(repos []*repo.Repo) (string, error)
}

type ClaudeAI struct {
	client anthropic.Client
}

func NewClaudeAI(apiKey string) *ClaudeAI {
	client := anthropic.NewClient(
		option.WithAPIKey(apiKey),
	)
	return &ClaudeAI{
		client: client,
	}
}

func (c *ClaudeAI) GenerateRepoReport(repoName string, commits []*repo.Commit) (string, error) {
	if len(commits) == 0 {
		return "No commits in this period.", nil
	}

	commitsText := formatCommitsForPrompt(commits)
	userPrompt := fmt.Sprintf(userPromptTemplate, repoName, commitsText)

	message, err := c.client.Messages.New(context.Background(), anthropic.MessageNewParams{
		Model:     anthropic.ModelClaude3_5HaikuLatest,
		MaxTokens: 1024,
		System: []anthropic.TextBlockParam{
			{Text: systemPrompt},
		},
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(userPrompt)),
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate report: %w", err)
	}

	var result strings.Builder
	for _, block := range message.Content {
		if block.Type == "text" {
			result.WriteString(block.Text)
		}
	}

	return result.String(), nil
}

func (c *ClaudeAI) GenerateFullReport(repos []*repo.Repo) (string, error) {
	// Check if there are any commits across all repos
	totalCommits := 0
	for _, r := range repos {
		totalCommits += len(r.Commits)
	}
	if totalCommits == 0 {
		return "No commits across any repositories in this period.", nil
	}

	reposText := formatAllReposForPrompt(repos)
	userPrompt := fmt.Sprintf(multiRepoUserPromptTemplate, reposText)

	message, err := c.client.Messages.New(context.Background(), anthropic.MessageNewParams{
		Model:     anthropic.ModelClaude3_5HaikuLatest,
		MaxTokens: 2048,
		System: []anthropic.TextBlockParam{
			{Text: multiRepoSystemPrompt},
		},
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(userPrompt)),
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate report: %w", err)
	}

	var result strings.Builder
	for _, block := range message.Content {
		if block.Type == "text" {
			result.WriteString(block.Text)
		}
	}

	return result.String(), nil
}

func formatCommitsForPrompt(commits []*repo.Commit) string {
	var sb strings.Builder
	for _, c := range commits {
		sb.WriteString(fmt.Sprintf("Commit: %s\n", c.Hash))
		sb.WriteString(fmt.Sprintf("Author: %s\n", c.Author))
		sb.WriteString(fmt.Sprintf("Date: %s\n", c.Date.Format("2006-01-02 15:04:05")))
		sb.WriteString(fmt.Sprintf("Message: %s\n", c.Message))
		if len(c.Content) > 0 {
			sb.WriteString(fmt.Sprintf("Changes:\n%s\n", c.Content))
		}
		sb.WriteString("\n---\n\n")
	}
	return sb.String()
}

func formatAllReposForPrompt(repos []*repo.Repo) string {
	var sb strings.Builder
	for _, r := range repos {
		sb.WriteString(fmt.Sprintf("=== Repository: %s ===\n\n", r.Name))
		if len(r.Commits) == 0 {
			sb.WriteString("No commits in this period.\n\n")
			continue
		}
		sb.WriteString(formatCommitsForPrompt(r.Commits))
		sb.WriteString("\n")
	}
	return sb.String()
}
