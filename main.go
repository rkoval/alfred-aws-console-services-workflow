package main

import (
	"flag"
	"log"
	"strings"

	aw "github.com/deanishe/awgo"
	"github.com/deanishe/awgo/update"
	"github.com/rkoval/alfred-aws-console-services-workflow/awsworkflow"
	"github.com/rkoval/alfred-aws-console-services-workflow/workflow"
)

var wf *aw.Workflow
var forceFetch bool
var query string
var ymlPath string
var openAll bool

func init() {
	flag.BoolVar(&forceFetch, "fetch", false, "force fetch via AWS instead of cache")
	flag.BoolVar(&openAll, "open_all", false, "open all URLs in a browser for the matching query")
	flag.StringVar(&query, "query", "", "query to use")
	flag.StringVar(&ymlPath, "yml_path", "console-services.yml", "query to use")
	flag.Parse()
	wf = aw.New(update.GitHub("rkoval/alfred-aws-console-services-workflow"))
}

func main() {
	wf.Run(func() {
		log.Printf("running workflow with query: `%s`", query)
		cfg := awsworkflow.NewWorkflowConfig(nil)
		query = strings.TrimLeft(query, " ")

		workflow.Run(wf, query, cfg, forceFetch, openAll, ymlPath)
	})
}
