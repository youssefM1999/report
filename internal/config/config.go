package config

import (
	"os"
	"path/filepath"
	"time"

	"github.com/goccy/go-yaml"
	"github.com/youssefM1999/report/internal/env"
	"github.com/youssefM1999/report/pkg/filesystem"
)

type Config struct {
	Env    string
	Mail   MailConfig
	Repos  ReposConfig
	Logger LoggerConfig
	AI     AIConfig
	User   UserConfig
	Range  time.Duration
}

type UserConfig struct {
	FullName string `yaml:"full_name"`
	Email    string `yaml:"email"`
}

type MailConfig struct {
	FromEmail string
}

type ReposConfig struct {
	YamlFilePath string //path to the yaml definition file
	Dir          string
	TargetRepos  []RepoConfig
}

type RepoConfig struct {
	Name   string `yaml:"name"`
	URL    string `yaml:"url"`
	Branch string `yaml:"branch"`
}

type LoggerConfig struct {
	Dir      string
	FilePath string
}

type AIConfig struct {
	Key string
}

// yamlFileConfig represents the structure of the YAML configuration file
type yamlFileConfig struct {
	User  UserConfig   `yaml:"user"`
	Repos []RepoConfig `yaml:"repos"`
}

func Load() (Config, error) {
	e := env.GetString("ENV", "development")
	if e == "development" {
		if err := env.LoadEnvFile(); err != nil {
			return Config{}, err
		}
	}

	logDir := env.GetString("LOG_DIR", "logs")
	logDir, err := filesystem.ResolvePath(logDir)
	if err != nil {
		return Config{}, err
	}

	aiKey := env.GetString("ANTHROPIC_API_KEY", "")

	repoDir := env.GetString("REPO_DIR", "repos")
	repoDir, err = filesystem.ResolvePath(repoDir)
	if err != nil {
		return Config{}, err
	}

	yamlFilePath := env.GetString("YAML_FILE_PATH", "config.yaml")
	yamlConfig, err := loadYAMLConfig(yamlFilePath)
	if err != nil {
		return Config{}, err
	}

	config := Config{
		Env: e,
		Logger: LoggerConfig{
			Dir: logDir,
		},
		AI: AIConfig{
			Key: aiKey,
		},
		Repos: ReposConfig{
			Dir:          repoDir,
			TargetRepos:  yamlConfig.Repos,
			YamlFilePath: yamlFilePath,
		},
		User:  yamlConfig.User,
		Range: env.GetDuration("REPORT_RANGE", 7*24*time.Hour),
	}

	return config, nil
}

func loadYAMLConfig(yamlFilePath string) (yamlFileConfig, error) {
	rootDir, err := filesystem.FindModuleRoot()
	if err != nil {
		return yamlFileConfig{}, err
	}
	yamlFile := filepath.Join(rootDir, yamlFilePath)
	if err := filesystem.CheckValidFile(yamlFile); err != nil {
		return yamlFileConfig{}, err
	}
	yamlBytes, err := os.ReadFile(yamlFile)
	if err != nil {
		return yamlFileConfig{}, err
	}
	var yamlConfig yamlFileConfig
	if err := yaml.Unmarshal(yamlBytes, &yamlConfig); err != nil {
		return yamlFileConfig{}, err
	}
	return yamlConfig, nil
}
