package bitrise

/**
Example structure:
  {
     "abort_reason": "string",
     "branch": "string",
     "build_number": 0,
     "commit_hash": "string",
     "commit_message": "string",
     "commit_view_url": "string",
     "environment_prepare_finished_at": "string",
     "finished_at": "string",
     "is_on_hold": true,
     "machine_type_id": "string",
     "original_build_params": "string",
     "pull_request_id": 0,
     "pull_request_target_branch": "string",
     "pull_request_view_url": "string",
     "slug": "string",
     "stack_identifier": "string",
     "started_on_worker_at": "string",
     "status": 0,
     "status_text": "string",
     "tag": "string",
     "triggered_at": "string",
     "triggered_by": "string",
     "triggered_workflow": "string"
   }
*/

import "fmt"

const (
	BuildNotFinished        = 0
	BuildSuccessful         = 1
	BuildFailed             = 2
	BuildAbortedWithFailure = 3
	BuildAbortedWithSuccess = 4
)

type Build struct {
	Number     int    `json:"build_number"`
	Slug       string `json:"slug"`
	CommitHash string `json:"commit_hash"`
	Branch     string `json:"branch"`
	Status     int    `json:"status"`
	OnHold     bool   `json:"is_on_hold"`
}

func (build Build) String() string {
	return fmt.Sprintf("bitrise.Build( %s, number: %d, status: %d )", build.Slug, build.Number, build.Status)
}

type BuildResponse struct {
	Data   []Build `json:"data"`
	Paging Paging  `json:"paging"`
}
