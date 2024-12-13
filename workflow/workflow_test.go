package workflow

import (
	"os"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/tests"
	"github.com/stretchr/testify/assert"
)

type testCase struct {
	query                       string
	fixtureName                 string
	deleteItemArgBeforeSnapshot bool
}

var tcs []testCase = []testCase{
	{
		query: "",
	},
	{
		query: " ",
	},
	{
		query: "$",
	},
	{
		query: "$us-",
	},
	{
		// autocomplete for this test is not working properly, but just keep track of it
		query: "$ us-",
	},
	{
		query: "$us-east-1",
	},
	{
		query: "$us-east-1 ",
	},
	{
		query: "$asdf asdf asdf",
	},
	{
		query: "@",
	},
	{
		query: "@prof",
	},
	{
		// autocomplete for this test is not working properly, but just keep track of it
		query: "@ prof",
	},
	{
		query: "@profile1",
	},
	{
		query: "@asdf asdf asdf",
	},
	{
		query: "@$",
	},
	{
		query: "@ $",
	},
	{
		query: "$ @",
	},
	{
		query: "$ @ adsf asdf",
	},
	{
		query: "$us-east 1",
	},
	{
		query: "$us-east-1 @",
	},
	{
		query: "$us-east-1 @prof",
	},
	{
		query: "$us-east-1 @prof asdf asdf",
	},
	{
		query: "$us-east-1 @profile3 elasticbeanstalk",
	},
	{
		query: "$cn-north-1 elasticbeanstalk",
	},
	{
		query: "@usgov elasticbeanstalk",
	},
	{
		query: "@china elasticbeanstalk",
	},
	{
		query: "alex",
	},
	{
		query: " alexa",
	},
	{
		query: "alexa",
	},
	{
		query: "alexa ",
	},
	{
		query: "$us-east-1 alexa ",
	},
	{
		query: "alexa $us-east-1",
	},
	{
		query: "alexa $us-east-1 ",
	},
	{
		query: "alexa h",
	},
	{
		query: "alexa home",
	},
	{
		query: "alexa home ",
	},
	{
		query: "elasticache",
	},
	{
		query: "elasticache ",
	},
	{
		query:       "elasticache redis",
		fixtureName: "../searchers/elasticache_redis_clusters_test", // reuse test fixture from this other test
	},
	{
		query:       "elasticache redis ",
		fixtureName: "../searchers/elasticache_redis_clusters_test", // reuse test fixture from this other test
	},
	{
		query:       "cloudwatch loggroups",
		fixtureName: "../searchers/cloudwatch_log_groups_test", // reuse test fixture from this other test
	},
	{
		query:       "cloudwatch loggroups ",
		fixtureName: "../searchers/cloudwatch_log_groups_test", // reuse test fixture from this other test
	},
	{
		query:       "cloudwatch loggroups log-group-aaa",
		fixtureName: "../searchers/cloudwatch_log_groups_test", // reuse test fixture from this other test
	},
	{
		query:       "cloudwatch ,",
		fixtureName: "../searchers/cloudwatch_log_groups_test", // reuse test fixture from this other test
	},
	{
		query:       "cloudwatch ,log-group-bbb",
		fixtureName: "../searchers/cloudwatch_log_groups_test", // reuse test fixture from this other test
	},
	{
		query:       "s3 home ",
		fixtureName: "../searchers/s3_buckets_test", // reuse test fixture from this other test
	},
	{
		query:       "s3 home bucket-1",
		fixtureName: "../searchers/s3_buckets_test", // reuse test fixture from this other test
	},
	{
		query:       "s3 buckets ",
		fixtureName: "../searchers/s3_buckets_test", // reuse test fixture from this other test
	},
	{
		query:       "s3 buckets bucket-1",
		fixtureName: "../searchers/s3_buckets_test", // reuse test fixture from this other test
	},
	{
		query: "$us-east-1 elasticbeanstalk",
	},
	{
		query: "$us-east-1 $us-east-2 elasticbeanstalk",
	},
	{
		query: "$us-east-1 elasticbeanstalk appli",
	},
	{
		query:       "$us-east-1 elasticbeanstalk applications",
		fixtureName: "../searchers/elastic_beanstalk_applications_test_us-east-1", // reuse test fixture from this other test
	},
	{
		query:       "$us-east-1 elasticbeanstalk applications ",
		fixtureName: "../searchers/elastic_beanstalk_applications_test_us-east-1", // reuse test fixture from this other test
	},
	{
		query:       "elasticbeanstalk $us-east-1 applications",
		fixtureName: "../searchers/elastic_beanstalk_applications_test_us-east-1", // reuse test fixture from this other test
	},
	{
		query:       "elasticbeanstalk $us-east-1 applications ",
		fixtureName: "../searchers/elastic_beanstalk_applications_test_us-east-1", // reuse test fixture from this other test
	},
	{
		query:       "elasticbeanstalk applications $us-east-1",
		fixtureName: "../searchers/elastic_beanstalk_applications_test_us-east-1", // reuse test fixture from this other test
	},
	{
		query:       "elasticbeanstalk applications $us-east-1 ",
		fixtureName: "../searchers/elastic_beanstalk_applications_test_us-east-1", // reuse test fixture from this other test
	},
	{
		query:       "elasticbeanstalk applications $us-east-1 Ap",
		fixtureName: "../searchers/elastic_beanstalk_applications_test_us-east-1", // reuse test fixture from this other test
	},
	{
		query:       "elasticbeanstalk applications $us-east-1 Ap ",
		fixtureName: "../searchers/elastic_beanstalk_applications_test_us-east-1", // reuse test fixture from this other test
	},
	{
		query:       "elasticbeanstalk applications App1 $us-east-1",
		fixtureName: "../searchers/elastic_beanstalk_applications_test_us-east-1", // reuse test fixture from this other test
	},
	{
		query:       "elasticbeanstalk applications App1 $us-",
		fixtureName: "../searchers/elastic_beanstalk_applications_test_us-east-1", // reuse test fixture from this other test
	},
	{
		query:       "elasticbeanstalk applications App1 $us-east-1 ",
		fixtureName: "../searchers/elastic_beanstalk_applications_test_us-east-1", // reuse test fixture from this other test
	},
	{
		query:       "elasticbeanstalk applications $us-east-1 arn:aws:elasticbeanstalk:us-east-1:0000000000:application/Ap",
		fixtureName: "../searchers/elastic_beanstalk_applications_test_us-east-1", // reuse test fixture from this other test
	},
	{
		query:       "elasticbeanstalk applications arn:aws:elasticbeanstalk:us-east-1:0000000000:application/Ap $us-east-1 ",
		fixtureName: "../searchers/elastic_beanstalk_applications_test_us-east-1", // reuse test fixture from this other test
	},
	{
		query: "@profile1 elasticbeanstalk",
	},
	{
		query: "@profile1 $us-east-2 elasticbeanstalk",
	},
	{
		query: "@profile1 elasticbeanstalk appli",
	},
	{
		query:       "@profile1 elasticbeanstalk applications",
		fixtureName: "../searchers/elastic_beanstalk_applications_test_us-east-1", // reuse test fixture from this other test
	},
	{
		query:       "@noregion elasticbeanstalk applications",
		fixtureName: "../searchers/elastic_beanstalk_applications_test_us-east-1", // reuse test fixture from this other test
	},
	{
		query:       "@profile1 elasticbeanstalk applications ",
		fixtureName: "../searchers/elastic_beanstalk_applications_test_us-east-1", // reuse test fixture from this other test
	},
	{
		query:       "elasticbeanstalk @profile1 applications",
		fixtureName: "../searchers/elastic_beanstalk_applications_test_us-east-1", // reuse test fixture from this other test
	},
	{
		query:       "elasticbeanstalk @profile1 applications ",
		fixtureName: "../searchers/elastic_beanstalk_applications_test_us-east-1", // reuse test fixture from this other test
	},
	{
		query:       "elasticbeanstalk applications @profile1",
		fixtureName: "../searchers/elastic_beanstalk_applications_test_us-east-1", // reuse test fixture from this other test
	},
	{
		query:       "elasticbeanstalk applications @profile1 ",
		fixtureName: "../searchers/elastic_beanstalk_applications_test_us-east-1", // reuse test fixture from this other test
	},
	{
		query:       "elasticbeanstalk applications @profile1 Ap",
		fixtureName: "../searchers/elastic_beanstalk_applications_test_us-east-1", // reuse test fixture from this other test
	},
	{
		query:       "elasticbeanstalk applications @profile1 Ap ",
		fixtureName: "../searchers/elastic_beanstalk_applications_test_us-east-1", // reuse test fixture from this other test
	},
	{
		query:       "elasticbeanstalk applications App1 @profile1",
		fixtureName: "../searchers/elastic_beanstalk_applications_test_us-east-1", // reuse test fixture from this other test
	},
	{
		query:       "elasticbeanstalk applications App1 @prof",
		fixtureName: "../searchers/elastic_beanstalk_applications_test_us-east-1", // reuse test fixture from this other test
	},
	{
		query:       "elasticbeanstalk applications App1 @profile1 ",
		fixtureName: "../searchers/elastic_beanstalk_applications_test_us-east-1", // reuse test fixture from this other test
	},
	{
		query:       "elasticbeanstalk applications @profile1 arn:aws:elasticbeanstalk:us-east-1:0000000000:application/Ap",
		fixtureName: "../searchers/elastic_beanstalk_applications_test_us-east-1", // reuse test fixture from this other test
	},
	{
		query:       "elasticbeanstalk applications arn:aws:elasticbeanstalk:us-east-1:0000000000:application/Ap @profile1 ",
		fixtureName: "../searchers/elastic_beanstalk_applications_test_us-east-1", // reuse test fixture from this other test
	},

	{
		query:       "elasticbeanstalk applications @profile3 $us-east-1 ",
		fixtureName: "../searchers/elastic_beanstalk_applications_test_us-east-1", // reuse test fixture from this other test
	},
	{
		query: "lambda",
	},
	{
		query: "lambda ",
	},
	{
		query: "lambda func",
	},
	{
		query: "cloudformation",
	},
	{
		query: "cloudformation ",
	},
	{
		query: "ecr ",
	},
	{
		query: "ecr repo",
	},
	{
		query: "cloudformation stacks",
	},
	{
		query:       "cloudformation ,",
		fixtureName: "../searchers/cloudformation_stacks_test", // reuse test fixture from this other test
	},
	{
		query:       "cloudformation stacks awseb-e-aaaaaaaaaa-",
		fixtureName: "../searchers/cloudformation_stacks_test", // reuse test fixture from this other test
	},
	{
		query:       "cloudformation ,Custom",
		fixtureName: "../searchers/cloudformation_stacks_test", // reuse test fixture from this other test
	},
	{
		query:       "lambda ,",
		fixtureName: "../searchers/lambda_functions_test", // reuse test fixture from this other test
	},
	{
		query:       "lambda ,Function2",
		fixtureName: "../searchers/lambda_functions_test", // reuse test fixture from this other test
	},
	{
		query:       "rds databases",
		fixtureName: "../searchers/rds_databases_test", // reuse test fixture from this other test
	},
	{
		query:       "rds databases ",
		fixtureName: "../searchers/rds_databases_test", // reuse test fixture from this other test
	},
	{
		query:       "rds ,instance",
		fixtureName: "../searchers/rds_databases_test", // reuse test fixture from this other test
	},
	{
		query: "cloudfront",
	},
	{
		query: "cloudfront ",
	},
	{
		query: "cloudfront fle",
	},
	{
		query:                       "OPEN_ALL",
		deleteItemArgBeforeSnapshot: true,
	},
	{
		query:                       "ec OPEN_ALL",
		deleteItemArgBeforeSnapshot: true,
	},
	{
		query:                       "ec OPEN_ALL ",
		deleteItemArgBeforeSnapshot: true,
	},
	{
		query:                       "OPEN_ALL ec2",
		deleteItemArgBeforeSnapshot: true,
	},
	{
		query: "eec2",
	},
	{
		query: "eec2 $us-east-1",
	},
	{
		query: "eec2 @profile1",
	},
	{
		query: "eec2 $us-east-1 ",
	},
	{
		query: "eec2 @profile1 ",
	},
	{
		query: "$us-east-1 eec2",
	},
	{
		query: "@profile1 eec2",
	},
	{
		query: "ec2",
	},
	{
		query: "ec2 ",
	},
	{
		query:                       "ec2 OPEN_ALL",
		deleteItemArgBeforeSnapshot: true,
	},
	{
		query: "ec2 secur",
	},
	{
		query: "ec2 securitygroups",
	},
	{
		query:       "ec2 securitygroups ",
		fixtureName: "../searchers/ec2_security_groups_test", // reuse test fixture from this other test
	},
	{
		query:       "ec2 securitygroups sg-000000",
		fixtureName: "../searchers/ec2_security_groups_test", // reuse test fixture from this other test
	},
	{
		query: "ec2 tags ",
	},
	{
		query: "ec2 tags asdf",
	},
	{
		query: "bean",
	},
	{
		query: "elasticbeanstalk",
	},
	{
		query: "app",
	},
	{
		query: "apprunner",
	},
	{
		query:       "elasticbeanstalk ",
		fixtureName: "../searchers/elastic_beanstalk_environments_test", // reuse test fixture from this other test
	},
	{
		query:       "elasticbeanstalk ,",
		fixtureName: "../searchers/elastic_beanstalk_environments_test", // reuse test fixture from this other test
	},
	{
		query:       "elasticbeanstalk environments e-cccccc",
		fixtureName: "../searchers/elastic_beanstalk_environments_test", // reuse test fixture from this other test
	},
	{
		query: "ec2 inst",
	},
	{
		query: "ec2 instances",
	},
	{
		query:       "ec2 instances ",
		fixtureName: "../searchers/ec2_instances_test", // reuse test fixture from this other test
	},
	{
		query:       "ec2 instances environment-name-1",
		fixtureName: "../searchers/ec2_instances_test", // reuse test fixture from this other test
	},
	{
		query:       "ec2 instances i-aaaaaaaaaa",
		fixtureName: "../searchers/ec2_instances_test", // reuse test fixture from this other test
	},
	{
		query:       "ec2 ,",
		fixtureName: "../searchers/ec2_instances_test", // reuse test fixture from this other test
	},
	{
		query:       "ec2 , ",
		fixtureName: "../searchers/ec2_instances_test", // reuse test fixture from this other test
	},
	{
		query:       "ec2 ,environment-name-1",
		fixtureName: "../searchers/ec2_instances_test", // reuse test fixture from this other test
	},
	{
		query:       "ec2 ,i-aaaaaaaaaa",
		fixtureName: "../searchers/ec2_instances_test", // reuse test fixture from this other test
	},
	{
		query:       "waf ",
		fixtureName: "../searchers/waf_ip_sets_test", // reuse test fixture from this other test
	},
	{
		query:       "waf ipsets ",
		fixtureName: "../searchers/waf_ip_sets_test", // reuse test fixture from this other test
	},
	{
		query:       "waf ipsets ipset-3",
		fixtureName: "../searchers/waf_ip_sets_test", // reuse test fixture from this other test
	},
	{
		query:       "waf webacls ",
		fixtureName: "../searchers/waf_web_acls_test", // reuse test fixture from this other test
	},
	{
		query:       "waf webacls webacl-2",
		fixtureName: "../searchers/waf_web_acls_test", // reuse test fixture from this other test
	},
	{
		query:       "waf ,webacl-2",
		fixtureName: "../searchers/waf_web_acls_test", // reuse test fixture from this other test
	},
	{
		query:       "sns",
		fixtureName: "../searchers/sns_topics_test", // reuse test fixture from this other test
	},
	{
		query:       "sns ",
		fixtureName: "../searchers/sns_topics_test", // reuse test fixture from this other test
	},
	{
		query:       "sns topics",
		fixtureName: "../searchers/sns_topics_test", // reuse test fixture from this other test
	},
	{
		query:       "sns topics ",
		fixtureName: "../searchers/sns_topics_test", // reuse test fixture from this other test
	},
	{
		query:       "sns topics sns-topic-name-1",
		fixtureName: "../searchers/sns_topics_test", // reuse test fixture from this other test
	},
	{
		query:       "sns subscriptions",
		fixtureName: "../searchers/sns_subscriptions_test", // reuse test fixture from this other test
	},
	{
		query:       "sns subscriptions ",
		fixtureName: "../searchers/sns_subscriptions_test", // reuse test fixture from this other test
	},
	{
		query:       "sns subscriptions paginated",
		fixtureName: "../searchers/sns_subscriptions_test", // reuse test fixture from this other test
	},
	{
		query:       "elasticache memcached",
		fixtureName: "../searchers/elasticache_memcached_clusters_test", // reuse test fixture from this other test
	},
	{
		query:       "elasticache memcached ",
		fixtureName: "../searchers/elasticache_memcached_clusters_test", // reuse test fixture from this other test
	},
	{
		query: "oss",
	},
	{
		query: "oss integrat",
	},
	{
		query:       "ec2 loadbalancers",
		fixtureName: "../searchers/ec2_load_balancers_test", // reuse test fixture from this other test
	},
	{
		query:       "ec2 loadbalancers ",
		fixtureName: "../searchers/ec2_load_balancers_test", // reuse test fixture from this other test
	},
	{
		query:       "ec2 loadbalancers arn:aws:elasticloadbalancing:us-west-2:0000000000:loadbalancer/net/awseb-AWSEB-BBBBBBBBBBBBB/bbbbbbbbbbbbbbbb",
		fixtureName: "../searchers/ec2_load_balancers_test", // reuse test fixture from this other test
	},
	{
		query:       "elasticbeanstalk applications",
		fixtureName: "../searchers/elastic_beanstalk_applications_test", // reuse test fixture from this other test
	},
	{
		query:       "elasticbeanstalk applications ",
		fixtureName: "../searchers/elastic_beanstalk_applications_test", // reuse test fixture from this other test
	},
	{
		query:       "elasticbeanstalk applications Ap p1",
		fixtureName: "../searchers/elastic_beanstalk_applications_test", // reuse test fixture from this other test
	},
	{
		query:       "elasticbeanstalk applications arn:aws:elasticbeanstalk:us-west-2:0000000000:application/App3",
		fixtureName: "../searchers/elastic_beanstalk_applications_test", // reuse test fixture from this other test
	},
	{
		query:       "route53",
		fixtureName: "../searchers/route53_hosted_zones_test", // reuse test fixture from this other test
	},
	{
		query:       "route53 ",
		fixtureName: "../searchers/route53_hosted_zones_test", // reuse test fixture from this other test
	},
	{
		query:       "route53 hostedzones",
		fixtureName: "../searchers/route53_hosted_zones_test", // reuse test fixture from this other test
	},
	{
		query:       "route53 hostedzones ",
		fixtureName: "../searchers/route53_hosted_zones_test", // reuse test fixture from this other test
	},
	{
		query:       "route53 hostedzones ZWWWWWWWWWWWWW",
		fixtureName: "../searchers/route53_hosted_zones_test", // reuse test fixture from this other test
	},
	{
		query:       "cloudwatch",
		fixtureName: "../searchers/cloud_watch_log_insights_queries_test", // reuse test fixture from this other test
	},
	{
		query:       "cloudwatch loginsights",
		fixtureName: "../searchers/cloud_watch_log_insights_queries_test", // reuse test fixture from this other test
	},
	{
		query:       "cloudwatch loginsights ",
		fixtureName: "../searchers/cloud_watch_log_insights_queries_test", // reuse test fixture from this other test
	},
	{
		query:       "codepipeline",
		fixtureName: "../searchers/codepipeline_pipelines_test", // reuse test fixture from this other test
	},
	{
		query:       "codepipeline pipelines",
		fixtureName: "../searchers/codepipeline_pipelines_test", // reuse test fixture from this other test
	},
	{
		query:       "codepipeline pipelines ",
		fixtureName: "../searchers/codepipeline_pipelines_test", // reuse test fixture from this other test
	},
	{
		query:       "codepipeline pipelines pipeline-name-1",
		fixtureName: "../searchers/codepipeline_pipelines_test", // reuse test fixture from this other test
	},
}

func testWorkflow(t *testing.T, tc testCase, forceFetch, snapshot bool) []*aw.Item {
	updater := &tests.MockAlfredUpdater{}
	wf := aw.New(aw.Update(updater))

	r := tests.NewAWSRecorderSession(tc.fixtureName)
	defer tests.PanicOnError(r.Stop)
	Run(wf, tc.query, r, forceFetch, false, "../console-services.yml")

	if tc.deleteItemArgBeforeSnapshot {
		for i := range wf.Feedback.Items {
			wf.Feedback.Items[i] = wf.Feedback.Items[i].Arg("[redacted]")
		}
	}

	if snapshot {
		cupaloy.SnapshotT(t, wf.Feedback.Items)
	}
	return wf.Feedback.Items
}

func TestRun(t *testing.T) {
	for _, tc := range tcs {
		t.Run(tc.query, func(t *testing.T) {
			testWorkflow(t, tc, true, true)
		})
	}
}

func TestRunWithCache(t *testing.T) {
	for _, tc := range tcs {
		t.Run(tc.query+"_cached", func(t *testing.T) {
			fetchedItems := testWorkflow(t, tc, true, false)
			cachedItems := testWorkflow(t, tc, false, false)
			assert.Equal(t, cachedItems, fetchedItems)
		})
	}
}

func TestRunWithoutRegion(t *testing.T) {
	tcs := []testCase{
		{
			query: "",
		},
	}
	awsProfile := os.Getenv("AWS_PROFILE")
	awsRegion := os.Getenv("AWS_REGION")
	os.Setenv("AWS_PROFILE", "bogus-test-profile")
	os.Setenv("AWS_REGION", "")
	for _, tc := range tcs {
		t.Run(tc.query, func(t *testing.T) {
			testWorkflow(t, tc, true, true)
		})
	}
	os.Setenv("AWS_PROFILE", awsProfile)
	os.Setenv("AWS_REGION", awsRegion)
}
