package main

import (
	"fmt"
	"os"

	"github.com/Holdapp/bitrise-step-jira-build/service"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("Usage %v <path> <branch> <hash>\n", os.Args[0])
		os.Exit(1)
	}

	repoPath := os.Args[1]
	branchName := os.Args[2]
	commitHashes := []string{os.Args[3]}
	worker := service.Open(repoPath, branchName, "origin", commitHashes)

	fmt.Println("Repo opened")
	commits := worker.LoadCommits()
	for _, commit := range commits {
		fmt.Println(commit.Id().String())
		fmt.Println(commit.Message())
	}
}
