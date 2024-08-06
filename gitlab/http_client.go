package gitlab

import (
	"github.com/maximilian-krauss/roerich/config"
	"io"
	"net/http"
	"net/url"
)

func Get(path string, config config.GitlabConfig) ([]byte, error) {
	httpClient := http.Client{}
	joinedUrl, err := url.JoinPath(config.BaseUrl, path)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("GET", joinedUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header = http.Header{
		"Content-Type":  []string{"application/json"},
		"Accept":        []string{"application/json"},
		"Authorization": []string{"Bearer " + config.Token},
	}

	response, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
