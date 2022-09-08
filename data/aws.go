package data

import (
	"context"
	"errors"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/thatoddmailbox/jobmgr/config"
)

var s3Client *s3.Client

func initAWS() error {
	if config.Current.AWS.AccessKeyID == "" || config.Current.AWS.SecretAccessKey == "" {
		return errors.New("AWS credentials are missing from config file")
	}

	cfg, err := awsConfig.LoadDefaultConfig(
		context.Background(),
		awsConfig.WithRegion(config.Current.AWS.Region),
		awsConfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				config.Current.AWS.AccessKeyID,
				config.Current.AWS.SecretAccessKey,
				"",
			),
		),
	)
	if err != nil {
		return err
	}

	s3Client = s3.NewFromConfig(cfg)

	return nil
}
