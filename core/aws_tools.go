package core

import (
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	aw "github.com/deanishe/awgo"
)

func LoadAWSConfig(transport http.RoundTripper) (*session.Session, *aws.Config) {
	sess := session.Must(session.NewSession())

	cfg := &aws.Config{
		Region: aws.String("us-west-2"),
	}
	if transport != nil {
		client := &http.Client{
			Transport: transport,
		}
		cfg.WithHTTPClient(client)
	}
	// cfg.WithLogLevel(aws.LogDebugWithHTTPBody)
	return sess, cfg
}

func GetImageIcon(id string) *aw.Icon {
	icon := &aw.Icon{Value: "images/" + id + ".png"}
	return icon
}
