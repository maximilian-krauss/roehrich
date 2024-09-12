package cmd

import (
	"log"
	"time"

	"github.com/maximilian-krauss/roehrich/config"
	"github.com/maximilian-krauss/roehrich/gitlab"
	"github.com/maximilian-krauss/roehrich/input"
	"github.com/maximilian-krauss/roehrich/utils"
	"github.com/spf13/cobra"
)

func onlyUrls(_ *cobra.Command, args []string) error {
	maybeUrl := args[0]
	return input.ValidateUrl(maybeUrl)
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
	log.Printf("%s  %s\n", utils.JobStatusToEmoji(job.Status), job.Name)
}

var rootCmd = &cobra.Command{
	Use:               "roehrich",
	Short:             "Tut das not?",
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	Args: cobra.MatchAll(
		cobra.ExactArgs(1),
		onlyUrls,
		cobra.OnlyValidArgs,
	),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadConfig()
		if err != nil {
			return err
		}
		mrInfo, err := input.GetMRInfo(args[0])
		if err != nil {
			return err
		} else {
			log.Printf("Found project name %s and merge request id %s", mrInfo.ProjectName, mrInfo.Id)
		}

		if err := gitlab.CheckToken(cfg.Gitlab); err != nil {
			return err
		} else {
			log.Println("access token verified")
		}

		mergeRequest, err := gitlab.GetMergeRequest(mrInfo, cfg.Gitlab)
		if err != nil {
			return err
		} else {
			log.Printf("resolved merge request: %s\n", mergeRequest.Title)
		}

		if mergeRequest.State != "opened" {
			log.Printf("merge request is already %s", mergeRequest.State)
			return nil
		} else {
			log.Printf("merge request is in valid state: %s\n", mergeRequest.State)
		}

		if mergeRequest.Pipeline.Status == "success" {
			log.Printf("%s pipeline did already succeed", "✅")
			return nil
		}

		jobs, err := gitlab.GetJobs(mergeRequest, cfg.Gitlab, nil)
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
			jobs, err := gitlab.GetJobs(mergeRequest, cfg.Gitlab, gitlab.FinishedJobStatuses)
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

			pipeline, err := gitlab.GetPipeline(mergeRequest, cfg.Gitlab)
			if err != nil {
				return err
			}
			if !pipeline.IsPendingOrRunning {
				log.Printf("pipeline changed status to %s", pipeline.Status)
				break
			}

			time.Sleep(10 * time.Second)
		}

		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
