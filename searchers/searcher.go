package searchers

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	aw "github.com/deanishe/awgo"
)

type Searcher interface {
	Search(wf *aw.Workflow, query string, cfg aws.Config, forceFetch bool, fullQuery string) error
}
