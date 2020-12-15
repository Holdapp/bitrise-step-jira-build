package main

import (
	"fmt"
	"os"

	"github.com/fredyshox/bitrise-step-jira-build/bitrise"
)

func main() {
	if len(os.Args) < 5 {
		fmt.Printf("Usage: %s <bitrise-token> <app-slug> <workflow> <branch>\n", os.Args[0])
		os.Exit(1)
	}

	token := os.Args[1]
	appSlug := os.Args[2]
	workflow := os.Args[3]
	branch := os.Args[4]

	client := bitrise.Client{
		Token: token,
	}
	builds, err := client.ListBuilds(appSlug, workflow, branch)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	fmt.Printf("Found %d builds\n", len(builds))
	for i, build := range builds {
		fmt.Printf("======== BUILD %d ========\n", i)
		fmt.Println(build)
	}
}
