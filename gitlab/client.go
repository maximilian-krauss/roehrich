package gitlab

import (
	"encoding/json"
	"errors"
	"github.com/maximilian-krauss/roerich/config"
)

type PersonalAccessTokenResponse struct {
	Active  bool     `json:"active"`
	Revoked bool     `json:"revoked"`
	Scopes  []string `json:"scopes"`
}

func CheckToken(config config.GitlabConfig) error {
	responseBody, err := Get("personal_access_tokens/self", config)
	if err != nil {
		return err
	}

	var personalAccessTokenResponse PersonalAccessTokenResponse
	if err := json.Unmarshal(responseBody, &personalAccessTokenResponse); err != nil {
		return err
	}

	if !personalAccessTokenResponse.Active || personalAccessTokenResponse.Revoked {
		return errors.New("access token is either revoked or not active")
	}

	//TODO: Check if response has scope access to: read_api and read_user

	return nil
}
