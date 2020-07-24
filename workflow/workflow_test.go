package workflow

import (
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	aw "github.com/deanishe/awgo"
	"github.com/go-test/deep"
	"github.com/rkoval/alfred-aws-console-services-workflow/tests"
)

type testCase struct {
	query       string
	fixtureName string
}

var tcs []testCase = []testCase{
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
		query:       "s3 home ",
		fixtureName: "s3_buckets_test", // reuse test fixture from this other test
	},
	{
		query:       "s3 home bucket-1",
		fixtureName: "s3_buckets_test", // reuse test fixture from this other test
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
		query: "ec2 secur",
	},
	{
		query: "ec2 securitygroups",
	},
	{
		query:       "ec2 securitygroups ",
		fixtureName: "ec2_security_groups_test", // reuse test fixture from this other test
	},
	{
		query:       "ec2 securitygroups sg-000000",
		fixtureName: "ec2_security_groups_test", // reuse test fixture from this other test
	},
	{
		query: "ec2 inst",
	},
	{
		query: "ec2 instances",
	},
	{
		query:       "ec2 instances ",
		fixtureName: "ec2_instances_test", // reuse test fixture from this other test
	},
	{
		query:       "ec2 instances environment-name-1",
		fixtureName: "ec2_instances_test", // reuse test fixture from this other test
	},
	{
		query:       "ec2 instances i-aaaaaaaaaa",
		fixtureName: "ec2_instances_test", // reuse test fixture from this other test
	},
	{
		query:       "ec2 $",
		fixtureName: "ec2_instances_test", // reuse test fixture from this other test
	},
	{
		query:       "ec2 $ ",
		fixtureName: "ec2_instances_test", // reuse test fixture from this other test
	},
	{
		query:       "ec2 $environment-name-1",
		fixtureName: "ec2_instances_test", // reuse test fixture from this other test
	},
	{
		query:       "ec2 $i-aaaaaaaaaa",
		fixtureName: "ec2_instances_test", // reuse test fixture from this other test
	},
}

func testWorkflow(t *testing.T, tc testCase, forceFetch, snapshot bool) []*aw.Item {
	wf := aw.New()
	r := tests.NewAWSRecorder(tc.fixtureName)
	defer tests.PanicOnError(r.Stop)
	Run(wf, tc.query, r, forceFetch, "../console-services.yml")

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
			if diff := deep.Equal(fetchedItems, cachedItems); diff != nil {
				t.Error(diff)
			}
		})
	}
}
