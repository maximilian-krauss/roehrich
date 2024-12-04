package gitlab

import (
	"errors"
	"fmt"
	"net/url"
	"slices"
	"strconv"
	"strings"

	"github.com/maximilian-krauss/roehrich/config"
	"github.com/maximilian-krauss/roehrich/input"
	"github.com/maximilian-krauss/roehrich/utils"
)

const SCOPE_READ_API string = "read_api"
const SCOPE_READ_USER string = "read_user"
const SCOPE_API string = "api"

var READ_SCOPES = []string{SCOPE_READ_API, SCOPE_READ_USER, SCOPE_API}
var WRITE_SCOPES = []string{SCOPE_API}

type PersonalAccessTokenResponse struct {
	Active  bool     `json:"active"`
	Revoked bool     `json:"revoked"`
	Scopes  []string `json:"scopes"`
}

func CheckToken(config config.GitlabConfig, needsWriteAccess bool) error {
	var accessToken PersonalAccessTokenResponse
	accessToken, err := Get("personal_access_tokens/self", config, accessToken, nil)
	if err != nil {
		return err
	}

	if !accessToken.Active || accessToken.Revoked {
		return errors.New("access token is either revoked or not active")
	}
	if !utils.ContainsAll(accessToken.Scopes, READ_SCOPES) {
		return fmt.Errorf("no read access, token needs at least %s scopes", strings.Join(READ_SCOPES, ","))
	}
	if needsWriteAccess && !utils.ContainsAll(accessToken.Scopes, WRITE_SCOPES) {
		return fmt.Errorf("no write access, token needs at least %s scopes", strings.Join(WRITE_SCOPES, ","))
	}

	return nil
}

type Pipeline struct {
	Id                 int    `json:"id"`
	Iid                int    `json:"iid"`
	Status             string `json:"status"`
	IsPendingOrRunning bool   `json:"-"`
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

func GetPipeline(mr MergeRequest, config config.GitlabConfig) (Pipeline, error) {
	var pipeline Pipeline
	var pipelinePath, err = url.JoinPath("projects", strconv.Itoa(mr.ProjectId), "pipelines", strconv.Itoa(mr.Pipeline.Id))
	if err != nil {
		return pipeline, err
	}
	pipeline, err = Get(pipelinePath, config, mr.Pipeline, nil)
	if err != nil {
		return pipeline, err
	}
	pipeline.IsPendingOrRunning = slices.Contains(PendingOrRunningJobStatuses, pipeline.Status)
	return pipeline, err
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
	var jobsPath = fmt.Sprintf("/projects/%d/pipelines/%d/jobs", mr.ProjectId, mr.Pipeline.Id)
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

func RetryJob(mr MergeRequest, config config.GitlabConfig, failedJob Job) (Job, error) {
	var path = fmt.Sprintf("/projects/%d/jobs/%d/retry", mr.ProjectId, failedJob.Id)
	var job Job
	job, err := Post(path, config, job)

	return job, err
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

func FilterFailedJobs(jobs []Job) []Job {
	return utils.Filter(jobs, func(job Job) bool { return job.Status == "failed" })
}
