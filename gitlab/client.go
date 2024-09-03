package gitlab

import (
	"errors"
	"github.com/maximilian-krauss/roehrich/config"
	"github.com/maximilian-krauss/roehrich/input"
	"github.com/maximilian-krauss/roehrich/utils"
	"net/url"
	"slices"
	"strconv"
)

type PersonalAccessTokenResponse struct {
	Active  bool     `json:"active"`
	Revoked bool     `json:"revoked"`
	Scopes  []string `json:"scopes"`
}

func CheckToken(config config.GitlabConfig) error {
	var accessToken PersonalAccessTokenResponse
	accessToken, err := Get("personal_access_tokens/self", config, accessToken, nil)
	if err != nil {
		return err
	}

	if !accessToken.Active || accessToken.Revoked {
		return errors.New("access token is either revoked or not active")
	}

	//TODO: Check if response has scope access to: read_api and read_user

	return nil
}

type Pipeline struct {
	Id     int    `json:"id"`
	Iid    int    `json:"iid"`
	Status string `json:"status"`
}

type MergeRequest struct {
	Title     string   `json:"title"`
	State     string   `json:"state"`
	Pipeline  Pipeline `json:"head_pipeline"`
	ProjectId int      `json:"project_id"`
}

var PendingOrRunningJobStatuses = []string{"created", "pending", "running", "waiting_for_resource"}
var FinishedJobStatuses = []string{"failed", "canceled", "skipped", "success", "manual"}

func GetMergeRequest(info *input.MergeRequestInfo, config config.GitlabConfig) (MergeRequest, error) {
	var mergeRequest MergeRequest
	var mrPath = "/projects/" + url.QueryEscape(info.ProjectName) + "/merge_requests/" + info.Id
	mergeRequest, err := Get(mrPath, config, mergeRequest, nil)

	return mergeRequest, err
}

type Job struct {
	Id        int               `json:"id"`
	Name      string            `json:"name"`
	Stage     string            `json:"stage"`
	Status    string            `json:"status"`
	CreatedAt utils.IsoDateTime `json:"created_at"`
}

func GetJobs(mr MergeRequest, config config.GitlabConfig, jobStatuses []string) ([]Job, error) {
	var jobs []Job
	var jobsPath = "/projects/" + strconv.Itoa(mr.ProjectId) + "/pipelines/" + strconv.Itoa(mr.Pipeline.Id) + "/jobs"
	params := make(map[string]string)
	jobs, err := GetMany(jobsPath, config, jobs, params)
	if err != nil {
		return nil, err
	}
	if jobStatuses != nil {
		return utils.Filter(jobs, func(job Job) bool {
			return slices.Contains(jobStatuses, job.Status)
		}), nil
	}

	return jobs, err
}

func isJobRunningOrPending(job Job) bool {
	return slices.Contains(PendingOrRunningJobStatuses, job.Status)
}

func FilterFinishedJobs(jobs []Job) []Job {
	return utils.Filter(jobs, func(job Job) bool {
		return !isJobRunningOrPending(job)
	})
}

func FilterPendingJobs(jobs []Job) []Job {
	return utils.Filter(jobs, isJobRunningOrPending)
}
