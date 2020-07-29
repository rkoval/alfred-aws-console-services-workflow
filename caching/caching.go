package caching

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	aw "github.com/deanishe/awgo"

	"github.com/cheekybits/genny/generic"
)

//go:generate genny -in=$GOFILE -out=gen-$GOFILE gen "Entity=ec2.Instance,s3.Bucket,ec2.SecurityGroup,elasticbeanstalk.EnvironmentDescription,wafv2.IPSetSummary"
type Entity = generic.Type

type EntityArrayFetcher = func(*session.Session) ([]Entity, error)

func LoadEntityArrayFromCache(wf *aw.Workflow, session *session.Session, cacheName string, fetcher EntityArrayFetcher, forceFetch bool, fullQuery string) []Entity {
	// TODO optimization: not all services have sa region associated with them, so cache can be reused across regions (e.g., s3 buckets are global)
	cacheName += "_" + *session.Config.Region

	results := []Entity{}
	lastFetchErrPath := wf.CacheDir() + "/last-fetch-err.txt"
	if forceFetch {
		log.Printf("fetching from aws ...")
		results, err := fetcher(session)

		if err != nil {
			log.Printf("fetch error occurred. writing to %s ...", lastFetchErrPath)
			_ = ioutil.WriteFile(lastFetchErrPath, []byte(err.Error()), 0600)
			panic(err)
		} else {
			os.Remove(lastFetchErrPath)
		}
		log.Printf("fetched %d results from aws", len(results))

		log.Printf("storing %d results with cache key `%s` to %s ...", len(results), cacheName, wf.CacheDir())
		if err := wf.Cache.StoreJSON(cacheName, results); err != nil {
			panic(err)
		}
		return results
	}

	err := handleExpiredCache(wf, cacheName, lastFetchErrPath, fullQuery)
	if err != nil {
		return []Entity{}
	}

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

	return results
}
