package service

import (
	"fmt"

	"github.com/Holdapp/bitrise-step-jira-build/bitrise"
)

func ScanRelatedCommits(client *bitrise.Client, appSlug string, buildSlug string, workflow string, branch string) ([]string, error) {
	builds, err := client.ListBuilds(appSlug, workflow, branch)
	if err != nil {
		return nil, err
	}

	currentBuildIndex := findCurrentBuild(builds, buildSlug)
	if currentBuildIndex < 0 {
		return nil, fmt.Errorf("bitrise: Current build not found")
	}

	commitHashes := []string{builds[currentBuildIndex].CommitHash}
	for _, build := range builds[currentBuildIndex+1:] {
		switch build.Status {
		case bitrise.BuildAbortedWithSuccess, bitrise.BuildAbortedWithFailure, bitrise.BuildFailed:
			commitHashes = append(commitHashes, build.CommitHash)
		default:
			break
		}
	}

	return commitHashes, nil
}

func findCurrentBuild(builds []bitrise.Build, buildSlug string) int {
	var currentBuildIndex int = -1
	for i, build := range builds {
		if build.Slug == buildSlug {
			currentBuildIndex = i
			break
		}
	}

	return currentBuildIndex
}
