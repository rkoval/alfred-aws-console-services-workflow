package caching

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awsconfig"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

var jobName = "fetch"

func handleExpiredCache(wf *aw.Workflow, cacheName string, lastFetchErrPath string, searchArgs searchutil.SearchArgs) error {
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
			cmd := exec.Command(os.Args[0], "-query="+searchArgs.FullQuery+"", "-fetch")
			log.Printf("running `%s` in background as job `%s` ...", cmd, jobName)
			if err := wf.RunInBackground(jobName, cmd); err != nil {
				panic(err)
			}
		} else {
			log.Printf("background job `%s` already running", jobName)
		}

		return handleFetchErr(wf, lastFetchErrPath, searchArgs)
	}

	return nil
}

func handleFetchErr(wf *aw.Workflow, lastFetchErrPath string, searchArgs searchutil.SearchArgs) error {
	data, err := os.ReadFile(lastFetchErrPath)
	if err != nil {
		if !os.IsNotExist(err) {
			// this file will often not exist, so don't spam logs if it doesn't
			log.Println(err)
		}
		return nil
	}

	// TODO need to fix "no results" display when there's really a fetch error

	userHomePath := os.Getenv("HOME")
	errString := string(data)
	wf.Configure(aw.SuppressUIDs(true))
	var profileDescription string
	if searchArgs.Profile == "" {
		profileDescription = "for default profile"
	} else {
		profileDescription = "for profile \"" + searchArgs.Profile + "\""
	}
	if strings.HasPrefix(errString, "NoCredentialProviders") {
		credentialsFilePath := strings.Replace(awsconfig.GetAwsCredentialsFilePath(), userHomePath, "~", 1)
		util.NewURLItem(wf, "AWS credentials not set in "+credentialsFilePath+" "+profileDescription).
			Subtitle("Press enter to open AWS docs on how to configure").
			Arg("https://aws.github.io/aws-sdk-go-v2/docs/configuring-sdk/#creating-the-credentials-file").
			Icon(aw.IconError).
			Valid(true)
	} else if strings.HasPrefix(errString, "MissingRegion") {
		configFilePath := strings.Replace(awsconfig.GetAwsProfileFilePath(), userHomePath, "~", 1)
		util.NewURLItem(wf, "AWS region not set in "+configFilePath+" "+profileDescription).
			Subtitle("Press enter to open AWS docs on how to configure").
			Arg("https://aws.github.io/aws-sdk-go-v2/docs/configuring-sdk/#creating-the-config-file").
			Icon(aw.IconError).
			Valid(true)
	} else {
		wf.NewItem(errString).
			Icon(aw.IconError)
	}

	return errors.New(errString)
}
