package workflow

import (
	"net/http"

	aw "github.com/deanishe/awgo"
)

type populater = func(*aw.Workflow, string, http.RoundTripper, bool, string) error

var PopulatersByServiceId map[string]populater = map[string]populater{
	"ec2":                PopulateEC2Instances,
	"ec2_instances":      PopulateEC2Instances,
	"ec2_securitygroups": PopulateEC2SecurityGroups,
	"s3":                 PopulateS3Buckets,
	"s3_home":            PopulateS3Buckets,
	// "elasticbeanstalk":              SearchElasticBeanstalkEnvironments,
	// "elasticbeanstalk_environments": SearchElasticBeanstalkEnvironments,
}
