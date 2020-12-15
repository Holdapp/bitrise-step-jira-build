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

func Open(path string, branch string, commits []string) *GitWorker {
	repo, err := git.OpenRepository(path)
	if err != nil {
		return nil
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

	return worker
}

func Load(url string, branch string, remote string) *GitWorker {
	return nil
}

func (self *GitWorker) LoadCommits() []*git.Commit {
	var commits = make([]*git.Commit, 0)
	for _, oid := range self.MergeCommits {
		rangeString := fmt.Sprintf("%s^..%s", oid.String(), oid.String())
		fmt.Println(rangeString)
		revwalk, err := self.Repo.Walk()
		spec, err := self.Repo.Revparse(rangeString)
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

func (self *GitWorker) ScanIssues() []string {
	commits := self.LoadCommits()
	issueKeysMap := make(map[string]bool)
	regex, err := regexp.Compile(self.IssuePattern)
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
