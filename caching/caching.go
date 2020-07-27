package caching

import (
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	aw "github.com/deanishe/awgo"

	"github.com/cheekybits/genny/generic"
)

//go:generate genny -in=$GOFILE -out=gen-$GOFILE gen "Entity=ec2.Instance,s3.Bucket,ec2.SecurityGroup,elasticbeanstalk.EnvironmentDescription"
type Entity = generic.Type

type EntityArrayFetcher = func(*session.Session) ([]Entity, error)

func LoadEntityArrayFromCache(wf *aw.Workflow, session *session.Session, cacheName string, fetcher EntityArrayFetcher, forceFetch bool, fullQuery string) []Entity {
	if *session.Config.Region == "" {
		panic(aws.ErrMissingRegion)
	}
	// TODO optimization: not all services have sa region associated with them, so cache can be reused across regions (e.g., s3 buckets are global)
	cacheName += "_" + *session.Config.Region

	results := []Entity{}
	var jobName = "fetch"
	if forceFetch {
		log.Printf("fetching from aws ...")
		results, err := fetcher(session)
		if err != nil {
			panic(err)
		}
		log.Printf("fetched %d results from aws", len(results))

		log.Printf("storing %d results with cache key `%s` to %s ...", len(results), cacheName, wf.CacheDir())
		if err := wf.Cache.StoreJSON(cacheName, results); err != nil {
			panic(err)
		}
		return results
	}

	if wf.Cache.Exists(cacheName) {
		log.Printf("using cache with key `%s` from %s ...", cacheName, wf.CacheDir())
		if err := wf.Cache.LoadJSON(cacheName, &results); err != nil {
			panic(err)
		}
		return results
	}

	maxCacheAge := 1 * time.Minute
	if wf.Cache.Expired(cacheName, maxCacheAge) {
		log.Printf("cache with key `%s` did not exist or was expired in %s", cacheName, wf.CacheDir())
		wf.Rerun(0.4)
		if !wf.IsRunning(jobName) {
			cmd := exec.Command(os.Args[0], "-query="+fullQuery+"", "-fetch")
			log.Printf("running `%s` in background as job `%s` ...", cmd, jobName)
			if err := wf.RunInBackground(jobName, cmd); err != nil {
				panic(err)
			}
		} else {
			log.Printf("background job `%s` already running", jobName)
		}
		wf.NewItem("Fetching ...").
			Icon(aw.IconInfo)
	}
	return results
}
