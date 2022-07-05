package S3

import (
	"bytes"
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

/***********************
Initializes the Amazon Simple Storage Service Session.
***********************/
func initialize() {
	s3Session = s3.New(session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(os.Getenv("AWS_REGION")),
		Credentials: credentials.NewStaticCredentials(os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"), ""),
	})))
}

/***********************
Displays the list of s3 buckets your account holds.
***********************/
func ListBuckets() (resp *s3.ListBucketsOutput) {
	initialize()

	resp, err := s3Session.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		fmt.Println(err)
	}

	return resp
}

/***********************
Uploads file on you AWS bucket.
***********************/
func UploadFile(file []byte, name string, fileSize int64) error {
	initialize()
	fmt.Println("sending...")
	res, err := s3Session.PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(os.Getenv("AWS_BUCKET_NAME")),
		Key:           aws.String(name),
		Body:          bytes.NewReader(file),
		ACL:           aws.String("public-read"),
		ContentLength: aws.Int64(fileSize),
	})
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(res)
	return nil
}
