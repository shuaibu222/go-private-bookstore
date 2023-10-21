package config

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func AwsS3Config() (*s3.S3, error) {
	Region := "af-south-1"

	config, err := LoadConfig()
	if err != nil {
		log.Println("Error while loading envs: ", err)
	}

	sess, err := session.NewSession(&aws.Config{
		Region:                        aws.String(Region),
		CredentialsChainVerboseErrors: aws.Bool(true),
		Credentials:                   credentials.NewStaticCredentials(config.Access, config.Secret, ""),
	})
	if err != nil {
		log.Println("Error creating session:", err)
	}
	svc := s3.New(sess)
	return svc, err
}
