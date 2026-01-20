package repo

import (
	"path/filepath"
	"time"

	"github.com/youssefM1999/report/internal/config"
	"github.com/youssefM1999/report/pkg/git"
)

type RepoManager struct {
	baseDir string
	repos   []*Repo
}

type Repository interface {
	Clone() error
	GetCommits() error
}

type Repo struct {
	Name    string
	URL     string
	Branch  string
	RepoDir string
	Commits []*Commit
}

func NewRepoManager(baseDir string) *RepoManager {
	return &RepoManager{
		baseDir: baseDir,
		repos:   []*Repo{},
	}
}

func (rm *RepoManager) CloneAll(reposConfig config.ReposConfig) error {
	for _, repoConfig := range reposConfig.TargetRepos {
		repo := rm.NewRepoFromConfig(repoConfig)
		if err := repo.Clone(); err != nil {
			return err
		}
		rm.repos = append(rm.repos, repo)
	}
	return nil
}

func (rm *RepoManager) NewRepoFromConfig(config config.RepoConfig) *Repo {
	return NewRepo(config.Name, config.URL, config.Branch, filepath.Join(rm.baseDir, config.Name))
}

func NewRepo(name, url, branch, repoDir string) *Repo {
	return &Repo{
		Name:    name,
		URL:     url,
		Branch:  branch,
		RepoDir: repoDir,
	}
}

func (r *Repo) Clone() error {
	return git.Clone(r.RepoDir, r.URL, r.Branch)
}

func (r *Repo) GetCommitsByAuthor(author config.UserConfig, since time.Time) error {
	output, err := git.GetCommitsByAuthor(r.RepoDir, author.Email, since)
	if err != nil {
		return err
	}
	commits, err := parseToCommits(output)
	if err != nil {
		return err
	}
	r.Commits = commits
	return nil
}

func (r *Repo) GetCommitsContents() error {
	for _, commit := range r.Commits {
		content, err := git.GetCommitContents(r.RepoDir, commit.Hash)
		if err != nil {
			return err
		}
		commit.Content = content
	}
	return nil
}

func (rm *RepoManager) Repos() []*Repo {
	return rm.repos
}

func (rm *RepoManager) GetAllCommitsByAuthor(author config.UserConfig, since time.Time) error {
	for _, repo := range rm.repos {
		if err := repo.GetCommitsByAuthor(author, since); err != nil {
			return err
		}
		if err := repo.GetCommitsContents(); err != nil {
			return err
		}
	}
	return nil
}
