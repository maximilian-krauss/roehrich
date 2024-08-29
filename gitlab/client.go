package gitlab

import (
	"errors"
	"net/url"
	"strconv"

	"github.com/maximilian-krauss/roehrich/config"
	"github.com/maximilian-krauss/roehrich/input"
	"github.com/maximilian-krauss/roehrich/utils"
)

type PersonalAccessTokenResponse struct {
	Active  bool     `json:"active"`
	Revoked bool     `json:"revoked"`
	Scopes  []string `json:"scopes"`
}

func CheckToken(config config.GitlabConfig) error {
	var accessToken PersonalAccessTokenResponse
	accessToken, err := Get("personal_access_tokens/self", config, accessToken)
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

func GetMergeRequest(info *input.MergeRequestInfo, config config.GitlabConfig) (MergeRequest, error) {
	var mergeRequest MergeRequest
	var mrPath = "/projects/" + url.QueryEscape(info.ProjectName) + "/merge_requests/" + info.Id
	mergeRequest, err := Get(mrPath, config, mergeRequest)

	return mergeRequest, err
}

type Job struct {
	Id        int               `json:"id"`
	Name      string            `json:"name"`
	Stage     string            `json:"stage"`
	Status    string            `json:"status"`
	CreatedAt utils.IsoDateTime `json:"created_at"`
}

func GetJobs(mr MergeRequest, config config.GitlabConfig) ([]Job, error) {
	var jobs = []Job{}
	var jobsPath = "/projects/" + strconv.Itoa(mr.ProjectId) + "/pipelines/" + strconv.Itoa(mr.Pipeline.Id) + "/jobs"
	jobs, err := GetMany(jobsPath, config, jobs)

	return jobs, err
}

func isJobRunningOrPending(job Job) bool {
	return job.Status == "created" || job.Status == "running" || job.Status == "waiting_for_resource"
}

func GetFinishedJobs(jobs []Job) []Job {
	return utils.Filter(jobs, func(job Job) bool {
		return !isJobRunningOrPending(job)
	})
}

func GetPendingJobs(jobs []Job) []Job {
	return utils.Filter(jobs, isJobRunningOrPending)
}
