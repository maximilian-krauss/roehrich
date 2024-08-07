package gitlab

import (
	"errors"
	"github.com/maximilian-krauss/roerich/config"
	"github.com/maximilian-krauss/roerich/input"
	"net/url"
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
	Title    string   `json:"title"`
	State    string   `json:"state"`
	Pipeline Pipeline `json:"head_pipeline"`
}

func GetMergeRequest(info *input.MergeRequestInfo, config config.GitlabConfig) (MergeRequest, error) {
	var mergeRequest MergeRequest
	var mrPath = "/projects/" + url.QueryEscape(info.ProjectName) + "/merge_requests/" + info.Id
	mergeRequest, err := Get(mrPath, config, mergeRequest)

	return mergeRequest, err
}
