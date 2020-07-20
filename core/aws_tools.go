package core

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	aw "github.com/deanishe/awgo"
)

func LoadAWSConfig() (*session.Session, *aws.Config) {
	sess := session.Must(session.NewSession())
	cfg := &aws.Config{Region: aws.String("us-west-2")}
	return sess, cfg
}

func GetImageIcon(id string) *aw.Icon {
	icon := &aw.Icon{Value: "images/" + id + ".png"}
	return icon
}
