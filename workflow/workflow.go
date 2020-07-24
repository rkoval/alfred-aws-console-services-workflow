package workflow

import (
	"log"

	"github.com/aws/aws-sdk-go/aws/session"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/parsers"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchtypes"
)

func Run(wf *aw.Workflow, query string, session *session.Session, forceFetch bool, ymlPath string) {
	awsServices := parsers.ParseConsoleServicesYml(ymlPath)
	fullQuery := query
	query, searchType, awsService := parsers.ParseQuery(awsServices, query)

	var err error
	if searchType == searchtypes.Services {
		log.Println("using searcher associated with services")
		SearchServices(wf, awsServices)
	} else if searchType == searchtypes.SubServices {
		log.Println("using searcher associated with sub-services")
		SearchSubServices(wf, *awsService)
	} else {
		log.Printf("using searcher associated with %d", searchType)
		searcher := searchers.SearchersBySearchType[searchType]
		err = searcher(wf, query, session, forceFetch, fullQuery)
	}

	if err != nil {
		wf.FatalError(err)
	}

	if query != "" {
		log.Printf("filtering with query %s", query)
		res := wf.Filter(query)

		log.Printf("%d results match %q", len(res), query)

		for i, r := range res {
			log.Printf("%02d. score=%0.1f sortkey=%s", i+1, r.Score, wf.Feedback.Keywords(i))
		}
	}

	wf.WarnEmpty("No matching services found", "Try a different query?")

	wf.SendFeedback()
}
