package awsworkflow

import (
	"context"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	aw "github.com/deanishe/awgo"
)

func NewWorkflowConfig(transport http.RoundTripper) aws.Config {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}

	if transport != nil {
		cfg.HTTPClient = &http.Client{
			Transport: transport,
		}
	}
	return cfg
}

func GetImageIcon(id string) *aw.Icon {
	icon := &aw.Icon{Value: "images/" + id + ".png"}
	return icon
}
