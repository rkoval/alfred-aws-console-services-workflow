package searchers

import (
	"github.com/rkoval/alfred-aws-console-services-workflow/searchtypes"
)

var SearchersBySearchType map[searchtypes.SearchType]Searcher = map[searchtypes.SearchType]Searcher{
	searchtypes.EC2Instances:                 &EC2InstanceSearcher{},
	searchtypes.EC2SecurityGroups:            &EC2SecurityGroupSearcher{},
	searchtypes.S3Buckets:                    &S3BucketSearcher{},
	searchtypes.ElasticBeanstalkEnvironments: &ElasticBeanstalkEnvironmentSearcher{},
	searchtypes.WAFIPSets:                    &WAFIPSetSearcher{},
	searchtypes.WAFWebACLs:                   &WAFWebACLSearcher{},
}
