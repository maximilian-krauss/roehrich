package gitlab

import (
	"encoding/json"
	"fmt"
	"github.com/maximilian-krauss/roerich/config"
	"io"
	"net/http"
	"net/url"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func parseAndReturnServerError(path string, responseBody []byte) error {
	var errorResponse ErrorResponse
	err := json.Unmarshal(responseBody, &errorResponse)
	if err != nil {
		return err
	}

	return fmt.Errorf("cannot get %s: %s", path, errorResponse.Message)
}

func Get[T any](path string, config config.GitlabConfig, responseType T) (T, error) {
	httpClient := http.Client{}
	joinedUrl, err := url.JoinPath(config.BaseUrl, path)
	if err != nil {
		return responseType, err
	}

	request, err := http.NewRequest("GET", joinedUrl, nil)
	if err != nil {
		return responseType, err
	}
	request.Header = http.Header{
		"Content-Type":  []string{"application/json"},
		"Accept":        []string{"application/json"},
		"Authorization": []string{"Bearer " + config.Token},
	}

	response, err := httpClient.Do(request)
	if err != nil {
		return responseType, err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return responseType, err
	}
	if response.StatusCode != http.StatusOK {
		return responseType, parseAndReturnServerError(joinedUrl, body)
	}

	var tBody T
	if err := json.Unmarshal(body, &tBody); err != nil {
		return responseType, err
	}
	return tBody, nil
}
