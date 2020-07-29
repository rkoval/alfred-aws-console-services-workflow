package searchtypes

var SearchTypesByServiceId map[string]SearchType = map[string]SearchType{
	"ec2":                           EC2Instances,
	"ec2_instances":                 EC2Instances,
	"ec2_securitygroups":            EC2SecurityGroups,
	"elasticbeanstalk":              ElasticBeanstalkEnvironments,
	"elasticbeanstalk_environments": ElasticBeanstalkEnvironments,
	"s3":                            S3Buckets,
	"s3_buckets":                    S3Buckets,
	"waf":                           WAFWebACLs,
	"waf_webacls":                   WAFWebACLs,
	"waf_ipsets":                    WAFIPSets,
}
