package searchers

var cloudWatchLogGroupSearcher = &CloudWatchLogGroupSearcher{}
var ec2InstanceSearcher = &EC2InstanceSearcher{}
var ec2SecurityGroupSearcher = &EC2SecurityGroupSearcher{}
var elasticBeanstalkEnvironmentSearcher = &ElasticBeanstalkEnvironmentSearcher{}
var s3BucketSearcher = &S3BucketSearcher{}
var wafWebACLSearcher = &WAFWebACLSearcher{}
var wafIPSetSearcher = &WAFIPSetSearcher{}

var SearchersByServiceId map[string]Searcher = map[string]Searcher{
	"cloudwatch":                    cloudWatchLogGroupSearcher,
	"cloudwatch_loggroups":          cloudWatchLogGroupSearcher,
	"ec2":                           ec2InstanceSearcher,
	"ec2_instances":                 ec2InstanceSearcher,
	"ec2_securitygroups":            ec2SecurityGroupSearcher,
	"elasticbeanstalk":              elasticBeanstalkEnvironmentSearcher,
	"elasticbeanstalk_environments": elasticBeanstalkEnvironmentSearcher,
	"s3":                            s3BucketSearcher,
	"s3_buckets":                    s3BucketSearcher,
	"waf":                           wafWebACLSearcher,
	"waf_webacls":                   wafWebACLSearcher,
	"waf_ipsets":                    wafIPSetSearcher,
}
