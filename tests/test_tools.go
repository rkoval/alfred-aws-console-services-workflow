package tests

import (
	"gopkg.in/dnaeon/go-vcr.v4/pkg/recorder"
)

type requestQueryParamSanitizer struct {
	ParamName      string
	SanitizedValue string
}

var requestQueryParamSanitizers = []requestQueryParamSanitizer{
	{
		ParamName:      "NextToken",
		SanitizedValue: "BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB",
	},
}

func NewAWSRecorderSession(fixtureName string) *recorder.Recorder {
	r, err := recorder.New(fixtureName+"_aws_fixture",
		recorder.WithMode(recorder.ModeRecordOnce),
		recorder.WithMatcher(CustomMatcher),
		recorder.WithHook(sanitizeAndFormatBodyHook, recorder.BeforeSaveHook),
	)
	if err != nil {
		panic(err)
	}

	return r
}

func PanicOnError(f func() error) {
	if err := f(); err != nil {
		panic(err)
	}
}
