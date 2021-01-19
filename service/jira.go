package service

import (
	"fmt"
	"strings"

	logger "github.com/bitrise-io/go-utils/log"

	"github.com/andygrunwald/go-jira"
)

type JIRAWorker struct {
	Auth              jira.BasicAuthTransport
	Client            *jira.Client
	CustomFieldID     int
	DestinationStatus string
}

func NewJIRAWorker(baseURL string, username string, password string, customFieldID int, destStatus string) (*JIRAWorker, error) {
	auth := jira.BasicAuthTransport{
		Username: username,
		Password: password,
	}

	client, err := jira.NewClient(auth.Client(), baseURL)
	if err != nil {
		return nil, err
	}

	worker := JIRAWorker{
		Auth:              auth,
		Client:            client,
		CustomFieldID:     customFieldID,
		DestinationStatus: destStatus,
	}

	return &worker, nil
}

func (worker *JIRAWorker) UpdateFieldValueForIssues(issueKeys []string, fieldValue string) {
	for _, key := range issueKeys {
		customFieldKey := fmt.Sprintf("customfield_%v", worker.CustomFieldID)

		fields := map[string]string{
			customFieldKey: fieldValue,
		}
		body := map[string]interface{}{
			"fields": fields,
		}

		_, err := worker.Client.Issue.UpdateIssue(key, body)
		if err != nil {
			logger.Warnf("Error for '%s': %v\n", key, err)
			// TODO Print response body
		}
	}
}

func (worker *JIRAWorker) TransitionEnabled() bool {
	return len(worker.DestinationStatus) > 0
}

func (worker *JIRAWorker) TransitionIssues(issueKeys []string) {
	if !worker.TransitionEnabled() {
		logger.Infof("Transitions disabled. Proceeding...\n")
		return
	}

	logger.Infof("Transitions enabled\n")
	for _, key := range issueKeys {
		transitions, _, err := worker.Client.Issue.GetTransitions(key)
		if err != nil {
			logger.Warnf("Error while retriving available transitions for '%s': %v\n", key, err)
			continue
		}

		destinationIndex := findTransitionWithName(transitions, worker.DestinationStatus)
		if destinationIndex < 0 && destinationIndex >= len(transitions) {
			logger.Warnf("Could not find trantistion with matching name '%s' for '%s'\n", worker.DestinationStatus, key)
			continue
		}

		destination := transitions[destinationIndex]
		_, err = worker.Client.Issue.DoTransition(key, destination.ID)
		if err != nil {
			logger.Warnf("Failed to perform transition for '%s': %v\n", key, err)
			// TODO Print response body
		}
	}
}

func findTransitionWithName(transitions []jira.Transition, substr string) int {
	index := -1
	substr = strings.ToLower(substr)
	for i, transition := range transitions {
		name := strings.ToLower(transition.Name)
		if strings.Contains(name, substr) {
			index = i
			break
		}
	}

	return index
}
