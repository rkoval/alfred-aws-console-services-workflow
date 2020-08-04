package caching

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

var jobName = "fetch"

func handleExpiredCache(wf *aw.Workflow, cacheName string, lastFetchErrPath string, rawQuery string) error {
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
	if wf.Cache.Expired(cacheName, maxCacheAge) {
		log.Printf("cache with key `%s` was expired (older than %d seconds) in %s", cacheName, maxCacheAgeSeconds, wf.CacheDir())
		wf.Rerun(0.5)
		if !wf.IsRunning(jobName) {
			cmd := exec.Command(os.Args[0], "-query="+rawQuery+"", "-fetch")
			log.Printf("running `%s` in background as job `%s` ...", cmd, jobName)
			if err := wf.RunInBackground(jobName, cmd); err != nil {
				panic(err)
			}
		} else {
			log.Printf("background job `%s` already running", jobName)
		}

		return handleFetchErr(wf, lastFetchErrPath)
	}

	return nil
}

func handleFetchErr(wf *aw.Workflow, lastFetchErrPath string) error {
	data, err := ioutil.ReadFile(lastFetchErrPath)
	if err != nil {
		if !os.IsNotExist(err) {
			// this file will often not exist, so don't spam logs if it doesn't
			log.Println(err)
		}
		return nil
	}

	errString := string(data)
	wf.Configure(aw.SuppressUIDs(true))
	if strings.HasPrefix(errString, "NoCredentialProviders") {
		util.NewURLItem(wf, "AWS credentials not configured in ~/.aws/credentials").
			Subtitle("Press enter to open AWS docs for how to configure").
			Arg("https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html#specifying-credentials").
			Icon(aw.IconError).
			Valid(true)
	} else if strings.HasPrefix(errString, "MissingRegion") {
		util.NewURLItem(wf, "AWS default region not configured in ~/.aws/config").
			Subtitle("Press enter to open AWS docs for how to configure").
			Arg("https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html#specifying-the-region").
			Icon(aw.IconError).
			Valid(true)
	} else {
		wf.NewItem(errString).
			Icon(aw.IconError)
	}

	return errors.New(errString)
}
