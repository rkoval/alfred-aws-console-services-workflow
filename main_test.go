package main

import (
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/tests"
)

type testCase struct {
	query       string
	fixtureName string
}

func testWorkflow(t *testing.T, tc testCase) {
	wf := aw.New()
	r := tests.NewAWSRecorder(tc.fixtureName)
	defer tests.PanicOnError(r.Stop)
	Run(wf, tc.query, r)

	cupaloy.SnapshotT(t, wf.Feedback.Items)
}

func TestRun(t *testing.T) {
	tcs := []testCase{
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
			fixtureName: "searchers/ec2_security_groups_test", // reuse test fixture from this other test
		},
		{
			query: "ec2 inst",
		},
		{
			query: "ec2 instances",
		},
		{
			query:       "ec2 instances ",
			fixtureName: "searchers/ec2_instances_test", // reuse test fixture from this other test
		},
		{
			query:       "ec2 instances environment-name-1",
			fixtureName: "searchers/ec2_instances_test", // reuse test fixture from this other test
		},
		{
			query:       "ec2 instances i-aaaaaaaaaa",
			fixtureName: "searchers/ec2_instances_test", // reuse test fixture from this other test
		},
		{
			query:       "ec2 $",
			fixtureName: "searchers/ec2_instances_test", // reuse test fixture from this other test
		},
		{
			query:       "ec2 $ ",
			fixtureName: "searchers/ec2_instances_test", // reuse test fixture from this other test
		},
		{
			query:       "ec2 $environment-name-1",
			fixtureName: "searchers/ec2_instances_test", // reuse test fixture from this other test
		},
		{
			query:       "ec2 $i-aaaaaaaaaa",
			fixtureName: "searchers/ec2_instances_test", // reuse test fixture from this other test
		},
	}

	for _, tc := range tcs {
		t.Run(tc.query, func(t *testing.T) {
			testWorkflow(t, tc)
		})
	}
}
