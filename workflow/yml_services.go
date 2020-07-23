package workflow

import (
	"strings"

	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/core"
)

func AddServiceToWorkflow(wf *aw.Workflow, awsService core.AwsService) {
	title := awsService.GetName()
	var match string
	if awsService.ShortName != "" {
		match = awsService.ShortName
	} else {
		match = title
	}

	if len(awsService.ExtraSearchTerms) > 0 {
		match += " " + strings.Join(awsService.ExtraSearchTerms, " ")
	}

	item := wf.NewItem(title).
		Autocomplete(awsService.Id + " ").
		UID(awsService.Id).
		Arg(awsService.Url).
		Subtitle(awsService.Description).
		Match(match).
		Valid(true)

	icon := &aw.Icon{Value: "images/" + awsService.Id + ".png"}
	item.Icon(icon)
}

func SearchServices(wf *aw.Workflow, awsServices []core.AwsService) {
	for _, awsService := range awsServices {
		AddServiceToWorkflow(wf, awsService)
	}
}

func SearchSubServices(wf *aw.Workflow, awsService core.AwsService) {
	for _, subService := range awsService.SubServices {
		var title string
		if subService.Id == "home" {
			title = awsService.GetName() + " – Home"
		} else if awsService.ShortName != "" {
			title = awsService.ShortName + " – " + subService.Name
		} else {
			title = awsService.GetName() + " – " + subService.Name
		}

		item := wf.NewItem(title).
			Autocomplete(awsService.Id + " " + subService.Id + " ").
			UID(subService.Id).
			Arg(subService.Url).
			Subtitle(subService.Description).
			Valid(true)

		item.Icon(core.GetImageIcon(awsService.Id))
	}
}
