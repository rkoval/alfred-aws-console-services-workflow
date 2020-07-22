package searchers

import (
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/tests"
)

func TestSearchEC2SecurityGroups(t *testing.T) {
	wf := aw.New()

	r := tests.NewAWSRecorder("ec2_security_groups_fixture")
	defer r.Stop()
	SearchEC2SecurityGroups(wf, "", r)

	cupaloy.SnapshotT(t, wf.Feedback.Items)
}
