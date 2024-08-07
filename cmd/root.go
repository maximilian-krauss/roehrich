package cmd

import (
	"fmt"
	"github.com/maximilian-krauss/roerich/config"
	"github.com/maximilian-krauss/roerich/gitlab"
	"github.com/maximilian-krauss/roerich/input"
	"github.com/maximilian-krauss/roerich/utils"
	"github.com/spf13/cobra"
	"log"
)

func onlyUrls(_ *cobra.Command, args []string) error {
	maybeUrl := args[0]
	return input.ValidateUrl(maybeUrl)
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
			return fmt.Errorf("merge request is %s. aborting", mergeRequest.State)
		} else {
			log.Printf("merge request is in valid state: %s\n", mergeRequest.State)
		}

		isPipelineRunning := mergeRequest.Pipeline.Status == "running" || mergeRequest.Pipeline.Status == "pending"
		if !isPipelineRunning {
			log.Printf("pipeline is not running. current state: %s", mergeRequest.Pipeline.Status)
			return nil
		}

		jobs, err := gitlab.GetJobs(mergeRequest, cfg.Gitlab)
		if err != nil {
			return err
		}
		jobsGroupedByStage := utils.GroupByProperty(jobs, func(j gitlab.Job) string {
			return j.Stage
		})

		for stage, group := range jobsGroupedByStage {
			log.Println(fmt.Sprintf("=== %s ===", stage))
			for _, job := range group {
				log.Printf("%s: %s\n", job.Name, job.Status)
			}
		}

		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
