package searchers

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awsworkflow"
	"github.com/rkoval/alfred-aws-console-services-workflow/caching"
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

type RDSDatabaseSearcher struct{}

func (s RDSDatabaseSearcher) Search(wf *aw.Workflow, query string, session *session.Session, forceFetch bool, fullQuery string) error {
	cacheName := util.GetCurrentFilename()
	es := caching.LoadRdsDBInstanceArrayFromCache(wf, session, cacheName, s.fetch, forceFetch, fullQuery)
	for _, e := range es {
		s.addToWorkflow(wf, query, session.Config, e)
	}
	return nil
}

func (s RDSDatabaseSearcher) fetch(session *session.Session) ([]rds.DBInstance, error) {
	svc := rds.New(session)

	resp, err := svc.DescribeDBInstances(&rds.DescribeDBInstancesInput{})
	if err != nil {
		return nil, err
	}

	databases := []rds.DBInstance{}
	for i := range resp.DBInstances {
		databases = append(databases, *resp.DBInstances[i])
	}
	return databases, nil
}

func (s RDSDatabaseSearcher) addToWorkflow(wf *aw.Workflow, query string, config *aws.Config, entity rds.DBInstance) {
	subtitleArray := []string{}
	var engineString string
	if entity.Engine != nil && *entity.Engine != "" {
		engineString += *entity.Engine
	}
	if entity.EngineVersion != nil && *entity.EngineVersion != "" {
		engineString += " " + *entity.EngineVersion
	}
	subtitleArray = util.AppendString(subtitleArray, &engineString)
	subtitleArray = util.AppendString(subtitleArray, entity.DBInstanceClass)

	title := *entity.DBInstanceIdentifier
	if entity.DBName != nil && *entity.DBName != title {
		subtitleArray = util.AppendString(subtitleArray, entity.DBName)
	}

	subtitle := strings.Join(subtitleArray, " – ")

	util.NewURLItem(wf, title).
		Subtitle(subtitle).
		Arg(fmt.Sprintf(
			"https://%s.console.aws.amazon.com/rds/home?region=%s#database:id=%s;is-cluster=false",
			*config.Region,
			*config.Region,
			*entity.DBInstanceIdentifier,
		)).
		Icon(awsworkflow.GetImageIcon("rds")).
		Valid(true)
}
