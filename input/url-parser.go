package input

import (
	"fmt"

	"github.com/pkg/errors"
)

const (
	MergeRequest = iota + 1
	Job
	Pipeline
)

type ParsedUrl struct {
	url   string
	uType int
}

func parseUrl(url string) (*ParsedUrl, error) {
	fmt.Println(url)
	return nil, errors.New("Nope")
}
