package S3

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var (
	s3Session *s3.S3
)

func initialize() {
	s3Session = s3.New(session.Must(session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
		Credentials: credentials.NewStaticCredentials(os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"), ""),
	})))
}

func ListBuckets() (resp *s3.ListBucketsOutput){
initialize()

	resp, err := s3Session.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		fmt.Println(err)
	}

	return resp
}