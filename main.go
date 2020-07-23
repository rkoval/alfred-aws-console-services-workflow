package main

import (
	"log"

	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/parsers"
)

var wf *aw.Workflow

func init() {
	wf = aw.New()
}

func run() {
	awsServices := parsers.ParseConsoleServicesYml()
	query, err := parsers.ParseQueryAndPopulateItems(wf, awsServices)

	if err != nil {
		wf.FatalError(err)
	} else if query != "" {
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

func main() {
	wf.Run(run)
}
