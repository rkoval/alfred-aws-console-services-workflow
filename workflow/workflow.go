package workflow

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awsworkflow"
	"github.com/rkoval/alfred-aws-console-services-workflow/parsers"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchtypes"
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

func Run(wf *aw.Workflow, query string, session *session.Session, forceFetch, openAll bool, ymlPath string) {
	log.Println("using workflow cacheDir: " + wf.CacheDir())
	log.Println("using workflow dataDir: " + wf.DataDir())
	awsServices := parsers.ParseConsoleServicesYml(ymlPath)
	fullQuery := query
	query, searchType, awsService, promptOpenAll := parsers.ParseQuery(awsServices, query)

	if promptOpenAll && awsService != nil {
		if openAll {
			for i := range awsService.SubServices {
				OpenServiceInBrowser(wf, &awsService.SubServices[i])
			}
		} else {
			cmd := fmt.Sprintf(`%s -query="%s" -open_all`, os.Args[0], fullQuery)
			wf.NewItem(fmt.Sprintf("Open the %d %s sub-services in browser", len(awsService.SubServices)+1, awsService.Id)).
				Subtitle("Fair warning: this may briefly overload your system").
				Icon(aw.IconWarning).
				Arg(cmd).
				Var("action", "run-script").
				Valid(true)
		}
		wf.SendFeedback()
		return
	}

	var err error
	if searchType == searchtypes.None {
		log.Println("no search type parsed")
		wf.NewItem("Search for an AWS Service ...").
			Subtitle("e.g., cloudformation, ec2, s3 ...")

		if wf.UpdateCheckDue() {
			if err := wf.CheckForUpdate(); err != nil {
				wf.FatalError(err)
			}
		}
		if wf.UpdateAvailable() {
			util.NewURLItem(wf, "Update available").
				Subtitle("Select this result to navigate to download").
				Arg("https://github.com/rkoval/alfred-aws-console-services-workflow/releases").
				Icon(aw.IconInfo)
		}
	} else if searchType == searchtypes.Services {
		log.Println("using searcher associated with services")
		if awsService == nil {
			SearchServices(wf, awsServices)
		} else {
			awsServices = []awsworkflow.AwsService{}
			awsServices = append(awsServices, *awsService)
			SearchServices(wf, awsServices)
		}
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

	if wf.IsEmpty() {
		wf.NewItem("No matching services found").
			Subtitle("Try another query (example: `aws ec2 instances`)").
			Icon(aw.IconNote)
	}

	wf.SendFeedback()
}

func OpenServiceInBrowser(wf *aw.Workflow, awsService *awsworkflow.AwsService) {
	cmd := exec.Command("open", awsService.Url)
	if err := wf.RunInBackground("open-sub-service-in-browser-"+awsService.Id, cmd); err != nil {
		panic(err)
	}
	time.Sleep(250 * time.Millisecond) // sleep so that tabs are more-or-less opened in the order by which this function is called
}
