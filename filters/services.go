package filters

import (
	"strings"

	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/core"
)

func Services(wf *aw.Workflow, awsServices []core.AwsService, query string) {
	for _, awsService := range awsServices {
		var title string
		var match string
		if awsService.ShortName != "" {
			title = awsService.ShortName + " - " + awsService.Name
			match = awsService.ShortName
		} else {
			title = awsService.Name
			match = title
		}

		if len(awsService.ExtraSearchTerms) > 0 {
			match += " " + strings.Join(awsService.ExtraSearchTerms, " ")
		}

		item := wf.NewItem(title).
			Autocomplete(awsService.Id).
			UID(awsService.Id).
			Arg(awsService.Url).
			Subtitle(awsService.Description).
			Match(match).
			Valid(true)

		icon := &aw.Icon{Value: "images/" + awsService.Id + ".png"}
		item.Icon(icon)
	}
}

func ServiceSections(wf *aw.Workflow, awsService core.AwsService, subquery string) {
	for _, section := range awsService.Sections {
		item := wf.NewItem(awsService.GetName() + " - " + section.Name).
			Autocomplete(section.Id).
			UID(section.Id).
			Arg(section.Url).
			Subtitle(section.Description).
			Valid(true)

		icon := &aw.Icon{Value: "images/" + awsService.Id + ".png"}
		item.Icon(icon)
	}
}
