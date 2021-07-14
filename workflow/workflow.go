package workflow

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awsworkflow"
	"github.com/rkoval/alfred-aws-console-services-workflow/parsers"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

func Run(wf *aw.Workflow, rawQuery string, cfg aws.Config, forceFetch, openAll bool, ymlPath string) {
	log.Println("using workflow cacheDir: " + wf.CacheDir())
	log.Println("using workflow dataDir: " + wf.DataDir())

	awsServices := parsers.ParseConsoleServicesYml(ymlPath)
	parser := parsers.NewParser(strings.NewReader(rawQuery))
	query := parser.Parse()
	defer finalize(wf)

	searchArgs := searchutil.SearchArgs{
		Cfg:        cfg,
		ForceFetch: forceFetch,
		FullQuery:  rawQuery,
		Profile:    util.GetProfile(cfg),
	}

	if query.IsEmpty() {
		handleEmptyQuery(wf, searchArgs)
		return
	}

	var awsService *awsworkflow.AwsService
	if query.ServiceId != "" {
		for i := range awsServices {
			if awsServices[i].Id == query.ServiceId {
				awsService = &awsServices[i]
				break
			}
		}
	}

	if query.HasOpenAll {
		handleOpenAll(wf, awsService, awsServices, openAll, rawQuery, cfg)
		return
	}

	if awsService == nil || (!query.HasTrailingWhitespace && query.SubServiceId == "" && !query.HasDefaultSearchAlias) {
		log.Println("using searcher associated with services")
		searchArgs.Query = query.ServiceId
		SearchServices(wf, awsServices, cfg)
	} else {
		if !query.HasDefaultSearchAlias && (awsService.SubServices == nil || len(awsService.SubServices) <= 0) {
			handleUnimplemented(wf, awsService, nil, fmt.Sprintf("%s doesn't have sub-services configured (yet)", awsService.Id), cfg)
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
				searchArgs.Query = query.RemainingQuery
				err := searcher.Search(wf, searchArgs)
				if err != nil {
					wf.FatalError(err)
				}
			} else {
				handleUnimplemented(wf, awsService, subService, fmt.Sprintf("No searcher for `%s %s` (yet)", query.ServiceId, query.SubServiceId), cfg)
				return
			}
		} else {
			log.Println("using searcher associated with sub-services")
			searchArgs.Query = query.SubServiceId
			SearchSubServices(wf, *awsService, cfg)
		}
	}

	if searchArgs.Query != "" {
		log.Printf("filtering with query %s", searchArgs.Query)
		res := wf.Filter(searchArgs.Query)
		log.Printf("%d results match %q", len(res), searchArgs.Query)
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

func handleEmptyQuery(wf *aw.Workflow, searchArgs searchutil.SearchArgs) {
	log.Println("no search type parsed")
	wf.NewItem("Search for an AWS Service ...").
		Subtitle("e.g., cloudformation, ec2, s3 ...")

	if searchArgs.Profile != "" {
		wf.NewItem("Using profile \"" + searchArgs.Profile + "\"").
			Icon(aw.IconAccount)
	}

	if searchArgs.Cfg.Region == "" {
		wf.NewItem("No region configured for this profile").
			Subtitle("Select this option to open AWS docs on how to configure").
			Arg("https://aws.github.io/aws-sdk-go-v2/docs/configuring-sdk/#creating-the-config-file").
			Icon(aw.IconWarning)
	} else {
		wf.NewItem("Using region \"" + searchArgs.Cfg.Region + "\"").
			Icon(aw.IconWeb)
	}

	if wf.UpdateCheckDue() {
		if err := wf.CheckForUpdate(); err != nil {
			wf.FatalError(err)
		}
	}
	handleUpdateAvailable(wf)
}

func handleUnimplemented(wf *aw.Workflow, awsService, subService *awsworkflow.AwsService, header string, cfg aws.Config) {
	if subService == nil {
		AddServiceToWorkflow(wf, *awsService, cfg)
	} else {
		AddSubServiceToWorkflow(wf, *awsService, *subService, cfg)
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

func handleOpenAll(wf *aw.Workflow, awsService *awsworkflow.AwsService, allAwsServices []awsworkflow.AwsService, openAll bool, rawQuery string, cfg aws.Config) {
	if openAll {
		if awsService == nil {
			for _, awsService := range allAwsServices {
				openServiceInBrowser(wf, &awsService, cfg)
			}
		} else {
			for _, subService := range awsService.SubServices {
				openServiceInBrowser(wf, &subService, cfg)
			}
		}
	} else {
		var length int
		var entityName string
		if awsService == nil {
			length = len(allAwsServices)
			entityName = "services"
		} else {
			length = len(awsService.SubServices)
			entityName = awsService.Id + " sub-services"
		}

		title := fmt.Sprintf("Open the %d %s in browser", length, entityName)
		cmd := fmt.Sprintf(`%s -query="%s" -open_all`, os.Args[0], rawQuery)
		wf.NewItem(title).
			Subtitle("Fair warning: this may briefly overload your system").
			Icon(aw.IconWarning).
			Arg(cmd).
			Var("action", "run-script").
			Valid(true)
	}
}

func openServiceInBrowser(wf *aw.Workflow, awsService *awsworkflow.AwsService, cfg aws.Config) {
	cmd := exec.Command("open", util.ConstructAWSConsoleUrl(awsService.Url, awsService.GetRegion(cfg)))
	if err := wf.RunInBackground("open-service-in-browser-"+awsService.Id, cmd); err != nil {
		panic(err)
	}
	time.Sleep(250 * time.Millisecond) // sleep so that tabs are more-or-less opened in the order by which this function is called
}
