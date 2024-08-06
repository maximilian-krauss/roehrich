package input

import (
	"fmt"
	"github.com/spf13/cobra"
	"net/url"

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

func ValidateUrl(_ *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("no url provided")
	}
	var maybeUrl = args[0]
	parsedUrl, err := url.ParseRequestURI(maybeUrl)
	if err != nil {
		return fmt.Errorf("is not a valid url: %s", err.Error())
	}
	if parsedUrl.Scheme == "" || parsedUrl.Host == "" {
		return fmt.Errorf("%s is not a valid url: Relative urls are not supported", maybeUrl)
	}
	return nil
}
