package searchers

import aw "github.com/deanishe/awgo"

type searcher = func(wf *aw.Workflow, query string) error

var SearchersByServiceId map[string]searcher = map[string]searcher{
	"ec2":                           SearchEC2Instances,
	"ec2_instances":                 SearchEC2Instances,
	"ec2_securityGroups":            SearchEC2SecurityGroups,
	"s3":                            SearchS3Buckets,
	"s3_buckets":                    SearchS3Buckets,
	"elasticbeanstalk":              SearchElasticBeanstalkEnvironments,
	"elasticbeanstalk_environments": SearchElasticBeanstalkEnvironments,
}
