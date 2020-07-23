package core

import (
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	aw "github.com/deanishe/awgo"
)

var maxCacheAge = 5 * time.Minute
var jobName = "fetch"

type Fetcher func(http.RoundTripper) ([]interface{}, error)

func HandleCache(wf *aw.Workflow, transport http.RoundTripper, cacheName string, results *[]interface{}, fetcher Fetcher, forceFetch bool, fullQuery string) {
	if forceFetch {
		wf.Configure(aw.TextErrors(true))
		log.Printf("fetching from aws ...")
		rs, err := fetcher(transport)
		log.Printf("fetched %d results from aws", len(rs))

		if err != nil {
			panic(err)
		}
		log.Printf("storing %d results in cache key `%s` ...", len(rs), cacheName)
		if err := wf.Cache.StoreJSON(cacheName, rs); err != nil {
			panic(err)
		}
		for i := range rs {
			*results = append(*results, rs[i])
		}
		return
	}

	if wf.Cache.Exists(cacheName) {
		log.Printf("using cache with key `%s` ...", cacheName)
		if err := wf.Cache.LoadJSON(cacheName, results); err != nil {
			panic(err)
		}
	}

	if wf.Cache.Expired(cacheName, maxCacheAge) {
		wf.Rerun(0.2)
		if !wf.IsRunning(jobName) {
			cmd := exec.Command(os.Args[0], "-fetch", "-query='"+fullQuery+"'")
			log.Printf("running `%s` in background as job `%s` ...", cmd, jobName)
			if err := wf.RunInBackground(jobName, cmd); err != nil {
				panic(err)
			}
		} else {
			log.Printf("background job `%s` already running", jobName)
		}
		// Cache is also "expired" if it doesn't exist. So if there are no
		// cached data, show a corresponding message and exit.
		if len(*results) == 0 {
			wf.NewItem("Fetching ...").
				Icon(aw.IconInfo)
			return
		}
	}
}
