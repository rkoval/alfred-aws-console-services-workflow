package filters

import aw "github.com/deanishe/awgo"

type Searcher = func(wf *aw.Workflow, query string) error

var SearchersByServiceId map[string]Searcher = map[string]Searcher{
	"ec2": SearchEC2Instances,
}
