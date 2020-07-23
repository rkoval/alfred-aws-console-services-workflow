package main

import (
	"flag"
	"log"
	"strings"

	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/workflow"
)

var wf *aw.Workflow
var forceFetch bool
var query string

func init() {
	flag.BoolVar(&forceFetch, "fetch", false, "force fetch via AWS instead of cache")
	flag.StringVar(&query, "query", "", "query to use")
	flag.Parse()
	wf = aw.New()
}

func main() {
	wf.Configure(aw.TextErrors(true))
	wf.Run(func() {
		log.Printf("running workflow with query: `%s`", query)
		query = strings.TrimLeft(query, " ")

		workflow.Run(wf, query, nil, forceFetch, "console-services.yml")
	})
}
