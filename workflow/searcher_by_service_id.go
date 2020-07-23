package workflow

import (
	"net/http"

	aw "github.com/deanishe/awgo"
)

type searcher = func(*aw.Workflow, string, http.RoundTripper, bool, string) error

var SearchersByServiceId map[string]searcher = map[string]searcher{
	"ec2":           PopulateEC2Instances,
	"ec2_instances": PopulateEC2Instances,
	// "ec2_securitygroups":            SearchEC2SecurityGroups,
	// "s3":                            SearchS3Buckets,
	// "s3_buckets":                    SearchS3Buckets,
	// "elasticbeanstalk":              SearchElasticBeanstalkEnvironments,
	// "elasticbeanstalk_environments": SearchElasticBeanstalkEnvironments,
}
