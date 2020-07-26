package workflow

import (
	"strings"

	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awsworkflow"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchtypes"
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

func AddServiceToWorkflow(wf *aw.Workflow, awsService awsworkflow.AwsService) {
	title := awsService.Id
	match := awsService.Id + " " + awsService.ShortName + " " + awsService.Name
	subtitle := ""
	if len(awsService.SubServices) > 0 {
		subtitle += "🗂 "
	}

	subtitle += awsService.Name
	if awsService.ShortName != "" {
		subtitle += " (" + awsService.ShortName + ")"
	}
	subtitle += " – " + awsService.Description

	if len(awsService.ExtraSearchTerms) > 0 {
		match += " " + strings.Join(awsService.ExtraSearchTerms, " ")
	}

	item := util.NewURLItem(wf, title).
		Subtitle(subtitle).
		Match(match).
		Autocomplete(awsService.Id + " ").
		UID(awsService.Id).
		Arg(awsService.Url)

	item.Icon(awsworkflow.GetImageIcon(awsService.Id))
}

func SearchServices(wf *aw.Workflow, awsServices []awsworkflow.AwsService) {
	for _, awsService := range awsServices {
		AddServiceToWorkflow(wf, awsService)
	}
}

func SearchSubServices(wf *aw.Workflow, awsService awsworkflow.AwsService) {
	for _, subService := range awsService.SubServices {
		title := awsService.Id + " " + subService.Id
		subtitle := ""

		searchType := searchtypes.SearchTypesByServiceId[awsService.Id+"_"+subService.Id]
		if searchType != searchtypes.None {
			// this subservice has a searcher, so denote that in the result
			if searchType == searchtypes.SearchTypesByServiceId[awsService.Id] {
				// this sub-service is the default searcher
				subtitle += "🔎⭐️ "
			} else {
				subtitle += "🔎 "
			}
		}

		if awsService.ShortName != "" {
			subtitle += awsService.ShortName + " – "
		} else {
			subtitle += awsService.GetName() + " – "
		}

		var match string
		if subService.Id == "home" {
			subtitle += "Home"
			match = subService.Id
		} else {
			subtitle += subService.Name
			match = subService.Id + " " + subService.Name
		}

		if subService.Description != "" {
			subtitle += " – " + subService.Description
		}

		item := util.NewURLItem(wf, title).
			Subtitle(subtitle).
			Match(match).
			Autocomplete(awsService.Id + " " + subService.Id + " ").
			UID(subService.Id).
			Arg(subService.Url)

		item.Icon(awsworkflow.GetImageIcon(awsService.Id))
	}
}
