package retry

import (
	"log"

	"github.com/maximilian-krauss/roehrich/config"
	"github.com/maximilian-krauss/roehrich/gitlab"
	"github.com/maximilian-krauss/roehrich/input"
)

type Args struct {
	SourceUrl  string
	ConfigPath string
}

func Run(args Args) error {
	cfg, err := config.LoadConfig(args.ConfigPath)
	if err != nil {
		return err
	}
	mrInfo, err := input.GetMRInfo(args.SourceUrl)
	if err != nil {
		return err
	}
	log.Printf("Found project name %s and merge request id %s", mrInfo.ProjectName, mrInfo.Id)

	gitlabConfig, err := config.GetConfigByHostname(mrInfo.HostName, *cfg)
	if err != nil {
		return err
	}

	if err := gitlab.CheckToken(*gitlabConfig); err != nil {
		return err
	}
	log.Println("access token verified")

	mergeRequest, err := gitlab.GetMergeRequest(mrInfo, *gitlabConfig)
	if err != nil {
		return err
	}
	log.Printf("resolved merge request: %s\n", mergeRequest.Title)

	if mergeRequest.State != "opened" {
		log.Printf("merge request is already %s", mergeRequest.State)
		return nil
	} else {
		log.Printf("merge request is in valid state: %s\n", mergeRequest.State)
	}

	if mergeRequest.Pipeline.Status == "success" {
		log.Println("pipeline did already succeed")
		return nil
	}

	jobs, err := gitlab.GetJobs(mergeRequest, *gitlabConfig, nil)
	if err != nil {
		return err
	}

	for _, job := range gitlab.FilterFailedJobs(jobs) {
		log.Printf("retrying failed job %s", job.Name)
		if _, err := gitlab.RetryJob(mergeRequest, *gitlabConfig, job); err != nil {
			return err
		}
	}

	return nil
}
