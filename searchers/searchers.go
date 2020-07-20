package searchers

import aw "github.com/deanishe/awgo"

type searcher = func(wf *aw.Workflow, query string) error

var SearchersByServiceId map[string]searcher = map[string]searcher{
	"ec2":              SearchEC2Instances,
	"elasticbeanstalk": SearchElasticBeanstalkEnvironments,
}
