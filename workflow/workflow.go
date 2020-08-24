package workflow

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awsworkflow"
	"github.com/rkoval/alfred-aws-console-services-workflow/parsers"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers"
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

func Run(wf *aw.Workflow, rawQuery string, session *session.Session, forceFetch, openAll bool, ymlPath string) {
	log.Println("using workflow cacheDir: " + wf.CacheDir())
	log.Println("using workflow dataDir: " + wf.DataDir())

	awsServices := parsers.ParseConsoleServicesYml(ymlPath)
	parser := parsers.NewParser(strings.NewReader(rawQuery))
	query := parser.Parse()
	defer finalize(wf)

	if query.IsEmpty() {
		handleEmptyQuery(wf)
		return
	}

	var awsService *awsworkflow.AwsService
	for i := range awsServices {
		if awsServices[i].Id == query.ServiceId {
			awsService = &awsServices[i]
			break
		}
	}

	if query.HasOpenAll {
		handleOpenAll(wf, awsService, openAll, rawQuery)
		return
	}

	var filterQuery string
	if awsService == nil || (!query.HasTrailingWhitespace && query.SubServiceId == "" && !query.HasDefaultSearchAlias) {
		log.Println("using searcher associated with services")
		filterQuery = query.ServiceId
		SearchServices(wf, awsServices)
	} else {
		if !query.HasDefaultSearchAlias && (awsService.SubServices == nil || len(awsService.SubServices) <= 0) {
			handleUnimplemented(wf, awsService, nil, fmt.Sprintf("%s doesn't have sub-services configured (yet)", awsService.Id))
			return
		}

		var subService *awsworkflow.AwsService
		for i := range awsService.SubServices {
			if awsService.SubServices[i].Id == query.SubServiceId {
				subService = &awsService.SubServices[i]
				break
			}
		}

		if query.HasDefaultSearchAlias || subService != nil && (query.HasTrailingWhitespace || query.RemainingQuery != "") {
			serviceId := query.ServiceId
			if query.SubServiceId != "" {
				serviceId += "_" + query.SubServiceId
			}
			log.Println("using searcher associated with " + serviceId)
			searcher := searchers.SearchersByServiceId[serviceId]
			if searcher != nil {
				filterQuery = query.RemainingQuery
				err := searcher.Search(wf, filterQuery, session, forceFetch, rawQuery)
				if err != nil {
					wf.FatalError(err)
				}
			} else {
				handleUnimplemented(wf, awsService, subService, fmt.Sprintf("No searcher for `%s %s` (yet)", query.ServiceId, query.SubServiceId))
				return
			}
		} else {
			log.Println("using searcher associated with sub-services")
			filterQuery = query.SubServiceId
			SearchSubServices(wf, *awsService)
		}
	}

	if filterQuery != "" {
		log.Printf("filtering with query %s", filterQuery)
		res := wf.Filter(filterQuery)
		log.Printf("%d results match %q", len(res), filterQuery)
	}
}

func finalize(wf *aw.Workflow) {
	if wf.IsEmpty() {
		wf.NewItem("No matching services found").
			Subtitle("Try another query (example: `aws ec2 instances`)").
			Icon(aw.IconNote)
		handleUpdateAvailable(wf)
	}
	wf.SendFeedback()
}

func handleEmptyQuery(wf *aw.Workflow) {
	log.Println("no search type parsed")
	wf.NewItem("Search for an AWS Service ...").
		Subtitle("e.g., cloudformation, ec2, s3 ...")

	if wf.UpdateCheckDue() {
		if err := wf.CheckForUpdate(); err != nil {
			wf.FatalError(err)
		}
	}
	handleUpdateAvailable(wf)
}

func handleUnimplemented(wf *aw.Workflow, awsService, subService *awsworkflow.AwsService, header string) {
	if subService == nil {
		AddServiceToWorkflow(wf, *awsService)
	} else {
		AddSubServiceToWorkflow(wf, *awsService, *subService)
	}
	util.NewURLItem(wf, header).
		Subtitle("Select this result to open the contributing guide to easily add them!").
		Arg("https://github.com/rkoval/alfred-aws-console-services-workflow/blob/master/CONTRIBUTING.md").
		Icon(aw.IconNote)
	handleUpdateAvailable(wf)
}

func handleUpdateAvailable(wf *aw.Workflow) {
	if wf.UpdateAvailable() {
		util.NewURLItem(wf, "Update available").
			Subtitle("Select this result to navigate to download").
			Arg("https://github.com/rkoval/alfred-aws-console-services-workflow/releases").
			Icon(aw.IconInfo)
	}
}

func handleOpenAll(wf *aw.Workflow, awsService *awsworkflow.AwsService, openAll bool, rawQuery string) {
	if openAll {
		for i := range awsService.SubServices {
			openServiceInBrowser(wf, &awsService.SubServices[i])
		}
	} else if awsService != nil {
		cmd := fmt.Sprintf(`%s -query="%s" -open_all`, os.Args[0], rawQuery)
		wf.NewItem(fmt.Sprintf("Open the %d %s sub-services in browser", len(awsService.SubServices), awsService.Id)).
			Subtitle("Fair warning: this may briefly overload your system").
			Icon(aw.IconWarning).
			Arg(cmd).
			Var("action", "run-script").
			Valid(true)
	}
}

func openServiceInBrowser(wf *aw.Workflow, awsService *awsworkflow.AwsService) {
	cmd := exec.Command("open", awsService.Url)
	if err := wf.RunInBackground("open-sub-service-in-browser-"+awsService.Id, cmd); err != nil {
		panic(err)
	}
	time.Sleep(250 * time.Millisecond) // sleep so that tabs are more-or-less opened in the order by which this function is called
}
