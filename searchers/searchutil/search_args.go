package searchutil

import (
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

type SearchArgs struct {
	Query                  string
	Cfg                    aws.Config
	ForceFetch             bool
	FullQuery              string
	Profile                string
	GetRegionFunc          func(cfg aws.Config) string
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

func (s *SearchArgs) GetRegion() string {
	return s.GetRegionFunc(s.Cfg)
}

func (s *SearchArgs) AddMatch(item *aw.Item, idPrefix, id, title string) *aw.Item {
	if idPrefix != "" && id != "" && strings.HasPrefix(s.Query, idPrefix) {
		item.Match(id).
			Autocomplete(s.GetAutocomplete(id))
	} else {
		item.Autocomplete(s.GetAutocomplete(title))
	}
	return item
}
