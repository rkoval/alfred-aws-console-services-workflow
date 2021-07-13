package searchers

//go:generate go run ../generators/searchers_by_service_id_sorter/main.go

var cloudFormationStackSearcher = &CloudFormationStackSearcher{}
var cloudWatchLogGroupSearcher = &CloudWatchLogGroupSearcher{}
var ec2InstanceSearcher = &EC2InstanceSearcher{}
var ec2SecurityGroupSearcher = &EC2SecurityGroupSearcher{}
var elasticBeanstalkEnvironmentSearcher = &ElasticBeanstalkEnvironmentSearcher{}
var lambdaFunctionSearcher = &LambdaFunctionSearcher{}
var rdsDatabaseSearcher = &RDSDatabaseSearcher{}
var s3BucketSearcher = &S3BucketSearcher{}
var snsTopicSearcher = &SNSTopicSearcher{}
var wafIPSetSearcher = &WAFIPSetSearcher{}
var wafWebACLSearcher = &WAFWebACLSearcher{}

var SearchersByServiceId map[string]Searcher = map[string]Searcher{
	"cloudformation":                cloudFormationStackSearcher,
	"cloudformation_stacks":         cloudFormationStackSearcher,
	"cloudwatch":                    cloudWatchLogGroupSearcher,
	"cloudwatch_loggroups":          cloudWatchLogGroupSearcher,
	"ec2":                           ec2InstanceSearcher,
	"ec2_instances":                 ec2InstanceSearcher,
	"ec2_securitygroups":            ec2SecurityGroupSearcher,
	"elasticbeanstalk":              elasticBeanstalkEnvironmentSearcher,
	"elasticbeanstalk_environments": elasticBeanstalkEnvironmentSearcher,
	"lambda":                        lambdaFunctionSearcher,
	"lambda_functions":              lambdaFunctionSearcher,
	"rds":                           rdsDatabaseSearcher,
	"rds_databases":                 rdsDatabaseSearcher,
	"s3":                            s3BucketSearcher,
	"s3_buckets":                    s3BucketSearcher,
	"sns":                           snsTopicSearcher,
	"sns_topics":                    snsTopicSearcher,
	"waf":                           wafWebACLSearcher,
	"waf_ipsets":                    wafIPSetSearcher,
	"waf_webacls":                   wafWebACLSearcher,
}
