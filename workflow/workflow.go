package workflow

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/aliases"
	"github.com/rkoval/alfred-aws-console-services-workflow/awsconfig"
	"github.com/rkoval/alfred-aws-console-services-workflow/awsworkflow"
	"github.com/rkoval/alfred-aws-console-services-workflow/parsers"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

func Run(wf *aw.Workflow, rawQuery string, transport http.RoundTripper, forceFetch, openAll bool, ymlPath string) {
	log.Println("using workflow cacheDir: " + wf.CacheDir())
	log.Println("using workflow dataDir: " + wf.DataDir())

	parser := parsers.NewParser(rawQuery)
	query, awsServices := parser.Parse(ymlPath)
	defer finalize(wf, query)

	log.Printf("using query: %#v", query)

	if query.RegionQuery != nil {
		for _, region := range awsconfig.AllAWSRegions {
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
	}

	if query.ProfileQuery != nil {
		for _, profile := range awsconfig.GetAwsProfiles() {
			autocomplete := strings.Replace(rawQuery, aliases.OverrideAwsProfile+*query.ProfileQuery, aliases.OverrideAwsProfile+profile.Name+" ", 1)
			item := wf.NewItem(profile.Name).
				Icon(aw.IconAccount).
				Autocomplete(autocomplete).
				UID(profile.Name)

			if profile.Region != "" {
				item.Subtitle(fmt.Sprintf("🌎 %s", profile.Region))
			} else {
				item.Subtitle("⚠️ This profile does not specify a region. Functionality will be limited")
			}
		}
		log.Printf("filtering with profile override %q", *query.ProfileQuery)
		res := wf.Filter(*query.ProfileQuery)
		log.Printf("%d results match %q", len(res), *query.ProfileQuery)
		return
	}

	cfg := awsworkflow.InitAWS(transport, query.ProfileOverride, query.GetRegionOverride(), wf)
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

	if query.HasOpenAll {
		handleOpenAll(wf, query.Service, awsServices, openAll, rawQuery, cfg)
		return
	}

	if query.Service == nil || (!query.HasTrailingWhitespace && query.SubService == nil && !query.HasDefaultSearchAlias && query.RemainingQuery == "") {
		if query.Service == nil {
			searchArgs.Query = query.RemainingQuery
		} else if query.Service.ShortName != "" {
			searchArgs.Query = strings.ToLower(query.Service.ShortName)
		} else {
			searchArgs.Query = query.Service.Id
		}
		log.Printf("using searcher associated with services with query %q", searchArgs.Query)
		SearchServices(wf, awsServices, searchArgs)
	} else {
		if !query.HasDefaultSearchAlias && (len(query.Service.SubServices) <= 0) {
			handleUnimplemented(wf, query.Service, nil, fmt.Sprintf("%s doesn't have sub-services configured (yet)", query.Service.Id), searchArgs)
			return
		}

		searchArgs.GetRegionFunc = query.Service.GetRegion

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
	if r := recover(); r != nil {
		// we don't want this deferred function to prevent panics from showing in the workflow
		panic(r)
	}
	if wf.IsEmpty() {
		title := ""
		subtitle := ""
		if query.RegionQuery != nil {
			title = "No matching regions found"
			subtitle = fmt.Sprintf("Try starting over with \"%s\" again to see the full list", aliases.OverrideAwsRegion)
		} else if query.ProfileQuery != nil {
			title = "No matching profiles found"
			subtitle = fmt.Sprintf("Try starting over with \"%s\" again to see the full list", aliases.OverrideAwsProfile)
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

	if searchArgs.Profile == "" {
		util.NewURLItem(wf, "No profile configured").
			Subtitle("Select this option to open AWS docs on how to configure").
			Arg("https://aws.github.io/aws-sdk-go-v2/docs/configuring-sdk/#creating-the-config-file").
			Icon(aw.IconNote)
	} else {
		wf.NewItem("Using profile \"" + searchArgs.Profile + "\"").
			Subtitle("Use \"" + aliases.OverrideAwsProfile + "\" to override for the current query").
			Icon(aw.IconAccount)
	}

	if searchArgs.Cfg.Region == "" {
		util.NewURLItem(wf, "No region configured for this profile").
			Subtitle("Select this option to open AWS docs on how to configure").
			Arg("https://aws.github.io/aws-sdk-go-v2/docs/configuring-sdk/#creating-the-config-file").
			Icon(aw.IconNote)
	} else {
		wf.NewItem("Using region \"" + searchArgs.Cfg.Region + "\"").
			Subtitle("Use \"" + aliases.OverrideAwsRegion + "\" to override for the current query").
			Icon(aw.IconWeb)
	}

	util.NewURLItem(wf, "Like this workflow? Consider donating! 😻").
		Subtitle("Select this option to open this project's Patreon").
		Arg("https://www.patreon.com/rkoval_alfred_aws_console_services_workflow").
		Icon(aw.IconFavorite)

	util.NewURLItem(wf, "Developers gotta eat! 🍕").
		Subtitle("Select this option to buy Ryan a pizza").
		Arg("https://ryankoval.pizza").
		Icon(aw.IconColor)

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
		util.NewURLItem(wf, fmt.Sprintf("Update available (current version: %s)", wf.Version())).
			Subtitle("Select this result to navigate to download").
			Arg("https://github.com/rkoval/alfred-aws-console-services-workflow/releases").
			Icon(aw.IconInfo)
	}
}

func handleOpenAll(wf *aw.Workflow, awsService *awsworkflow.AwsService, allAwsServices []awsworkflow.AwsService, openAll bool, rawQuery string, cfg aws.Config) {
	if openAll {
		var urls []string
		if awsService == nil {
			for _, awsService := range allAwsServices {
				urls = append(urls, openServiceInBrowser(wf, awsService, awsService, cfg))
			}
		} else {
			for _, subService := range awsService.SubServices {
				urls = append(urls, openServiceInBrowser(wf, subService, *awsService, cfg))
			}
		}
		// write these to disk so that we can use this file to run automated validation against the pages as they redirect (or error) in AWS
		lastOpenedUrlsPath := wf.CacheDir() + "/last-opened-urls.json"
		urlBytes, err := json.Marshal(urls)
		if err != nil {
			panic(err)
		}
		err = os.WriteFile(lastOpenedUrlsPath, urlBytes, 0600)
		if err != nil {
			panic(err)
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
			Icon(aw.IconNote).
			Arg(cmd).
			Var("action", "run-script").
			Valid(true)
	}
}

func openServiceInBrowser(wf *aw.Workflow, awsService, awsServiceForRegion awsworkflow.AwsService, cfg aws.Config) string {
	url := util.ConstructAWSConsoleUrl(awsService.Url, awsServiceForRegion.GetRegion(cfg))
	log.Printf("opening url: %s ...", url)
	cmd := exec.Command("open", url)
	if err := wf.RunInBackground("open-service-in-browser-"+awsService.Id, cmd); err != nil {
		panic(err)
	}
	time.Sleep(250 * time.Millisecond) // sleep so that tabs are more-or-less opened in the order by which this function is called
	return url
}
