package config

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

var defaultSession *aws.Config

func DefaultAwsSession() aws.Config {
	if defaultSession != nil {
		region := os.Getenv("AWS_REGION")
		cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
		if err != nil {
			log.Fatalf("unable to load AWS SDK config, %v", err)
		}
		defaultSession = &cfg
	}
	return *defaultSession

}
