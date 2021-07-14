package searchutil

import "github.com/aws/aws-sdk-go-v2/aws"

type SearchArgs struct {
	Query      string
	Cfg        aws.Config
	ForceFetch bool
	FullQuery  string
}
