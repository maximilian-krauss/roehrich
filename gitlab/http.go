package gitlab

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strconv"

	"github.com/maximilian-krauss/roehrich/config"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func parseAndReturnServerError(stautsCode int, path string, responseBody []byte) error {
	var errorResponse ErrorResponse
	err := json.Unmarshal(responseBody, &errorResponse)
	if err != nil {
		return err
	}

	return fmt.Errorf("request failed %s: %d %s", path, stautsCode, errorResponse.Message)
}

func makeRequest(method string, config config.GitlabConfig, requestUrl string, queryParameter map[string]string) (*http.Response, []byte, error) {
	httpClient := http.Client{}
	requestUri, err := url.ParseRequestURI(requestUrl)
	if err != nil {
		return nil, nil, err
	}
	if queryParameter != nil {
		query := requestUri.Query()
		for key, value := range queryParameter {
			query.Add(key, value)
		}
		requestUri.RawQuery = query.Encode()
	}

	request, err := http.NewRequest(method, requestUri.String(), nil)
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

	if !slices.Contains([]int{http.StatusOK, http.StatusCreated, http.StatusNoContent}, response.StatusCode) {
		return nil, nil, parseAndReturnServerError(
			response.StatusCode,
			fmt.Sprintf("[%s] %s", method, requestUrl),
			body,
		)
	}
	return response, body, nil
}

func Post[TResponse any](path string, config config.GitlabConfig, responseType TResponse) (TResponse, error) {
	joinedUrl, err := url.JoinPath(config.BaseUrl, path)
	if err != nil {
		return responseType, err
	}

	var tBody TResponse
	_, body, err := makeRequest("POST", config, joinedUrl, nil)
	if err != nil {
		return responseType, err
	}

	if err := json.Unmarshal(body, &tBody); err != nil {
		return responseType, err
	}
	return tBody, nil
}

func Get[T any](path string, config config.GitlabConfig, responseType T, queryParameter map[string]string) (T, error) {
	joinedUrl, err := url.JoinPath(config.BaseUrl, path)
	if err != nil {
		return responseType, err
	}

	var tBody T
	_, body, err := makeRequest("GET", config, joinedUrl, queryParameter)
	if err != nil {
		return responseType, err
	}

	if err := json.Unmarshal(body, &tBody); err != nil {
		return responseType, err
	}
	return tBody, nil
}

func GetMany[T any](path string, config config.GitlabConfig, responseType []T, additionalQueryParameter map[string]string) ([]T, error) {
	const perPage = 50
	currentPage := 1
	items := make([]T, 0)

	for {
		joinedUrl, err := url.JoinPath(config.BaseUrl, path)
		if err != nil {
			return responseType, err
		}

		queryParameter := map[string]string{"page": strconv.Itoa(currentPage), "per_page": strconv.Itoa(perPage)}
		for key, value := range additionalQueryParameter {
			queryParameter[key] = value
		}

		var tBody []T
		response, body, err := makeRequest(
			"GET",
			config,
			joinedUrl,
			queryParameter,
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
