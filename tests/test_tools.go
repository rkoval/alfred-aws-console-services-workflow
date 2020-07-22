package tests

import (
	"os"

	"github.com/dnaeon/go-vcr/cassette"
	"github.com/dnaeon/go-vcr/recorder"
)

func NewAWSRecorder(fixtureLocation string) *recorder.Recorder {
	var mode recorder.Mode
	if os.Getenv("RECORD_VCR") != "" {
		mode = recorder.ModeRecording
	} else {
		mode = recorder.ModeReplaying
	}
	r, err := recorder.NewAsMode(fixtureLocation, mode, nil)
	if err != nil {
		panic(err)
	}

	r.AddFilter(func(i *cassette.Interaction) error {
		delete(i.Request.Headers, "Authorization")
		delete(i.Request.Headers, "X-Amz-Date")
		delete(i.Response.Headers, "X-Amzn-Requestid")
		delete(i.Response.Headers, "Date")
		return nil
	})

	return r
}
