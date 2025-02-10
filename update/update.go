package update

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/hashicorp/go-version"
)

const ReleaseUrl string = "https://api.github.com/repos/maximilian-krauss/roehrich/releases/latest"

type GithubRelease struct {
	TagName string `json:"tag_name"`
	Url     string `json:"html_url"`
}

func getLatestGithubRelease() (*GithubRelease, error) {
	httpClient := http.Client{}
	requestUri, err := url.ParseRequestURI(ReleaseUrl)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("GET", requestUri.String(), nil)
	if err != nil {
		return nil, err
	}
	request.Header = http.Header{
		"Content-Type": []string{"application/json"},
		"Accept":       []string{"application/json"},
	}

	response, err := httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request to %s failed with status code %d", ReleaseUrl, response.StatusCode)
	}

	var githubRelease GithubRelease
	if err := json.Unmarshal(body, &githubRelease); err != nil {
		return nil, err
	}
	return &githubRelease, nil
}

type VersionInfo struct {
	Version string
	IsNewer bool
	Url     string
}

func FindLatestVersion(currentVersionString string) (*VersionInfo, error) {
	latestRelease, err := getLatestGithubRelease()
	if err != nil {
		return nil, err
	}

	currentVersion, err := version.NewVersion(currentVersionString)
	if err != nil {
		return nil, err
	}
	releaseVersion, err := version.NewVersion(latestRelease.TagName)
	if err != nil {
		return nil, err
	}

	return &VersionInfo{
		Version: latestRelease.TagName,
		IsNewer: releaseVersion.GreaterThan(currentVersion),
		Url:     latestRelease.Url,
	}, nil
}
