package searchtypes

var SearchTypesByServiceId map[string]SearchType = map[string]SearchType{
	"ec2":                           EC2Instances,
	"ec2_instances":                 EC2Instances,
	"ec2_securitygroups":            EC2SecurityGroups,
	"s3":                            S3Buckets,
	"s3_home":                       S3Buckets,
	"elasticbeanstalk":              ElasticBeanstalkEnvironments,
	"elasticbeanstalk_environments": ElasticBeanstalkEnvironments,
}
