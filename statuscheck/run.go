package statuscheck

import (
	"github.com/fatih/color"
	"github.com/maximilian-krauss/roehrich/static"
	"github.com/maximilian-krauss/roehrich/update"
	"log"
	"time"

	"github.com/maximilian-krauss/roehrich/config"
	"github.com/maximilian-krauss/roehrich/gitlab"
	"github.com/maximilian-krauss/roehrich/input"
	"github.com/maximilian-krauss/roehrich/utils"
)

type Args struct {
	SourceUrl                string
	PollingIntervalInSeconds int
	ConfigPath               string
	SkipVersionCheck         bool
}

func printGroupedJobs(jobs []gitlab.Job) {
	jobsGroupedByStage := utils.GroupByProperty(jobs, func(j gitlab.Job) string {
		return j.Stage
	})

	for stage, group := range jobsGroupedByStage {
		log.Printf("=== %s ===", stage)
		for _, job := range group {
			printJob(job)
		}
	}
}

func printJob(job gitlab.Job) {
	statusColor := utils.JobStatusToColor(job.Status)
	log.Printf("%s  %s (%s)\n", statusColor.SprintFunc()("["+job.Status+"]"), job.Name, job.Stage)
}

func versionCheck(args Args) {
	if args.SkipVersionCheck {
		log.Printf("Version check is skipped")
		return
	}
	remoteVersion, err := update.FindLatestVersion(static.ApplicationVersion)
	if err != nil || remoteVersion == nil {
		log.Printf("Version check failed: %s", err)
		return
	}
	if remoteVersion.IsNewer {
		yellow := color.New(color.FgYellow)
		log.Printf("current version: %s [%s]", static.ApplicationVersion, yellow.Sprintf("%s available", remoteVersion.Version))
	} else {
		green := color.New(color.FgGreen)
		log.Printf("current version: %s [%s]", static.ApplicationVersion, green.Sprint("up to date"))
	}
}

func Run(args Args) error {
	cfg, err := config.LoadConfig(args.ConfigPath)
	if err != nil {
		return err
	}
	versionCheck(args)

	mrInfo, err := input.GetMRInfo(args.SourceUrl)
	if err != nil {
		return err
	}
	log.Printf("Found project name %s and merge request id %s", mrInfo.ProjectName, mrInfo.Id)

	gitlabConfig, err := config.GetConfigByHostname(mrInfo.HostName, *cfg)
	if err != nil {
		return err
	}

	if err := gitlab.CheckToken(*gitlabConfig, false); err != nil {
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
	printGroupedJobs(jobs)

	finishedJobs := make(map[int]gitlab.Job)
	for _, job := range gitlab.FilterFinishedJobs(jobs) {
		finishedJobs[job.Id] = job
	}

	log.Println("waiting for other jobs to complete...")

	for {
		jobs, err := gitlab.GetJobs(mergeRequest, *gitlabConfig, gitlab.FinishedJobStatuses)
		if err != nil {
			return err
		}
		for _, job := range jobs {
			if _, exists := finishedJobs[job.Id]; exists {
				continue
			}
			finishedJobs[job.Id] = job
			printJob(job)
		}

		pipeline, err := gitlab.GetPipeline(mergeRequest, *gitlabConfig)
		if err != nil {
			return err
		}
		if !pipeline.IsPendingOrRunning {
			statusColor := utils.JobStatusToColor(pipeline.Status)
			log.Printf("pipeline changed status to %s", statusColor.SprintFunc()(pipeline.Status))
			break
		}

		time.Sleep(time.Second * time.Duration(args.PollingIntervalInSeconds))
	}

	return nil
}
