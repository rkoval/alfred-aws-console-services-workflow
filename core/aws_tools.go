package core

import (
	"net/http"

	"github.com/aws/aws-sdk-go/aws/session"
	aw "github.com/deanishe/awgo"
)

func LoadAWSConfig(transport http.RoundTripper) *session.Session {
	sess := session.Must(session.NewSession())

	if transport != nil {
		client := &http.Client{
			Transport: transport,
		}
		sess.Config.WithHTTPClient(client)
	}
	// cfg.WithLogLevel(aws.LogDebugWithHTTPBody)
	return sess
}

func GetImageIcon(id string) *aw.Icon {
	icon := &aw.Icon{Value: "images/" + id + ".png"}
	return icon
}
