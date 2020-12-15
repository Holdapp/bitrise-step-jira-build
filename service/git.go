package service

import (
	"fmt"
	"log"
	"regexp"

	git "github.com/libgit2/git2go/v30"
)

const DefaultIssuePattern string = `([A-Z]{1,10}-[0-9]+)`

type GitWorker struct {
	Repo         *git.Repository
	Branch       string
	Remote       string
	MergeCommits []*git.Oid
	IssuePattern string
}

func GitOpen(path string, branch string, commits []string) (*GitWorker, error) {
	repo, err := git.OpenRepository(path)
	if err != nil {
		return nil, err
	}

	repo.SetHead(branch)

	commitHashes := make([]*git.Oid, 0, len(commits))
	for _, commit := range commits {
		hash, err := git.NewOid(commit)
		if err == nil {
			commitHashes = append(commitHashes, hash)
		}
	}

	worker := new(GitWorker)
	worker.Repo = repo
	worker.Branch = branch
	worker.MergeCommits = commitHashes
	worker.IssuePattern = DefaultIssuePattern

	return worker, nil
}

func GitLoad(url string, branch string, remote string) (*GitWorker, error) {
	log.Fatalln("Not implemented!")
	return nil, nil
}

func (worker *GitWorker) LoadCommits() []*git.Commit {
	var commits = make([]*git.Commit, 0)
	for _, oid := range worker.MergeCommits {
		mergeCommit, err := worker.Repo.LookupCommit(oid)
		if mergeCommit.ParentCount() < 2 {
			fmt.Printf("%s is not merge commit!\n", oid.String())
			commits = append(commits, mergeCommit)
			continue
		}

		rangeString := fmt.Sprintf("%s^..%s", oid.String(), oid.String())
		fmt.Printf("Revparse range: %s\n", rangeString)
		revwalk, err := worker.Repo.Walk()
		spec, err := worker.Repo.Revparse(rangeString)
		if err != nil {
			fmt.Println("Revparse error: ", err)
			continue
		}

		fromID := spec.From().Id()
		toID := spec.To().Id()
		if err := revwalk.Hide(fromID); err != nil {
			fmt.Println("revwalk.Hide error", err)
			continue
		}

		if err := revwalk.Push(toID); err != nil {
			fmt.Println("revwalk.Push error", err)
			continue
		}

		revwalk.Iterate(func(commit *git.Commit) bool {
			commits = append(commits, commit)
			return true
		})
	}

	return commits
}

func (worker *GitWorker) ScanIssues() []string {
	commits := worker.LoadCommits()
	issueKeysMap := make(map[string]bool)
	regex, err := regexp.Compile(worker.IssuePattern)
	if err != nil {
		log.Fatal(err)
	}

	for _, commit := range commits {
		keys := regex.FindAllString(commit.Message(), -1)
		for _, key := range keys {
			issueKeysMap[key] = true
		}
	}

	issueKeys := make([]string, 0, len(issueKeysMap))
	for k := range issueKeysMap {
		issueKeys = append(issueKeys, k)
	}

	return issueKeys
}
