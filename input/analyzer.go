package input

import "errors"

type MergeRequestInfo struct {
	Id          string
	ProjectName string
}

func GetMRInfo(url string) (*MergeRequestInfo, error) {
	return nil, errors.New("not implemented")
}
