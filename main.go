package main

import (
	"os"

	"github.com/Holdapp/bitrise-step-jira-build/bitrise"
	"github.com/Holdapp/bitrise-step-jira-build/config"
	"github.com/Holdapp/bitrise-step-jira-build/service"
	logger "github.com/bitrise-io/go-utils/log"

	"github.com/bitrise-io/go-steputils/stepconf"
)

type StepConfig struct {
	// Generar info
	AppVersion string `env:"app_version,required"`

	// JIRA
	JiraHost         string          `env:"jira_host,required"`
	JiraUsername     string          `env:"jira_username,required"`
	JiraToken        stepconf.Secret `env:"jira_access_token,required"`
	JiraFieldID      int             `env:"jira_custom_field_id,required"`
	JiraIssuePattern string          `env:"jira_issue_pattern,required"`

	// Bitrise API
	BitriseToken stepconf.Secret `env:"bitrise_api_token,required"`

	// Fields provided by Bitrise
	BuildNumber string `env:"BITRISE_BUILD_NUMBER,required"`
	Workflow    string `env:"BITRISE_TRIGGERED_WORKFLOW_TITLE,required"`
	SourceDir   string `env:"BITRISE_SOURCE_DIR,required"`
	Branch      string `env:"BITRISE_GIT_BRANCH,required"`
	BuildSlug   string `env:"BITRISE_BUILD_SLUG,required"`
	AppSlug     string `env:"BITRISE_APP_SLUG,required"`
}

func (config *StepConfig) JiraTokenString() string {
	return string(config.JiraToken)
}

func (config *StepConfig) BitriseTokenString() string {
	return string(config.BitriseToken)
}

func main() {
	// Parse config
	var stepConfig = StepConfig{}
	if err := stepconf.Parse(&stepConfig); err != nil {
		logger.Errorf("Configuration error: %s", err)
		os.Exit(1)
	}

	build := config.Build{
		Version: stepConfig.AppVersion,
		Number:  stepConfig.BuildNumber,
	}

	// get commit hashes from bitrise
	logger.Infof("Scanning Bitrise API for previous failed/aborted builds\n")
	bitriseClient := bitrise.Client{Token: stepConfig.BitriseTokenString()}
	hashes, err := service.ScanRelatedCommits(
		&bitriseClient, stepConfig.AppSlug,
		stepConfig.BuildSlug, stepConfig.Workflow,
		stepConfig.Branch,
	)
	if err != nil {
		logger.Errorf("Bitrise error: %s\n", err)
		os.Exit(2)
	}

	// scan repo for related issue keys
	logger.Infof("Scanning git repo for JIRA issues (%d anchor[s])\n", len(hashes))
	gitWorker, err := service.GitOpen(
		stepConfig.SourceDir, stepConfig.Branch,
		stepConfig.JiraIssuePattern, hashes,
	)
	if err != nil {
		logger.Errorf("Git error: %s\n", err)
		os.Exit(3)
	}

	issueKeys := gitWorker.ScanIssues()

	// update custom field on issues with current build number
	logger.Infof("Updating build status for issues: %v\n", issueKeys)
	jiraWorker, err := service.NewJIRAWorker(
		stepConfig.JiraHost, stepConfig.JiraUsername,
		stepConfig.JiraTokenString(), stepConfig.JiraFieldID,
	)
	if err != nil {
		logger.Errorf("JIRA error: %s\n", err)
		os.Exit(4)
	}

	jiraWorker.UpdateBuildForIssues(issueKeys, build)

	// exit with success code
	os.Exit(0)
}
