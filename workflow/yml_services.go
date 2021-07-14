package workflow

import (
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awsworkflow"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers"
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

func AddServiceToWorkflow(wf *aw.Workflow, awsService awsworkflow.AwsService, cfg aws.Config) {
	title := awsService.Id
	match := awsService.Id
	if awsService.ShortName == "" {
		match += " " + awsService.Name
	}

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
		Arg(util.ConstructAWSConsoleUrl(awsService.Url, awsService.GetRegion(cfg)))

	item.Icon(awsworkflow.GetImageIcon(awsService.Id))
}

func SearchServices(wf *aw.Workflow, awsServices []awsworkflow.AwsService, cfg aws.Config) {
	for _, awsService := range awsServices {
		AddServiceToWorkflow(wf, awsService, cfg)
	}
}

func SearchSubServices(wf *aw.Workflow, awsService awsworkflow.AwsService, cfg aws.Config) {
	for _, subService := range awsService.SubServices {
		AddSubServiceToWorkflow(wf, awsService, subService, cfg)
	}
}

func AddSubServiceToWorkflow(wf *aw.Workflow, awsService, subService awsworkflow.AwsService, cfg aws.Config) {
	title := awsService.Id + " " + subService.Id
	subtitle := ""

	searcher := searchers.SearchersByServiceId[awsService.Id+"_"+subService.Id]
	if searcher != nil {
		// this subservice has a searcher, so denote that in the result
		if searcher == searchers.SearchersByServiceId[awsService.Id] {
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

	subtitle += subService.Name
	match := subService.Id + " " + subService.Name

	if subService.Description != "" {
		subtitle += " – " + subService.Description
	}

	item := util.NewURLItem(wf, title).
		Subtitle(subtitle).
		Match(match).
		Autocomplete(awsService.Id + " " + subService.Id + " ").
		UID(subService.Id).
		Arg(util.ConstructAWSConsoleUrl(subService.Url, awsService.GetRegion(cfg)))

	item.Icon(awsworkflow.GetImageIcon(awsService.Id))
}
