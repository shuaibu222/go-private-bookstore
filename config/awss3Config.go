package config

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func AwsS3Config() (*s3.S3, error) {
	Region := "af-south-1"

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(Region),
	})
	if err != nil {
		log.Println("Error creating session:", err)
	}
	svc := s3.New(sess)
	return svc, err
}
