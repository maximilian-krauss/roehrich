package gitlab

import (
	"encoding/json"
	"fmt"
	"github.com/maximilian-krauss/roehrich/config"
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

func makeRequest(config config.GitlabConfig, requestUrl string) (*http.Response, []byte, error) {
	httpClient := http.Client{}
	request, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		return nil, nil, err
	}
	request.Header = http.Header{
		"Content-Type":  []string{"application/json"},
		"Accept":        []string{"application/json"},
		"Authorization": []string{"Bearer " + config.Token},
	}

	response, err := httpClient.Do(request)
	if err != nil {
		return nil, nil, err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, nil, err
	}
	if response.StatusCode != http.StatusOK {
		return nil, nil, parseAndReturnServerError(requestUrl, body)
	}
	return response, body, nil
}

func Get[T any](path string, config config.GitlabConfig, responseType T) (T, error) {
	joinedUrl, err := url.JoinPath(config.BaseUrl, path)
	if err != nil {
		return responseType, err
	}

	var tBody T
	_, body, err := makeRequest(config, joinedUrl)
	if err != nil {
		return responseType, err
	}

	if err := json.Unmarshal(body, &tBody); err != nil {
		return responseType, err
	}
	return tBody, nil
}

func GetMany[T any](path string, config config.GitlabConfig, responseType []T) ([]T, error) {
	const perPage = 10
	currentPage := 1
	items := make([]T, 0)

	for {
		joinedUrl, err := url.JoinPath(config.BaseUrl, path)
		if err != nil {
			return responseType, err
		}

		var tBody []T
		response, body, err := makeRequest(
			config,
			fmt.Sprintf("%s?page=%d&per_page=%d", joinedUrl, currentPage, perPage),
		)
		if err != nil {
			return responseType, err
		}
		if err := json.Unmarshal(body, &tBody); err != nil {
			return responseType, err
		}

		items = append(items, tBody...)

		nextPage := response.Header.Get("x-next-page")
		if nextPage == "" {
			break
		}

		currentPage++
	}

	return items, nil
}
