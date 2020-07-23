package workflow

import (
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/tests"
)

func TestSearchElasticBeanstalkEnvironments(t *testing.T) {
	wf := aw.New()

	r := tests.NewAWSRecorder("elastic_beanstalk_environments_test")
	defer tests.PanicOnError(r.Stop)
	err := SearchElasticBeanstalkEnvironments(wf, "elasticbeanstalk", r)
	if err != nil {
		t.Errorf("got error from search: %v", err)
	}

	cupaloy.SnapshotT(t, wf.Feedback.Items)
}
