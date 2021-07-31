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
	"github.com/rkoval/alfred-aws-console-services-workflow/aliases"
	"github.com/rkoval/alfred-aws-console-services-workflow/awsworkflow"
	"github.com/rkoval/alfred-aws-console-services-workflow/parsers"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

func Run(wf *aw.Workflow, rawQuery string, cfg aws.Config, forceFetch, openAll bool, ymlPath string) {
	log.Println("using workflow cacheDir: " + wf.CacheDir())
	log.Println("using workflow dataDir: " + wf.DataDir())

	parser := parsers.NewParser(rawQuery)
	query, awsServices := parser.Parse(ymlPath)
	defer finalize(wf, query)

	searchArgs := searchutil.SearchArgs{
		Cfg:        cfg,
		ForceFetch: forceFetch,
		FullQuery:  rawQuery,
		Profile:    util.GetProfile(cfg),
	}

	log.Printf("using query: %#v", query)

	if query.IsEmpty() {
		handleEmptyQuery(wf, searchArgs)
		return
	}

	if query.HasOpenAll {
		handleOpenAll(wf, query.Service, awsServices, openAll, rawQuery, cfg)
		return
	}

	// TODO need to autocomplete correctly within services/subservices/searchers with respect to region

	if query.RegionQuery != nil {
		for _, region := range awsworkflow.AllAWSRegions {
			autocomplete := strings.Replace(rawQuery, aliases.OverrideAwsRegion+*query.RegionQuery, aliases.OverrideAwsRegion+region.Name+" ", 1)
			wf.NewItem(region.Name).
				Subtitle(region.Description).
				Icon(aw.IconWeb).
				Autocomplete(autocomplete).
				UID(region.Name)
		}
		log.Printf("filtering with region override %q", *query.RegionQuery)
		res := wf.Filter(*query.RegionQuery)
		log.Printf("%d results match %q", len(res), *query.RegionQuery)
		return
	} else if query.RegionOverride != nil && (query.Service == nil || !query.Service.HasGlobalRegion) {
		cfg.Region = query.RegionOverride.Name
		searchArgs.Cfg.Region = query.RegionOverride.Name
	}

	if query.Service == nil || (!query.HasTrailingWhitespace && query.SubService == nil && !query.HasDefaultSearchAlias && query.RemainingQuery == "") {
		if query.Service == nil {
			searchArgs.Query = query.RemainingQuery
		} else if query.Service.ShortName != "" {
			searchArgs.Query = query.Service.ShortName
		} else {
			searchArgs.Query = query.Service.Name
		}
		log.Printf("using searcher associated with services with query %q", searchArgs.Query)
		SearchServices(wf, awsServices, searchArgs)
	} else {
		if !query.HasDefaultSearchAlias && (query.Service.SubServices == nil || len(query.Service.SubServices) <= 0) {
			handleUnimplemented(wf, query.Service, nil, fmt.Sprintf("%s doesn't have sub-services configured (yet)", query.Service.Id), searchArgs)
			return
		}

		if query.HasDefaultSearchAlias || query.SubService != nil && (query.HasTrailingWhitespace || query.RemainingQuery != "") {
			serviceId := query.Service.Id
			if query.SubService != nil {
				serviceId += "_" + query.SubService.Id
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
				handleUnimplemented(wf, query.Service, query.SubService, fmt.Sprintf("No searcher for `%s %s` (yet)", query.Service.Id, query.SubService.Id), searchArgs)
				return
			}
		} else {
			log.Println("using searcher associated with sub-services")
			if query.SubService != nil {
				searchArgs.Query = query.SubService.Id
			} else {
				searchArgs.Query = query.RemainingQuery
			}
			SearchSubServices(wf, *query.Service, searchArgs)
		}
	}

	if searchArgs.Query != "" {
		log.Printf("filtering with query %q", searchArgs.Query)
		res := wf.Filter(searchArgs.Query)
		log.Printf("%d results match %q", len(res), searchArgs.Query)
	}
}

func finalize(wf *aw.Workflow, query *parsers.Query) {
	if wf.IsEmpty() {
		title := ""
		subtitle := ""
		if query.RegionQuery != nil {
			title = "No matching regions found"
			subtitle = "Try starting over with `$` again to see the full list"
		} else {
			title = "No matching services found"
			subtitle = "Try another query (example: `aws ec2 instances`)"
		}
		wf.NewItem(title).
			Subtitle(subtitle).
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
		util.NewURLItem(wf, "No region configured for this profile").
			Subtitle("Select this option to open AWS docs on how to configure").
			Arg("https://aws.github.io/aws-sdk-go-v2/docs/configuring-sdk/#creating-the-config-file").
			Icon(aw.IconWarning)
	} else {
		wf.NewItem("Using region \"" + searchArgs.Cfg.Region + "\"").
			Subtitle("Use \"" + aliases.OverrideAwsRegion + "\" to override for the current query").
			Icon(aw.IconWeb)
	}

	if wf.UpdateCheckDue() {
		if err := wf.CheckForUpdate(); err != nil {
			wf.FatalError(err)
		}
	}
	handleUpdateAvailable(wf)
}

func handleUnimplemented(wf *aw.Workflow, awsService, subService *awsworkflow.AwsService, header string, searchArgs searchutil.SearchArgs) {
	searchArgs.IgnoreAutocompleteTerm = true
	if subService == nil {
		AddServiceToWorkflow(wf, *awsService, searchArgs)
	} else {
		AddSubServiceToWorkflow(wf, *awsService, *subService, searchArgs)
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
