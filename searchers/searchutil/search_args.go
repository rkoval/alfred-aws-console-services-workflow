package searchutil

import (
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

type SearchArgs struct {
	Query                  string
	Cfg                    aws.Config
	ForceFetch             bool
	FullQuery              string
	Profile                string
	IgnoreAutocompleteTerm bool
}

func (s *SearchArgs) GetAutocomplete(replaced string) string {
	if s.IgnoreAutocompleteTerm {
		return s.FullQuery
	}
	if s.Query == "" {
		return s.FullQuery + replaced + " "
	}
	autocomplete := util.ReplaceRight(s.FullQuery, s.Query, replaced, 1)
	if !strings.HasSuffix(autocomplete, " ") {
		autocomplete += " "
	}
	return autocomplete
}
