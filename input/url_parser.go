package input

import (
	"fmt"
	"net/url"
)

const (
	MergeRequest = iota + 1
	Job
	Pipeline
)

func ValidateUrl(maybeUrl string) error {
	parsedUrl, err := url.ParseRequestURI(maybeUrl)
	if err != nil {
		return fmt.Errorf("is not a valid url: %s", err.Error())
	}
	if parsedUrl.Scheme == "" || parsedUrl.Host == "" {
		return fmt.Errorf("%s is not a valid url: Relative urls are not supported", maybeUrl)
	}
	return nil
}
