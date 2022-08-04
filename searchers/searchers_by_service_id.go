package searchers

//go:generate go run ../generators/searchers_by_service_id_sorter/main.go

var cloudFormationStackSearcher = &CloudFormationStackSearcher{}
var cloudWatchLogGroupSearcher = &CloudWatchLogGroupSearcher{}
var cloudwatchLogInsightsQuerySearcher = &CloudWatchLogInsightsQuerySearcher{}
var codePipelinePipelineSearcher = &CodePipelinePipelinesSearcher{}
var ec2InstanceSearcher = &EC2InstanceSearcher{}
var ec2LoadBalancerSearcher = &EC2LoadBalancerSearcher{}
var ec2SecurityGroupSearcher = &EC2SecurityGroupSearcher{}
var elasticBeanstalkEnvironmentSearcher = &ElasticBeanstalkEnvironmentSearcher{}
var elasticacheMemcachedClusterSearcher = &ElasticacheMemcachedClusterSearcher{}
var elasticacheRedisClusterSearcher = &ElasticacheRedisClusterSearcher{}
var elasticbeanstalkApplicationSearcher = &ElasticBeanstalkApplicationSearcher{}
var lambdaFunctionSearcher = &LambdaFunctionSearcher{}
var rdsDatabaseSearcher = &RDSDatabaseSearcher{}
var route53HostedZoneSearcher = &Route53HostedZoneSearcher{}
var s3BucketSearcher = &S3BucketSearcher{}
var snsSubscriptionSearcher = &SNSSubscriptionSearcher{}
var snsTopicSearcher = &SNSTopicSearcher{}
var wafIPSetSearcher = &WAFIPSetSearcher{}
var wafWebACLSearcher = &WAFWebACLSearcher{}

var SearchersByServiceId map[string]Searcher = map[string]Searcher{
	"cloudformation":                cloudFormationStackSearcher,
	"cloudformation_stacks":         cloudFormationStackSearcher,
	"cloudwatch":                    cloudWatchLogGroupSearcher,
	"cloudwatch_loggroups":          cloudWatchLogGroupSearcher,
	"cloudwatch_loginsights":        cloudwatchLogInsightsQuerySearcher,
	"codepipeline_pipelines":        codePipelinePipelineSearcher,
	"ec2":                           ec2InstanceSearcher,
	"ec2_instances":                 ec2InstanceSearcher,
	"ec2_loadbalancers":             ec2LoadBalancerSearcher,
	"ec2_securitygroups":            ec2SecurityGroupSearcher,
	"elasticache":                   elasticacheRedisClusterSearcher,
	"elasticache_memcached":         elasticacheMemcachedClusterSearcher,
	"elasticache_redis":             elasticacheRedisClusterSearcher,
	"elasticbeanstalk":              elasticBeanstalkEnvironmentSearcher,
	"elasticbeanstalk_applications": elasticbeanstalkApplicationSearcher,
	"elasticbeanstalk_environments": elasticBeanstalkEnvironmentSearcher,
	"lambda":                        lambdaFunctionSearcher,
	"lambda_functions":              lambdaFunctionSearcher,
	"rds":                           rdsDatabaseSearcher,
	"rds_databases":                 rdsDatabaseSearcher,
	"route53":                       route53HostedZoneSearcher,
	"route53_hostedzones":           route53HostedZoneSearcher,
	"s3":                            s3BucketSearcher,
	"s3_buckets":                    s3BucketSearcher,
	"sns":                           snsTopicSearcher,
	"sns_subscriptions":             snsSubscriptionSearcher,
	"sns_topics":                    snsTopicSearcher,
	"waf":                           wafWebACLSearcher,
	"waf_ipsets":                    wafIPSetSearcher,
	"waf_webacls":                   wafWebACLSearcher,
}
