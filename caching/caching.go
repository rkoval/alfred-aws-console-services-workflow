package caching

import (
	"log"
	"os"
	"os/exec"
	"strconv"
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

	maxCacheAgeSeconds := 180
	m := os.Getenv("ALFRED_AWS_CONSOLE_SERVICES_WORKFLOW_MAX_CACHE_AGE_SECONDS")
	if m != "" {
		converted, err := strconv.Atoi(m)
		if err != nil {
			panic(err)
		}
		if converted != 0 {
			log.Printf("using custom max cache age of %v seconds", converted)
			maxCacheAgeSeconds = converted
		}
	}

	maxCacheAge := time.Duration(maxCacheAgeSeconds) * time.Second
	if wf.Cache.Exists(cacheName) {
		log.Printf("using cache with key `%s` in %s ...", cacheName, wf.CacheDir())
		if err := wf.Cache.LoadJSON(cacheName, &results); err != nil {
			panic(err)
		}
	} else {
		log.Printf("cache with key `%s` did not exist in %s ...", cacheName, wf.CacheDir())
		wf.NewItem("Fetching ...").
			Icon(aw.IconInfo)
	}

	if wf.Cache.Expired(cacheName, maxCacheAge) {
		log.Printf("cache with key `%s` was expired (older than %d seconds) in %s", cacheName, maxCacheAge, wf.CacheDir())
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
	}

	return results
}
