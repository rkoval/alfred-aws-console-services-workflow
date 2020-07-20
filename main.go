package main

import (
	"io/ioutil"
	"log"
	"strings"

	aw "github.com/deanishe/awgo"
	"gopkg.in/yaml.v2"
)

var wf *aw.Workflow

func init() {
	wf = aw.New()
}

type AwsService struct {
	Id               string   `yaml:"id"`
	Name             string   `yaml:"name"`
	ShortName        string   `yaml:"short_name"`
	Description      string   `yaml:"description"`
	Url              string   `yaml:"url"`
	ExtraSearchTerms []string `yaml:"extra_search_terms"`
}

func parseYaml() []AwsService {
	awsServices := []AwsService{}
	yamlFile, err := ioutil.ReadFile("console-services.yml")
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(yamlFile, &awsServices)
	if err != nil {
		log.Fatal(err)
	}
	return awsServices
}

func run() {
	var query string
	if args := wf.Args(); len(args) > 0 {
		query = args[0]
	}

	awsServices := parseYaml()
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
			UID(awsService.Id).
			Arg(awsService.Url).
			Subtitle(awsService.Description).
			Match(match).
			Valid(true)

		icon := &aw.Icon{Value: "images/" + awsService.Id + ".png"}
		item.Icon(icon)
	}

	if query != "" {
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
