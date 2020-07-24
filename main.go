package main

import (
	"flag"
	"log"
	"strings"

	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/core"
	"github.com/rkoval/alfred-aws-console-services-workflow/workflow"
)

var wf *aw.Workflow
var forceFetch bool
var query string
var ymlPath string

func init() {
	flag.BoolVar(&forceFetch, "fetch", false, "force fetch via AWS instead of cache")
	flag.StringVar(&query, "query", "", "query to use")
	flag.StringVar(&ymlPath, "yml_path", "console-services.yml", "query to use")
	flag.Parse()
	wf = aw.New()
}

func main() {
	wf.Run(func() {
		log.Printf("running workflow with query: `%s`", query)
		session := core.LoadAWSConfig(nil)
		query = strings.TrimLeft(query, " ")

		workflow.Run(wf, query, session, forceFetch, ymlPath)
	})
}
