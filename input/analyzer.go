package input

import (
	"fmt"
	"net/url"
	"regexp"
)

type MergeRequestInfo struct {
	Id          string
	ProjectName string
}

func GetMRInfo(suppliedUrl string) (*MergeRequestInfo, error) {
	parsedUrl, err := url.ParseRequestURI(suppliedUrl)
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`\/*(?P<name>[a-z\/\-_]+)\/-\/merge_requests\/(?P<id>[0-9]+)`)
	matches := re.FindStringSubmatch(parsedUrl.Path)

	if matches == nil || len(matches) != 3 {
		return nil, fmt.Errorf("%s is not a valid merge request url", suppliedUrl)
	}

	mergeRequestInfo := &MergeRequestInfo{
		Id:          matches[re.SubexpIndex("id")],
		ProjectName: matches[re.SubexpIndex("name")],
	}

	return mergeRequestInfo, nil
}
