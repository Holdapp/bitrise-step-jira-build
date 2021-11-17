package service

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/Holdapp/bitrise-step-jira-build/config"
	logger "github.com/bitrise-io/go-utils/log"

	"github.com/andygrunwald/go-jira"
)

const MultiFieldSeparator string = ", "

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

		log.Printf("New custom field content: \"%v\" \n", buildString)

		err := worker.sendUpdateCustomFieldRequest(key, buildString)
		if err != nil {
			logger.Warnf("Error while updating '%s': %v\n", key, err)
		}
	}
}

func (worker *JIRAWorker) UpdateBuildForIssuesMultiField(issueKeys []string, build config.Build) {
	for _, key := range issueKeys {
		fields, _, err := worker.Client.Issue.GetCustomFields(key)
		if err != nil {
			logger.Warnf("Error while getting custom fields of %s: %v\n", key, err)
			continue
		}

		// parse existing builds
		customFieldKey := fmt.Sprintf("customfield_%v", worker.CustomFieldID)
		currentFieldContent := fields[customFieldKey]
		currentBuildStrings := strings.Split(currentFieldContent, MultiFieldSeparator)
		currentBuilds := make([]config.Build, 0)
		for _, s := range currentBuildStrings {
			build, err := config.ParseBuild(s)
			if err == nil {
				currentBuilds = append(currentBuilds, *build)
			}
		}

		// create new build string
		newBuildStrings := make([]string, 0)
		newBuildStrings = append(newBuildStrings, build.String())
		for _, currentBuild := range currentBuilds {
			if currentBuild.Version != build.Version {
				newBuildStrings = append(newBuildStrings, currentBuild.String())
			}
		}
		sort.Strings(newBuildStrings)
		newFieldContent := strings.Join(newBuildStrings, MultiFieldSeparator)

		log.Printf("Current custom field content: \"%v\" \n", currentFieldContent)
		log.Printf("Current build list: %v , len: %d\n", currentBuildStrings, len(currentBuildStrings))
		log.Printf("New build list: %v , len: %d\n", newBuildStrings, len(newBuildStrings))
		log.Printf("New custom field content: \"%v\" \n", newFieldContent)

		err = worker.sendUpdateCustomFieldRequest(key, newFieldContent)
		if err != nil {
			logger.Warnf("Error while updating '%s': %v\n", key, err)
		}
	}
}

func (worker *JIRAWorker) sendUpdateCustomFieldRequest(issueKey string, buildString string) error {
	customFieldKey := fmt.Sprintf("customfield_%v", worker.CustomFieldID)
	fields := map[string]string{
		customFieldKey: buildString,
	}
	body := map[string]interface{}{
		"fields": fields,
	}

	_, err := worker.Client.Issue.UpdateIssue(issueKey, body)

	return err
}
