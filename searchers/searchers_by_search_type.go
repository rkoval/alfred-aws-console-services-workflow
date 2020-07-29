package searchers

import (
	"github.com/aws/aws-sdk-go/aws/session"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchtypes"
)

type Searcher = func(*aw.Workflow, string, *session.Session, bool, string) error

var SearchersBySearchType map[searchtypes.SearchType]Searcher = map[searchtypes.SearchType]Searcher{
	searchtypes.EC2Instances:                 SearchEC2Instances,
	searchtypes.EC2SecurityGroups:            SearchEC2SecurityGroups,
	searchtypes.S3Buckets:                    SearchS3Buckets,
	searchtypes.ElasticBeanstalkEnvironments: SearchElasticBeanstalkEnvironments,
	searchtypes.WAFIPSets:                    SearchWAFIPSets,
}
