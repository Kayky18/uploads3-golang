package main

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var (
	s3Client *s3.S3
	s3Bucket string
)

func init() {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
		Credentials: credentials.NewStaticCredentials(
			"your-access-key-id",
			"your-secret-access-key",
			""),
	})
	if err != nil {
		log.Fatalf("Failed to create session: %v", err)
	}
	s3Client = s3.New(sess)
	s3Bucket = "your-bucket-name"
}

func main() {

}
