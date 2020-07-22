package searchers

import (
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/tests"
)

func TestSearchElasticBeanstalkEnvironments(t *testing.T) {
	wf := aw.New()

	r := tests.NewAWSRecorder("elastic_beanstalk_environments_fixture")
	defer r.Stop()
	SearchElasticBeanstalkEnvironments(wf, "elasticbeanstalk", r)

	cupaloy.SnapshotT(t, wf.Feedback.Items)
}
