package service

import (
	"fmt"

	"github.com/Holdapp/bitrise-step-jira-build/config"
	logger "github.com/bitrise-io/go-utils/log"

	"github.com/andygrunwald/go-jira"
)

type JIRAWorker struct {
	Auth          jira.BasicAuthTransport
	Client        *jira.Client
	CustomFieldID int
}

func NewJIRAWorker(baseURL string, username string, password string, customFieldID int) (*JIRAWorker, error) {
	auth := jira.BasicAuthTransport{
		Username: username,
		Password: password,
	}

	client, err := jira.NewClient(auth.Client(), baseURL)
	if err != nil {
		return nil, err
	}

	worker := JIRAWorker{
		Auth:          auth,
		Client:        client,
		CustomFieldID: customFieldID,
	}

	return &worker, nil
}

func (worker *JIRAWorker) UpdateBuildForIssues(issueKeys []string, build config.Build) {
	for _, key := range issueKeys {
		buildString := build.String()
		customFieldKey := fmt.Sprintf("customfield_%v", worker.CustomFieldID)

		fields := map[string]string{
			customFieldKey: buildString,
		}
		body := map[string]interface{}{
			"fields": fields,
		}

		_, err := worker.Client.Issue.UpdateIssue(key, body)
		if err != nil {
			logger.Warnf("Error for '%s': %v\n", key, err)
			// TODO Response body
		}
	}
}
