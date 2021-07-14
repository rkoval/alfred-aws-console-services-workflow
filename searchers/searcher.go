package searchers

import (
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
)

type Searcher interface {
	Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error
}
