package workflow

import (
	"github.com/aws/aws-sdk-go/aws/session"
	aw "github.com/deanishe/awgo"
)

type searcher = func(*aw.Workflow, string, *session.Session, bool, string) error

var SearchersByServiceId map[string]searcher = map[string]searcher{
	"ec2":                           SearchEC2Instances,
	"ec2_instances":                 SearchEC2Instances,
	"ec2_securitygroups":            SearchEC2SecurityGroups,
	"s3":                            SearchS3Buckets,
	"s3_home":                       SearchS3Buckets,
	"elasticbeanstalk":              SearchElasticBeanstalkEnvironments,
	"elasticbeanstalk_environments": SearchElasticBeanstalkEnvironments,
}
