package workflow

import (
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
		query: "elasticache redis",
	},
	{
		query: "elasticache redis ",
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
		query: "cloudfront",
	},
	{
		query: "cloudfront ",
	},
	{
		query: "cloudfront fle",
	},
	{
		query: "OPEN_ALL",
	},
	{
		query: "ec OPEN_ALL",
	},
	{
		query: "ec OPEN_ALL ",
	},
	{
		query:                       "OPEN_ALL ec2",
		deleteItemArgBeforeSnapshot: true,
	},
	{
		query: "eec2",
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
}

func testWorkflow(t *testing.T, tc testCase, forceFetch, snapshot bool) []*aw.Item {
	updater := &tests.MockAlfredUpdater{}
	wf := aw.New(aw.Update(updater))

	session, r := tests.NewAWSRecorderSession(tc.fixtureName)
	defer tests.PanicOnError(r.Stop)
	Run(wf, tc.query, session, forceFetch, false, "../console-services.yml")

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
